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

	valueType, err := h.store.GetType(key)
	if err != nil {
		_, err = conn.Write(h.parser.WriteString("none"))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
		}
	}
	fmt.Println("Key: ", key)
	fmt.Println("Type: ", valueType)

	_, err = conn.Write(h.parser.WriteString(valueType))
	if err != nil {
		return
	}

}
