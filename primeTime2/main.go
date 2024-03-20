package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"
)

type primeRequest struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type primeResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	fmt.Println("Starting TCP Connection!")
	listener, err := net.Listen("tcp", ":11000")
	if err != nil {
		fmt.Println("Could not establish connection!", err.Error())
		return
	}
	defer listener.Close()
	cnt := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not establish listener!", err.Error())
			return
		}
		cnt++
		go handleConnection(conn, cnt)
	}
}
func handleConnection(conn net.Conn, cnt int) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Bytes()

		log.Println("received:", string(line))

		resBytes, valid, err := handleLine(line)
		if err != nil {
			return
		}

		if _, err := conn.Write(resBytes); err != nil {
			return
		}

		// stop processing if the request was invalid
		if !valid {
			break
		}
	}
}

func handleLine(line []byte) ([]byte, bool, error) {
	var req primeRequest
	if err := json.Unmarshal(line, &req); err != nil || !isValidPrimeRequest(req) {
		return []byte("invalid request\n"), false, nil
	}
	log.Println("response:", primeResponse{Method: "isPrime", Prime: isPrime(*req.Number)})
	resBytes, err := json.Marshal(primeResponse{Method: "isPrime", Prime: isPrime(*req.Number)})
	if err != nil {
		return nil, true, err
	}

	return append(resBytes, []byte("\n")...), true, nil
}

func isValidPrimeRequest(req primeRequest) bool {
	return req.Method == "isPrime" && req.Number != nil
}

func isPrime(n float64) bool {
	// prime numbers are positive integers
	if n < 0 || n != math.Trunc(n) {
		return false
	}
	return big.NewInt(int64(n)).ProbablyPrime(20)
}
