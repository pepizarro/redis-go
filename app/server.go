package main

import (
	"fmt"
	"strings"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on port 6379")
	for {

		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	// Add your code to handle the connection here
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if strings.Contains(string(buffer[:n]), "PING") {
			_, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				return
			}
		}

		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s\n", buffer[:n])
	}
}
