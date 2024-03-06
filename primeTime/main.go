package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
)

type jsonMessage struct {
	Method string      `json:"method"`
	Number interface{} `json:"number"`
}

type respMessage struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	fmt.Println("Starting TCP Connection!")
	listener, err := net.Listen("tcp", ":11000")
	if err != nil {
		fmt.Println("Could not establish connection!", err.Error())
		return
	}
	defer listener.Close()
	cnt := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not establish listener!", err.Error())
			return
		}
		cnt++
		go handleConnection(conn, cnt)
	}
}
func handleConnection(conn net.Conn, cnt int) {
	defer conn.Close()
	buffer := make([]byte, 5 * 1024)
	for {
		// fmt.Printf("size of buffer %p\n", &buffer)
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("could not read data form connection! => ", err.Error(), cnt)
			}
			return
		}
		fmt.Printf("Received data => :%d -> %s\n", cnt, buffer[:n-1])
		var reqString = strings.Split(string(buffer[:n-1]), "\n")

		for _, val := range reqString {
			var reqData jsonMessage
			err = json.Unmarshal([]byte(val), &reqData)
			if err != nil {
				fmt.Println("vvp -> ", string(val))
				fmt.Println("data cannot be unmarshalled", err.Error(), cnt)
				conn.Write([]byte("malformed"))
				return
			}

			fmt.Println("Fetched values from json count :", reqData.Method, reqData.Number, cnt)

			if !reqData.isRequestDataValid() {
				conn.Write([]byte("malformed"))
				fmt.Println("data malformed", cnt)
				return
			} else {
				fmt.Println("data valid", cnt)
			}
			respData := respMessage{
				Method: "isPrime",
				Prime:  reqData.isNumberPrime(),
			}
			fmt.Println("Final respData => ", respData, cnt)
			respByteData, err := json.Marshal(respData)
			respByteData = append(respByteData, []byte("\n")...)
			if err != nil {
				fmt.Println("cannot marshal data to send response", cnt)
				conn.Write([]byte("malformed"))
				return
			}

			_, err = conn.Write(respByteData)
			if err != nil {
				fmt.Println("cannot send response")
				conn.Write([]byte("malformed"))
				return
			}
			fmt.Println("response sent successfully", cnt)
		}
	}
}

func (msg jsonMessage) isRequestDataValid() bool {
	// when populating data from json into an interface, Go converts the data to float64 by default, because json
	// cannot distinguish between int and float, hence to save data losses, default type is float
	if msg.Number == nil || reflect.TypeOf(msg.Number).Kind() != reflect.Float64 {
		fmt.Println("isRequestDataValid 1", reflect.TypeOf(msg.Number).Kind() != reflect.Float64)
		return false
	}
	if msg.Method != "isPrime" {
		fmt.Println("isRequestDataValid 2")
		return false
	}
	fmt.Println("isRequestDataValid 3")
	return true
}

func (msg jsonMessage) isNumberPrime() bool {
	integerVal, _ := msg.Number.(float64)

	fmt.Println("isNumberPrime 0", integerVal, float64(int(integerVal)) == float64(integerVal))
	if float64(int(integerVal)) == float64(integerVal) {
		return isPrime(int(integerVal))
	}
	fmt.Println("isNumberPrime 2")
	return false
}

func isPrime(val int) bool {
	fmt.Println("isPrime 0", val)
	if val <= 1 {
		return false
	}
	fmt.Println("isPrime 1")
	if val <= 3 {
		return true
	}
	fmt.Println("isPrime 2")
	if val%2 == 0 || val%3 == 0 {
		return false
	}
	index := 5
	for index*index <= val {
		if val%index == 0 || val%(index+2) == 0 {
			return false
		}
		index += 6
	}
	return true
}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	dataByte, err := io.ReadAll(conn)
// 	if err != nil {
// 		if err != io.EOF {
// 			fmt.Println("Error in decoding message", err.Error())
// 		} else {
// 			fmt.Println("End of file came-> ")
// 		}
// 		return
// 	}

// 	var validData = string(dataByte)
// 	// var jsonData jsonMessage
// 	// err = json.Unmarshal([]byte(validData), &jsonData)
// 	// if err != nil {
// 	// 	fmt.Println("Error in Unmarshaling", err.Error())
// 	// 	return
// 	// }
// 	fmt.Println("Received data -> ", strings.Split(validData, "\n"))
// 	// decoder := json.NewDecoder(conn)
// 	// var message jsonMessage
// 	// err := decoder.Decode(&message)
// 	// if err != nil {
// 	// 	if err != io.EOF {
// 	// 		fmt.Println("Error in decoding message", err.Error())
// 	// 	}
// 	// 	return
// 	// }
// 	// fmt.Println("Received data -> ", message)

// 	// }
// }
