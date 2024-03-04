package main

import (
    "encoding/json"
    "fmt"
    "net"
    "time"
)

type Message struct {
	Method string      `json:"method"`
	Number int `json:"number"`
}

func sendJSONRequest(conn net.Conn, message Message) error {
    // Marshal the message data to JSON
    jsonData, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %w", err)
    }

    // Append the newline character to the JSON data
    jsonData = append(jsonData, '\n')

    // Send the JSON data with newline character over the connection
    _, err = conn.Write(jsonData)
    if err != nil {
        return fmt.Errorf("failed to write data: %w", err)
    }

    return nil
}

func main() {
    // Server address
    serverAddr := "localhost:11000"

    // Dial the TCP server
    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer conn.Close()

    // Create a sample Message object
    message := Message{Method: "isPrime", Number: 2}

    // Send 10 JSON requests with a delay of 1 second between each
    for i := 0; i < 1; i++ {
        err = sendJSONRequest(conn, message)
        if err != nil {
            fmt.Println("Error sending JSON request:", err)
            return
        }
        fmt.Printf("Sent request %d\n", i+1)
        time.Sleep(1 * time.Second)
    }

    fmt.Println("Successfully sent all 10 JSON requests over TCP.")
}
