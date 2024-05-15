package handlers

import (
	"fmt"
	"net"
)

func PingHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client

	_, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Error writing PONG to client: ", err)
		return
	}
}
