package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
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

func receiveResponse(conn net.Conn) (Response, error) {
	// Create a buffer to hold the response data
	buffer := make([]byte, 1024)

	// Read data from the connection until a newline character is found
	n, err := conn.Read(buffer)
	if err != nil {
		return Response{}, fmt.Errorf("failed to read data: %w", err)
	}

	// Find the first occurrence of the newline character
	newlineIndex := bytes.Index(buffer[:n], []byte("\n"))
	if newlineIndex == -1 {
		return Response{}, fmt.Errorf("could not find newline character in response")
	}

	// Extract the response data (excluding the newline character)
	responseData := make([]byte, newlineIndex)
	copy(responseData, buffer[:newlineIndex])

	// Unmarshal the response data to a Response object
	var response Response
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return Response{}, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return response, nil
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
	for i := 0; i < 10; i++ {
		err = sendJSONRequest(conn, message)
		if err != nil {
			fmt.Println("Error sending JSON request:", err)
			return
		}
		fmt.Printf("Sent request %d\n", i+1)

		// Receive response from server
		response, err := receiveResponse(conn)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			continue // Skip to next request if there's an error
		}

		fmt.Printf("Received response: Method: %s, Prime: %v\n", response.Method, response.Prime)

		time.Sleep(1 * time.Second)
	}

	fmt.Println("Successfully sent 10 JSON requests and received responses over TCP.")
}
