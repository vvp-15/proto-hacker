package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"sync"
)

type client struct {
	name string
	conn net.Conn
}

var (
	clients = make(map[net.Conn]*client)
	mutex   sync.Mutex
)


func main() {
	fmt.Println("starting server")

	listener, err := net.Listen("tcp", ":15004")
	if err != nil {
		fmt.Println("tcp chala hi nahi")
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("sunnte wqt fat gya")
			return
		}
		fmt.Println("new connection aaya ", conn.LocalAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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
				mutex.Lock()
				clients[conn] = &client{
					name: msg,
					conn: conn,
				}
				mutex.Unlock()
				fmt.Println("username is valid")
				broadcastMsg(fmt.Sprintf("* %v has entered the room", msg), conn)
				sendUserRoomStatus(conn)
			} else {
				fmt.Println("username is not valid")
				conn.Write([]byte("Invalid username"))
				return
			}
		} else {
			fmt.Println("msg ye likha-", msg)
			broadcastMsg(fmt.Sprintf("[%v] %v", clients[conn].name, msg), conn)
		}
	}
	// fmt.Println("connection close hogya", clients[conn].name)
	mutex.Lock()
	name := clients[conn].name
	delete(clients, conn)
	mutex.Unlock()
	broadcastMsg(fmt.Sprintf("* %v has left the room",name), conn)
}

func checkIfValidUserName(msg string) bool {
	pattern := "^[a-zA-Z0-9]*$"
	match, _ := regexp.MatchString(pattern, msg)
	return match
}

func broadcastMsg(msg string, conn net.Conn) {
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

func sendUserRoomStatus(conn net.Conn) {
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
