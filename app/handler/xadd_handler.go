package handler

import (
	"fmt"
	"net"
)

func (h *Handler) XaddHandler(conn net.Conn, buffer []byte) {
	fmt.Println("XADD Handler")

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}

	if len(params) < 9 {
		fmt.Println("Invalid number of arguments")
		return
	}

	key := string(params[4])
	entryID := string(params[6])
	if entryID == "*" {
		fmt.Println("Entry ID is *")
	}

	// get the key-values pairs
	args, err := h.parser.GetArgs(buffer, 6)
	if err != nil {
		fmt.Println("Error getting args: ", err)
		return
	}

	fmt.Println("args: ", args)

	// create stream
	entry := h.store.NewEntry(entryID, args)

	err = h.store.SetStream(key, entry)
	if err != nil {
		fmt.Println("Error setting stream: ", err)
		_, err = conn.Write(h.parser.WriteString(err.Error()))
	}

	_, err = conn.Write(h.parser.WriteString(entryID))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}
