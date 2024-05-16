package main

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/handler"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func main() {

	store := storage.NewKeySpace()
	go store.LogKeySpace()
	parser := protocol.NewRESP()

	handler := handler.NewHandler(store, parser)

	server := NewRedisServer("0.0.0.0", "6379", handler)

	fmt.Println("Starting redis server...")

	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server: ", err.Error())
	}
}
