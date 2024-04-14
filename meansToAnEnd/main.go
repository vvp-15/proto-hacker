package main

import (
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

	for {
		conn, err2 := listener.Accept()
		if err2 != nil {
			fmt.Printf("Listen krte wqt fatt gya")
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 9)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("could not read data form connection! => ", err.Error())
			}
			return
		}

		if n != 9 {
			continue
		}

		// first byte ASCII = I or Q
		
	
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