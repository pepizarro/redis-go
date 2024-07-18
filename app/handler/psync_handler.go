package handler

import (
	"encoding/hex"
	"fmt"
	"net"
)

func (h *Handler) PsyncHandler(conn net.Conn, buffer []byte) {

	fmt.Println("Psync Handler")
	response := fmt.Sprintf("FULLRESYNC %s %d", h.config.replicationID, h.config.replicationOffset)
	_, err := conn.Write(h.parser.WriteSimpleString(response))
	// err := h.sendAndRead(conn, h.parser.WriteSimpleString(response))
	if err != nil {
		fmt.Println("Error sending: ", err)
		return
	}

	// send the empty RDB file
	file := "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
	bytes, err := hex.DecodeString(file)
	if err != nil {
		fmt.Println("Error decoding hex string: ", err)
		return
	}

	_, err = conn.Write(h.parser.WriteFile(bytes))
	// err = h.sendAndRead(conn, h.parser.WriteFile(bytes))
	if err != nil {
		fmt.Println("Error sending and reading: ", err)
		return
	}

}
