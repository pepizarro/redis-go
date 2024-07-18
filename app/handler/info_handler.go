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

	var info map[string]string

	switch section {
	case "replication":
		replicationInfo := h.config.getReplicationInfo()
		info = replicationInfo
	}

	fmt.Println("Info: ", info)
	var response string
	for key, value := range info {
		response = response + key + ":" + value + "\r\n"
	}

	// fmt.Println("Response: ", string(response))

	_, err = conn.Write(h.parser.WriteString(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}
