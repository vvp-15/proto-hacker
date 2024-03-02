package main

import (
	"fmt"
	"io"
	"net"
)

func main() {

	fmt.Println("Starting TCP Connection!")

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("Could not establish connection!", err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not establish listener!", err.Error())
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			fmt.Println("could not read data form connection! => ", err.Error())
		}
		return
	}
	fmt.Printf("Received data => : %s\n", buffer[:n])

	respData := string(buffer[:n])

	_, err = conn.Write([]byte(respData))
	if err != nil {
		fmt.Println("Error sending back data => :", err.Error())
	} else {
		fmt.Println("dats sent succesfully")
	}
}