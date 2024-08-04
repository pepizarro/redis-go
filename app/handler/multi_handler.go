package handler

import (
	"fmt"
	"net"
)

func (h *Handler) MultiHandler(conn net.Conn, buffer []byte) {
	fmt.Println("MultiHandler..")

	h.appendTransaction(conn)

	fmt.Println("Transaction appended..")

	_, _ = conn.Write(h.parser.WriteOk())
}

func (h *Handler) ExecHandler(conn net.Conn) {
	fmt.Println("ExecHandler..")

	var transaction *Transaction
	for _, t := range h.transactions {
		if t.conn == conn {
			transaction = t
			break
		}
	}

	if transaction == nil {
		_, _ = conn.Write(h.parser.WriteError("ERR EXEC without MULTI"))
		return
	}

	fmt.Println("Transaction: ", transaction.conn)
	commands := transaction.commands
	for _, cmd := range commands {
		fmt.Println("Command: ", string(cmd))
	}

	if len(commands) == 0 {
		fmt.Println("No commands in transaction..")
		_, _ = conn.Write(h.parser.WriteArray([][]byte{}))
		h.removeTransaction(conn)
		return
	}

	var responses [][]byte
	conn1, conn2 := net.Pipe()

	fmt.Println("Pipes created..")
	// fmt.Println("conn1: ", conn1)
	// fmt.Println("conn2: ", conn2)

	// conn1.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	for _, cmd := range commands {
		go h.Handle(conn2, cmd)
		buffer := make([]byte, 1024)
		n, err := conn1.Read(buffer)
		buffer = buffer[:n]

		if err != nil {
			responses = append(responses, h.parser.WriteError("Error reading from pipe"))
		}
		responses = append(responses, buffer)
	}
	fmt.Println("Responses: ", responses)

	_, _ = conn.Write(h.parser.WriteArray(responses))

}
