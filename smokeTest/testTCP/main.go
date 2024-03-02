package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

var wq sync.WaitGroup

func main() {
	for i := 0; i < 1000;  i++ {
		go dialTCPConn(i)
		wq.Add(1)
	}
	wq.Wait()
}

func dialTCPConn(index int) {
	defer wq.Done()
	conn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		fmt.Println("Could not connect to sent data!")
		return
	}
	defer conn.Close()

	data := []byte(fmt.Sprintf("Hello world guys %s", strconv.Itoa(index)))
	n, err := conn.Write(data)

	if err != nil {
		fmt.Println("Could not write!")
		return
	}
	fmt.Printf("data written succesfully %v\n", n)
}
