package main

import (
	"fmt"
	"net"
)

func main() {

	fmt.Println("Starting TCP Connection!")

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("Could not establish connection!")
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not establish listener!")
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
		fmt.Println("could not read data form connection!")
		return
	}
	fmt.Printf("Received data => : %s\n", buffer[:n])
}