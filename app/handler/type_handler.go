package handler

import (
	"fmt"
	"net"
)

func (h *Handler) TypeHandler(conn net.Conn, buffer []byte) {
	fmt.Println("TypeHandler")

	// get the requested key
	params, err := h.parser.GetParams(buffer)
	if err != nil {
		return
	}
	key := string(params[4])

	value, err := h.store.Get(key)
	if err != nil {
		none := h.parser.WriteString("none")
		_, err := conn.Write([]byte(none))
		if err != nil {
			return
		}
		return
	}

	valueType, err := h.parser.GetType(value)
	_, err = conn.Write([]byte(h.parser.WriteString(valueType)))
	if err != nil {
		return
	}

}
