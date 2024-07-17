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

		_, _ = conn.Write(h.parser.WriteOk())
	case "capa":
		fmt.Println("Capa: ", string(params[6]))
		_, _ = conn.Write(h.parser.WriteOk())
	case "psync":
		response := fmt.Sprintf("FULLRESYNC %s %d", h.config.replicationID, h.config.replicationOffset)
		_, _ = conn.Write(h.parser.WriteSimpleString(response))
	}

}
