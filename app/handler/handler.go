package handler

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

const (
	PING   = "ping"
	ECHO   = "echo"
	SET    = "set"
	GET    = "get"
	CONFIG = "config"
	KEYS   = "keys"
	TYPE   = "type"

	XADD   = "xadd"
	XRANGE = "xrange"
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
		conn.Close()
		return
	}

	switch command {
	case PING:
		h.PingHandler(conn, buffer)
	case ECHO:
		h.EchoHandler(conn, buffer)
	case SET:
		h.SetHandler(conn, buffer)
	case GET:
		h.GetHandler(conn, buffer)
	case CONFIG:
		h.ConfigHandler(conn, buffer)
	case KEYS:
		h.KeysHandler(conn, buffer)
	case TYPE:
		h.TypeHandler(conn, buffer)
	case XADD:
		h.XaddHandler(conn, buffer)
	case XRANGE:
		h.XrangeHandler(conn, buffer)

	default:
		fmt.Println("Unknown command: ", command)
		conn.Close()
		return
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

func (h *Handler) KeysHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params in KeysHandler: ", err)
		return
	}

	if string(params[len(params)-2][0]) != "*" {
		fmt.Println("Param not recognized")
		return
	}

	keys, err := h.store.GetAllKeys()
	if err != nil {
		return
	}

	var keysArray [][]byte
	for _, key := range keys {
		keysArray = append(keysArray, h.parser.WriteString(key))
	}

	keysByteArray := h.parser.WriteArray(keysArray)
	_, err = conn.Write(keysByteArray)
	if err != nil {
		return
	}
}
