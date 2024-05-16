package handler

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type Handler struct {
	store  *storage.KeySpace
	parser protocol.Parser
}

func NewHandler(store *storage.KeySpace, parser protocol.Parser) *Handler {
	return &Handler{
		store:  store,
		parser: parser,
	}
}

func (h *Handler) Handle(conn net.Conn, buffer []byte) {
	command, err := h.parser.GetCommand(buffer)
	if err != nil {
		fmt.Println("Error getting command: ", err)
		return
	}

	switch command {
	case "ping":
		h.PingHandler(conn, buffer)
	case "echo":
		h.EchoHandler(conn, buffer)
	case "set":
		h.SetHandler(conn, buffer)
	case "get":
		h.GetHandler(conn, buffer)
	default:
		fmt.Println("Unknown command: ", command)
	}
}

func (h *Handler) PingHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client
	_, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Error writing PONG to client: ", err)
		return
	}
}

func (h *Handler) EchoHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}

	message := string(params[len(params)-2])

	echo := h.parser.WriteString(message)
	fmt.Println("Echoing: ", echo)

	_, err = conn.Write([]byte(echo))

	if err != nil {
		fmt.Println("Error writing echo to client: ", err)
		return
	}
}

func (h *Handler) GetHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client

	params, err := h.parser.GetParams(buffer)
	if len(params) != 6 || err != nil {
		fmt.Println("Invalid number of arguments")
		return
	}

	key := string(params[4])

	item, err := h.store.Get(key)

	fmt.Println("Getting key: ", key, err)
	if err != nil {
		nullBulkString := h.parser.WriteNull()
		fmt.Println("Error getting key, writing: ", nullBulkString)
		_, err := conn.Write(nullBulkString)
		if err != nil {
			fmt.Println("Error writing to client: ", err)
			return
		}
		return
	}

	response := h.parser.WriteString(string(item))
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}
