package handler

import (
	"fmt"
	"net"
	"time"
)

func (h *Handler) connectToMaster() {

	masterSocket := h.config.getMasterSocket()
	if masterSocket == "" {
		fmt.Println("No master socket provided")
		return
	}

	for {
		fmt.Println("Checking connection to master...")
		conn, err := net.Dial("tcp", masterSocket)

		if err != nil {
			fmt.Println("Error connecting to master: ", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		connectionMaster := h.connectionToMaster(conn)
		if connectionMaster != nil {
			fmt.Println("Error connecting to master: ", connectionMaster)
		} else {
			fmt.Println("Connected to master")
		}

		time.Sleep(2 * time.Second)
	}
}

func (h *Handler) connectionToMaster(conn net.Conn) error {
	fmt.Println("Setting connection to master: ", conn)
	listeningPort := h.config.getListeningPort()

	var messages [][]byte
	messages = append(messages, h.parser.WriteArray([][]byte{h.parser.WriteString("PING")}))
	messages = append(messages, h.parser.WriteArray([][]byte{h.parser.WriteString("REPLCONF"), h.parser.WriteString("listening-port"), h.parser.WriteString(listeningPort)}))
	messages = append(messages, h.parser.WriteArray([][]byte{h.parser.WriteString("REPLCONF"), h.parser.WriteString("capa"), h.parser.WriteString("psync2")}))
	messages = append(messages, h.parser.WriteArray([][]byte{h.parser.WriteString("PSYNC"), h.parser.WriteString("?"), h.parser.WriteString("-1")}))

	for _, message := range messages {
		err := h.sendAndRead(conn, message)
		if err != nil {
			return fmt.Errorf("Error connecting to master: %s", err)
		}
	}

	return h.listenFromMaster(conn)
}

func (h *Handler) listenFromMaster(conn net.Conn) error {

	buffer := make([]byte, 1024)
	for {
		// fmt.Println("Waiting for messages from master...")
		if conn == nil {
			return fmt.Errorf("Disconnected from master")
		}

		// conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		n, err := conn.Read(buffer)
		// h.parser
		if err != nil {
			continue
		}
		if n > 0 {
			fmt.Println("Received bytes from master: ", buffer[:n])
			_, err := h.parser.GetCommand(buffer[:n])
			if err == nil {
				// fmt.Println("Is a command: ", command)
				// parse every array
				commands, err := h.parser.GetArrays(buffer[:n])
				if err != nil {
					return fmt.Errorf("Error getting commands: %s", err)
				}
				fmt.Println("Received commands from master: ", commands)
				for _, c := range commands {
					h.Handle(conn, c)
				}
			}
		}
	}
}

func (h *Handler) sendAndRead(conn net.Conn, buffer []byte) error {
	fmt.Println("Sending in sendAndRead: ", string(buffer))
	_, err := conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("Error writing: %s", err)
	}
	buffer = make([]byte, 1024)
	fmt.Println("Reading...")
	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("Error reading: %s", err)
	}
	fmt.Println("Received in sendAndRead: ", string(buffer))

	return nil
}

func (h *Handler) send(conn net.Conn, buffer []byte) error {
	fmt.Println("Sending in sendAndRead: ", string(buffer))
	_, err := conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("Error writing: %s", err)
	}
	return nil
}

// Master methods

func (h *Handler) replicate(buffer []byte) error {
	// fmt.Println("Replicating")

	for _, replica := range h.replicas {
		fmt.Println("Replicating to replica: ", replica.RemoteAddr())
		_, err := replica.Write(buffer)
		if err != nil {
			fmt.Println("Error writing to replica")
		}
	}

	return nil
}
