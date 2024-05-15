package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server/handlers"
)

type Server struct {
	address  string
	port     string
	protocol resp.RESP
}

func NewServer(address string, port string) *Server {
	return &Server{
		address: address,
		port:    port,
	}
}

func (s *Server) Start() error {

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error starting server in: ", s.address, s.port)
		return err
	}

	defer listener.Close()

	fmt.Println("Server started at: ", s.address, s.port)

	for {

		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go route(conn, s.protocol)
	}
}

func route(conn net.Conn, protocol resp.RESP) {
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

		// Get the command
		command, err := protocol.GetCommand(buffer[:n])
		if err != nil {
			fmt.Println("Error getting command:", err)
			return
		}

		switch command {
		case "PING":
			handlers.PingHandler(conn, buffer[:n])
		case "ECHO":
			handlers.EchoHandler(conn, buffer[:n])
		default:
			fmt.Println("Command not found")
		}

	}
}
