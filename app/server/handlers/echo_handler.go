package handlers

import (
	"bytes"
	"fmt"
	"net"
)

func EchoHandler(conn net.Conn, buffer []byte) {
	// Write the data back to the client

	message := getMessage(buffer)

	echo := fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
	fmt.Println("Echoing: ", echo)
	_, err := conn.Write([]byte(echo))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}

func getMessage(payload []byte) string {

	lines := bytes.Split(payload, []byte{'\r', '\n'})

	// return the last line as a string
	return string(lines[len(lines)-2])
}
