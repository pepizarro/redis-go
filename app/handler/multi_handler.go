package handler

import (
	"fmt"
	"net"
)

func (h *Handler) MultiHandler(conn net.Conn, buffer []byte) {
	fmt.Println("MultiHandler..")

	_, _ = conn.Write(h.parser.WriteOk())
}

func (h *Handler) ExecHandler(conn net.Conn, buffer []byte) {

	_, _ = conn.Write(h.parser.WriteError("ERR EXEC without MULTI"))
}
