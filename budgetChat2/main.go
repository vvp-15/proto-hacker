package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
)

type client struct {
	name string
	conn net.Conn
}

func main() {
	fmt.Println("starting server")

	listener, err := net.Listen("tcp", ":15004")
	if err != nil {
		fmt.Println("tcp chala hi nahi")
		return
	}
	defer listener.Close()
	clients := make(map[net.Conn]*client)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("sunnte wqt fat gya")
			return
		}
		fmt.Println("new connection aaya ", conn.LocalAddr().String())
		go handleConnection(conn, clients)
	}
}

func handleConnection(conn net.Conn, clients map[net.Conn]*client) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	_, err := conn.Write([]byte("Welcome to budgetchat! What shall I call you?"))
	if err != nil {
		fmt.Println("msg broadcast nahi kr pae")
	}
	gotName := false
	for scanner.Scan() {
		msg := scanner.Text()
		if !gotName {
			fmt.Println("name kya hai-", msg)
			if checkIfValidUserName(msg) {
				gotName = true
				clients[conn] = &client{
					name: msg,
					conn: conn,
				}
				fmt.Println("username is valid")
				broadcastMsg(fmt.Sprintf("* %v has entered the room", msg), conn, clients)
				sendUserRoomStatus(conn, clients)
			} else {
				fmt.Println("username is not valid")
				conn.Write([]byte("Invalid username"))
				return
			}
		} else {
			fmt.Println("msg ye likha-", msg)
			broadcastMsg(fmt.Sprintf("[%v] %v", clients[conn].name, msg), conn, clients)
		}
	}
	// fmt.Println("connection close hogya", clients[conn].name)
	broadcastMsg(fmt.Sprintf("* %v has left the room", clients[conn].name), conn, clients)
	delete(clients, conn)
}

func checkIfValidUserName(msg string) bool {
	pattern := "^[a-zA-Z0-9]*$"
	match, _ := regexp.MatchString(pattern, msg)
	return match
}

func broadcastMsg(msg string, conn net.Conn, clients map[net.Conn]*client) {
	// _, err := conn.Write([]byte(msg))
	// if err != nil {
	// 	fmt.Println("msg broadcast nahi kr pae")
	// }
	for key, _ := range clients {
		if key != conn {
			_, err := key.Write([]byte(msg))
			if err != nil {
				fmt.Println("msg broadcast nahi kr pae")
			}
		}
	}
}

func sendUserRoomStatus(conn net.Conn, clients map[net.Conn]*client) {
	msg := "* The room contains:"
	for key := range clients {
		if key == conn {
			continue
		}
		msg += fmt.Sprintf(" %v,", clients[key].name)
	}
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Room status bhejte wqt fatt gaess")
	}
}
