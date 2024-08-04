package handler

import (
	"fmt"
	"net"
)

func (h *Handler) MultiHandler(conn net.Conn, buffer []byte) {
	fmt.Println("MultiHandler..")

	_, _ = conn.Write(h.parser.WriteOk())
}
