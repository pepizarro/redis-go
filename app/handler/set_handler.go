package handler

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func (h *Handler) SetHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}

	key := string(params[4])
	value := params[6]

	args, err := h.parser.GetArgs(buffer, 6)
	if err != nil {
		fmt.Println("Error getting args: ", err)
		return
	}

	fmt.Println("args: ", args)

	if len(args) == 0 {
		h.store.Set(key, "string", value)
	}

	// Check for known arguments
	expArgs := []string{"ex", "px"}

	for _, arg := range expArgs {
		if _, ok := args[arg]; ok {
			expTime, err := strconv.Atoi(string(args[arg]))
			if err != nil {
				null := h.parser.WriteNull()
				_, err := conn.Write(null)
				if err != nil {
					fmt.Println("Error writing to client: ", err)
				}
				return
			}
			timeToAdd := time.Duration(0)
			switch arg {
			case "ex":
				timeToAdd = time.Duration(expTime) * time.Second
			case "px":
				timeToAdd = time.Duration(expTime) * time.Millisecond
			}

			h.store.SetWithExpiration(key, "string", value, time.Now().Add(timeToAdd))
		}
	}

	success := h.parser.WriteOk()
	_, err = conn.Write(success)
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}
