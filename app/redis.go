package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/handler"
)

type RedisServer struct {
	address string
	port    string
	handler *handler.Handler
}

func NewRedisServer(address string, port string, handler *handler.Handler) *RedisServer {
	return &RedisServer{
		address: address,
		port:    port,
		handler: handler,
	}
}

func (rs *RedisServer) Start() error {

	listener, err := net.Listen("tcp", rs.address+":"+rs.port)
	if err != nil {
		fmt.Println("Error starting server in: ", rs.address, rs.port)
		return err
	}

	defer listener.Close()

	fmt.Println("Server started at: ", rs.address, rs.port)

	for {

		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go route(conn, rs.handler)
	}
}

func route(conn net.Conn, handler *handler.Handler) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			fmt.Println("Error:", err)
			return
		}
		if n == 0 {
			return
		}

		// fmt.Println("\nReceived:\n", string(buffer[:n]))

		handler.Handle(conn, buffer[:n])

	}
}
