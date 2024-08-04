package handler

import (
	"fmt"
	"net"
)

func (h *Handler) MultiHandler(conn net.Conn, buffer []byte) {
	fmt.Println("MultiHandler..")
	fmt.Println("conn: ", conn)
	fmt.Println("&conn: ", &conn)

	h.appendTransaction(conn)

	fmt.Println("Transaction appended..")
	for _, t := range h.transactions {
		fmt.Println("Transaction: ", t.conn)
	}

	_, _ = conn.Write(h.parser.WriteOk())
}

func (h *Handler) ExecHandler(conn net.Conn, buffer []byte) {
	fmt.Println("ExecHandler..")

	var transaction *Transaction
	for _, t := range h.transactions {
		if t.conn == conn {
			fmt.Println("Found transaction: ", t)
			transaction = t
			break
		}
	}

	fmt.Println("Transaction: ", transaction)

	if transaction == nil {
		_, _ = conn.Write(h.parser.WriteError("ERR EXEC without MULTI"))
		return
	}

	commands := transaction.commands
	if len(commands) == 0 {
		_, _ = conn.Write(h.parser.WriteArray([][]byte{}))
		h.removeTransaction(conn)
		return
	}

}
