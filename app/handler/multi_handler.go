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

	commands := transaction.commands

	if len(commands) == 0 {
		fmt.Println("No commands in transaction..")
		_, _ = conn.Write(h.parser.WriteArray([][]byte{}))
		h.removeTransaction(conn)
		return
	}

	var responses [][]byte
	conn1, conn2 := net.Pipe()

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

	_, _ = conn.Write(h.parser.WriteArray(responses))

}

func (h *Handler) DiscardHandler(conn net.Conn) {
	fmt.Println("DiscardHandler..")

	for _, t := range h.transactions {
		if t.conn == conn {
			h.removeTransaction(conn)
			_, _ = conn.Write(h.parser.WriteOk())
			return
		}
	}

	_, _ = conn.Write(h.parser.WriteError("ERR DISCARD without MULTI"))
}
