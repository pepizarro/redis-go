package handler

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"
)

type Handler struct {
	mu       sync.RWMutex
	keyspace map[string][]byte
}

func NewHandler() *Handler {
	return &Handler{
		keyspace: make(map[string][]byte),
	}
}

func (h *Handler) logKeySpace() {
	for {
		h.mu.RLock()
		fmt.Println("KeySpace: ")
		for key, value := range h.keyspace {
			fmt.Println(key, ":", string(value))
		}
		h.mu.RUnlock()

		time.Sleep(5 * time.Second)
	}
}

func (h *Handler) getParams(buffer []byte) [][]byte {
	lines := bytes.Split(buffer, []byte{'\r', '\n'})

	return lines
}

func (h *Handler) EchoHandler(conn net.Conn, buffer []byte) {

	lines := h.getParams(buffer)
	message := string(lines[len(lines)-2])

	echo := fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
	fmt.Println("Echoing: ", echo)

	_, err := conn.Write([]byte(echo))

	if err != nil {
		fmt.Println("Error writing to client: ", err)
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

func (h *Handler) SetHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client

	lines := h.getParams(buffer)
	if len(lines) != 8 {
		fmt.Println("Invalid number of arguments")
		return
	}

	key := string(lines[4])
	value := lines[6]

	h.mu.Lock()
	defer h.mu.Unlock()
	h.keyspace[key] = value

	// fmt.Println("key: ", key)
	// fmt.Println("value: ", value)

	response := "+OK\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}

func (h *Handler) GetHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client

	params := h.getParams(buffer)
	if len(params) != 6 {
		fmt.Println("Invalid number of arguments")
		return
	}

	key := string(params[4])

	h.mu.RLock()
	defer h.mu.RUnlock()
	value, ok := h.keyspace[key]

	if !ok {
		response := "$_\r\n"
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
			return
		}
		return
	}

	response := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}
