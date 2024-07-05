package handler

import (
	"fmt"
	"net"
)

func (h *Handler) InfoHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		return
	}

	fmt.Println("Info Handler")
	fmt.Println("Params: ", params)
	var section string
	if len(params) == 6 {
		section = string(params[4])
	}

	info, err := h.store.GetInfo(section)
	if err != nil {
		_, err = conn.Write(h.parser.WriteError(err.Error()))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
		}
		return
	}

	fmt.Println("Info: ", info)
	var response [][]byte
	for key, value := range info {
		response = append(response, h.parser.WriteString(fmt.Sprintf("%s:%s", key, value)))
	}

	_, err = conn.Write(response[0])
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}

}
