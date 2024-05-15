package main

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/server"
	// Uncomment this block to pass the first stage
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.

	server := server.NewServer("0.0.0.0", "6379")

	fmt.Println("Starting redis server...")

	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server: ", err.Error())
	}
}
