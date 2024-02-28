package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		fmt.Println("Could not connect to sent data!")
		return
	}
	defer conn.Close()

	data := []byte("Hello World guys")
	n, err := conn.Write(data)

	if err != nil {
		fmt.Println("Could not write!")
		return
	}
	fmt.Printf("data written succesfully %v", n)
}
