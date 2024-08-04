package handler

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func (h *Handler) WaitHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		return
	}

	numReplicas, err := strconv.Atoi(string(params[4]))
	if err != nil {
		return
	}
	timeout, err := strconv.Atoi(string(params[6]))
	timeoutDuration := time.Duration(timeout) * time.Millisecond
	if err != nil {
		return
	}

	for _, replica := range h.replicas {

		err := h.send(replica.conn, h.parser.WriteArray([][]byte{h.parser.WriteString("REPLCONF"), h.parser.WriteString("GETACK"), h.parser.WriteString("*")}))
		if err != nil {
			fmt.Println("Error sending REPLCONF GETACK: ", err)
			continue
		}
	}

	replicasUpdated := 0
	if replicasUpdated < numReplicas {
		now := time.Now()
		for {
			replicasUpdated = 0
			for _, replica := range h.replicas {
				replica.mu.Lock()
				if replica.updated {
					replicasUpdated++
				}
				replica.mu.Unlock()
			}
			if replicasUpdated >= numReplicas {
				break
			}
			if time.Since(now) > timeoutDuration {
				fmt.Println("Timeout waiting for replicas to update..")
				break
			}
		}
	}

	response := h.parser.WriteInteger(replicasUpdated)
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}
