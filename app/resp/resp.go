package resp

import (
	"bytes"
	"fmt"
	"strings"
)

type RESP struct{}

func (r *RESP) GetCommand(buffer []byte) (string, error) {

	lines := bytes.Split(buffer, []byte{'\r', '\n'})
	if len(lines) < 3 {
		return "", fmt.Errorf("Invalid command: %s", buffer)
	}
	command := string(lines[2])
	fmt.Println("Command: ", command)

	switch {
	case strings.EqualFold(command, "PING"):
		return "PING", nil
	case strings.EqualFold(command, "ECHO"):
		return "ECHO", nil

	default:
		return "", fmt.Errorf("Unknown command: %s", command)
	}
}
