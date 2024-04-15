package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
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

		// first byte ASCII = I or Q
		reqType := string(buffer[0])
		firstValue := int32(binary.BigEndian.Uint32(buffer[1:5]))
		secondValue := int32(binary.BigEndian.Uint32(buffer[6:9]))

		fmt.Println("splitted response data -> ", reqType, firstValue, secondValue)
		
		fmt.Println("response data n -> ", n);
		respData := string(buffer[:n])
		fmt.Println("response data -> ", respData);
		_, err = conn.Write([]byte(respData))
		if err != nil {
			fmt.Println("Error sending back data => :", err.Error())
		} else {
			fmt.Println("dats sent succesfully")
		}
	}
}