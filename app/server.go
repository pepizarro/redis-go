package main

import (
	"flag"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/handler"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func main() {

	// get flag --dir
	// get flag --filename
	dirPtr := flag.String("dir", "/var/snap/redis/1446", "The directory to store the database files.")
	dbfilenamePtr := flag.String("dbfilename", "dump.rdb", "The name of the database file.")

	flag.Parse()

	config := storage.NewConfig(*dirPtr, *dbfilenamePtr)

	store := storage.NewKeySpace(config)
	parser := protocol.NewRESP()

	handler := handler.NewHandler(store, parser)

	server := NewRedisServer("0.0.0.0", "6379", handler)

	fmt.Println("Starting redis server...")

	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server: ", err.Error())
	}
}
