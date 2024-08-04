package handler

import (
	"fmt"
	"net"
	"strings"
)

func (h *Handler) ReplconfHandler(conn net.Conn, buffer []byte) {

	fmt.Println("Replconf Handler")
	params, err := h.parser.GetParams(buffer)
	if err != nil {
		return
	}

	subCommand := strings.ToLower(string(params[4]))

	switch subCommand {
	case "listening-port":
		fmt.Println("Listening port: ", string(params[6]))
		fmt.Println("Replica: ", conn.RemoteAddr())
		h.replicas = append(h.replicas, &Replica{
			conn:    conn,
			updated: false,
		})
		_, _ = conn.Write(h.parser.WriteOk())

	case "capa":
		fmt.Println("Capa: ", string(params[6]))
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
		fmt.Println("ACK handler")
		fmt.Println("Replica: ", conn.RemoteAddr())
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
