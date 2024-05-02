package main

import (
	"bufio"
	"fmt"
	"net"
)

type client struct {
	name   string
	conn   net.Conn
	reader *bufio.Reader
}

func main() {
	ln, err := net.Listen("tcp", ":8080") // Listen on port 8080
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	clients := make(map[net.Conn]*client) // Store connected clients

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			fmt.Println("Client connected:", conn.RemoteAddr())

			// Add new client to map
			reader := bufio.NewReader(conn)
			clients[conn] = &client{conn: conn, reader: reader}

			go handleClient(conn, clients)
		}
	}()

	fmt.Println("Server listening on port 8080...")
	scan := bufio.NewScanner(ln) // For potential future server commands
	for scan.Scan() {
		// Handle server commands (not implemented in this example)
	}
}

func handleClient(conn net.Conn, clients map[net.Conn]*client) {
	defer conn.Close()

	client := clients[conn]
	fmt.Println("Waiting for name...")
	name, err := client.reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	client.name = name[:len(name)-1] // Remove newline from name

	fmt.Println(client.name, "joined the chat!")
	broadcast(client.name+" joined the chat!", clients)

	for {
		message, err := client.reader.ReadString('\n')
		if err != nil {
			fmt.Println(client.name, "disconnected:", err)
			delete(clients, conn)
			broadcast(client.name+" left the chat!", clients)
			break
		}

		fmt.Println(client.name + ": " + message)
		broadcast(message, clients)
	}
}

func broadcast(message string, clients map[net.Conn]*client) {
	for conn, c := range clients {
		if conn != c.conn { // Don't send to self
			fmt.Fprintln(c.conn, message)
		}
	}
}
