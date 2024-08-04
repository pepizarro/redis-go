package handler

import (
	"fmt"
	"net"
	"strings"
)

func (h *Handler) ReplconfHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		return
	}

	subCommand := strings.ToLower(string(params[4]))

	switch subCommand {
	case "listening-port":
		h.replicas = append(h.replicas, &Replica{
			conn:    conn,
			updated: false,
		})
		_, _ = conn.Write(h.parser.WriteOk())

	case "capa":
		_, _ = conn.Write(h.parser.WriteOk())
		for _, replica := range h.replicas {
			if replica.conn.RemoteAddr() == conn.RemoteAddr() {
				replica.mu.Lock()
				fmt.Println("Replica updated: ", replica)
				replica.updated = true
				replica.mu.Unlock()
				break
			}
		}

	case "ack":
		for _, replica := range h.replicas {
			if replica.conn.RemoteAddr() == conn.RemoteAddr() {
				replica.mu.Lock()
				fmt.Println("Replica updated: ", replica)
				replica.updated = true
				replica.mu.Unlock()
				break
			}
		}
	}

}
