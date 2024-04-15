package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

type InsertData struct {
	ts int32
	price int32
}

type QueryData struct {
	minTs int32
	maxTs int32
}

func main() {
	fmt.Println("Call in main")

	listener, err := net.Listen("tcp", ":11000")

	if err != nil {
		fmt.Printf("Kuch fat gya")
		return
	}
	defer listener.Close()
	for {
		conn, err2 := listener.Accept()
		if err2 != nil {
			fmt.Printf("Listen krte wqt fatt gya")
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("handleConnection 1")
	data := make([]InsertData, 0)
	defer conn.Close()
	buffer := make([]byte, 9)
	for {
		n, err := io.ReadFull(conn, buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("could not read data form connection! => ", err.Error())
			}
			return
		}
		fmt.Println("response data n -> ", n);
		// first byte ASCII = I or Q
		reqType := string(buffer[0])
		firstValue := int32(binary.BigEndian.Uint32(buffer[1:5]))
		secondValue := int32(binary.BigEndian.Uint32(buffer[5:9]))

		fmt.Println("splitted response data -> ", reqType, firstValue, secondValue)
		
		if reqType == "I" {
			data = append(data, InsertData{
				ts : firstValue,
				price: secondValue,
			})
		} else if reqType == "Q" {
			var respData, cnt int
			for temp := range data {
				fmt.Println("Element", data[temp].ts, data[temp].price)
				dataTs := data[temp].ts
				if dataTs >= firstValue && dataTs <= secondValue {
					cnt ++;
					respData += int(data[temp].price)
				}
			}
			fmt.Println("response mean -> ", strconv.Itoa(respData / cnt));
			binary.Write(conn, binary.BigEndian, int32(respData / cnt))
		} else {
			_, err = conn.Write([]byte("failed"))
			if err != nil {
				fmt.Println("Error l sending back data => :", err.Error())
			} else {
				fmt.Println("dats  l sent succesfully")
			}
		}
	}
}