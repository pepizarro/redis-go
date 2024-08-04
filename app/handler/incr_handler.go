package handler

import (
	"fmt"
	"net"
	"strconv"
)

func (h *Handler) IncrHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}

	key := string(params[4])

	value, err := h.store.Get(key)
	if err != nil {
		fmt.Println("Error getting key: ", err)
		// defaulting to 1
		h.store.Set(key, "string", []byte("1"))
		response := h.parser.WriteInteger(1)
		_, _ = conn.Write(response)
		return
	}

	if !IsReal(value) {
		fmt.Println("Value is not a number")
		errMsg := h.parser.WriteError("ERR value is not an integer or out of range")
		_, _ = conn.Write(errMsg)
		return
	}

	num, err := strconv.Atoi(string(value))
	if err != nil {
		fmt.Println("Error converting value to int: ", err)
		return
	}

	h.store.Set(key, "string", []byte(strconv.Itoa(num+1)))

	response := h.parser.WriteInteger(num + 1)
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}

}

func IsReal(n []byte) bool {
	if len(n) > 0 && n[0] == '-' {
		n = n[1:]
	}
	if len(n) == 0 {
		return false
	}
	var point bool
	for _, c := range n {
		if '0' <= c && c <= '9' {
			continue
		}
		if c == '.' && len(n) > 1 && !point {
			point = true
			continue
		}
		return false
	}
	return true
}
