package protocol

import (
	"bytes"
	"fmt"
	"strings"
)

type RESP struct {
}

func NewRESP() *RESP {
	return &RESP{}
}

func (r *RESP) GetCommand(buffer []byte) (string, error) {
	params, err := r.GetParams(buffer)
	if err != nil {
		return "", err
	}
	command := strings.ToLower(string(params[2]))

	return command, nil
}

func (r *RESP) GetSubCommand(buffer []byte) (string, error) {

	params, err := r.GetParams(buffer)
	if err != nil {
		return "", err
	}

	// get the subcommand from a sub slice
	params = append(params[:1], params[3:]...)
	var newBuffer []byte

	for _, item := range params {
		newBuffer = append(newBuffer, item...)
		newBuffer = append(newBuffer, '\r', '\n')
	}

	subCommand, err := r.GetCommand(newBuffer)
	if err != nil {
		return "", err
	}

	return subCommand, nil
}

func (r *RESP) GetParams(buffer []byte) ([][]byte, error) {
	newBuffer := make([]byte, len(buffer))
	copy(newBuffer, buffer)
	params := bytes.Split(newBuffer, []byte{'\r', '\n'})
	if len(params) < 3 {
		return nil, fmt.Errorf("Invalid command: %s", newBuffer)
	}
	return params, nil
}

func (r *RESP) GetArgs(buffer []byte, start int) (map[string][]byte, error) {
	params, err := r.GetParams(buffer)
	if err != nil {
		return nil, err
	}

	if len(params) < start+2 {
		return nil, fmt.Errorf("Invalid command: %s", buffer)
	}

	args := make(map[string][]byte)

	for i := start + 2; i < len(params); i += 4 {
		key := strings.ToLower(string(params[i]))
		if len(params) <= i+2 {
			args[key] = nil
			break
		}
		args[key] = params[i+2]
	}

	return args, nil
}

func (r *RESP) GetType(buffer []byte) (string, error) {
	fmt.Println("In GetType, buffer: ", buffer)
	// valueType := buffer[0]

	// switch valueType {
	// case '+':
	// 	return "string", nil
	// case '~':
	// 	return "set", nil
	// case '%':
	// 	return "hash", nil
	// }

	return "string", nil
}

func (r *RESP) WriteString(s string) []byte {
	length := len(s)
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", length, s))
}

func (r *RESP) WriteNull() []byte {
	return []byte("$-1\r\n")
}

func (r *RESP) WriteOk() []byte {
	return []byte("+OK\r\n")
}

func (r *RESP) WriteError(message string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", message))
}

func (r *RESP) WriteArray(array [][]byte) []byte {
	length := len(array)

	var buffer bytes.Buffer
	for _, item := range array {
		buffer.Write(item)
	}

	return []byte(fmt.Sprintf("*%d\r\n%s", length, buffer.String()))

}
