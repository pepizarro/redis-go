package handler

import (
	"fmt"
	"net"
)

func (h *Handler) PsyncHandler(conn net.Conn, buffer []byte) {

	response := fmt.Sprintf("FULLRESYNC %s %d", h.config.replicationID, h.config.replicationOffset)
	fmt.Println("Response: ", response)
	_, _ = conn.Write(h.parser.WriteSimpleString(response))
}
