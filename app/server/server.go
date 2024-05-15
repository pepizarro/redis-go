package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server/handler"
)

type Server struct {
	address  string
	port     string
	protocol resp.RESP
	h        *handler.Handler
}

func NewServer(address string, port string, h *handler.Handler) *Server {
	return &Server{
		address: address,
		port:    port,
		h:       h,
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

		go route(conn, s.protocol, s.h)
	}
}

func route(conn net.Conn, protocol resp.RESP, h *handler.Handler) {
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

		fmt.Println("\nReceived:\n", string(buffer[:n]))

		// Get the command
		command, err := protocol.GetCommand(buffer[:n])
		if err != nil {
			errMsg := fmt.Sprintf("-%s\r\n", err.Error())
			_, err := conn.Write([]byte(errMsg))
			if err != nil {
				fmt.Println("Error writing to client: ", err)
			}

			return
		}

		switch command {
		case "PING":
			h.PingHandler(conn, buffer[:n])
		case "ECHO":
			h.EchoHandler(conn, buffer[:n])
		case "SET":
			h.SetHandler(conn, buffer[:n])
		case "GET":
			h.GetHandler(conn, buffer[:n])
		default:
			fmt.Println("Command not found")
		}

	}
}
