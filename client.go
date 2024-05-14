package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	// Send data to the server

	fmt.Println("Sending ping...")
	data := []byte("*1\r\n$4\r\nPING\r\n")
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Read and process data from the server
	// ...
	// Read and process data from the server
	response := make([]byte, 1024)
	for {
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Response from server:", string(response[:n]))
	}
}
