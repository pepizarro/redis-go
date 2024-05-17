package handler

import (
	"fmt"
	"net"
)

func (h *Handler) ConfigHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client
	subCommand, err := h.parser.GetSubCommand(buffer)
	if err != nil {
		return
	}

	fmt.Println("SubCommand: ", subCommand)

	switch subCommand {
	case "get":
		h.ConfigGetHandler(conn, buffer)
	default:
		return
	}
}

func (h *Handler) ConfigGetHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client

	params, err := h.parser.GetParams(buffer)
	key := string(params[6])

	var configArray [][]byte
	configArray = append(configArray, h.parser.WriteString(key))

	config := h.store.GetConfig()
	switch string(key) {
	case "dir":
		configArray = append(configArray, h.parser.WriteString(config.Dir))
	case "dbfilename":
		configArray = append(configArray, h.parser.WriteString(config.DBfilename))
	default:
		configArray = append(configArray, h.parser.WriteNull())
	}

	response := h.parser.WriteArray(configArray)

	_, err = conn.Write(response)
	if err != nil {
		return
	}
}
