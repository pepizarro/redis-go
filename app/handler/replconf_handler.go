package handler

import (
	"fmt"
	"net"
)

func (h *Handler) ReplconfHandler(conn net.Conn, buffer []byte) {

	fmt.Println("Replconf Handler")
	params, err := h.parser.GetParams(buffer)
	if err != nil {
		return
	}

	if string(params[4]) == "listening-port" {
		fmt.Println("Listening port: ", string(params[6]))
	}

	fmt.Println("Params: ", params)
}
