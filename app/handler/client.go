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

		connectionMaster := h.SetConnectionToMaster(conn)
		if connectionMaster != nil {
			fmt.Println("Error connecting to master: ", connectionMaster)
		} else {
			fmt.Println("Connected to master")
		}

		time.Sleep(5 * time.Second)
	}
}

func (h *Handler) SetConnectionToMaster(conn net.Conn) error {
	fmt.Println("Setting connection to master: ", conn)
	listeningPort := h.config.getListeningPort()

	buffer := make([]byte, 1024)

	buffer = h.parser.WriteArray([][]byte{h.parser.WriteString("PING")})
	_, err := conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("Error writing to master: %s", err)
	}
	buffer = make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("Error reading from master: %s", err)
	}

	buffer = h.parser.WriteArray([][]byte{h.parser.WriteString("REPLCONF"), h.parser.WriteString("listening-port"), h.parser.WriteString(listeningPort)})
	_, err = conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("Error writing to master: %s", err)
	}
	fmt.Println("Sent REPLCONF listening-port to master ", listeningPort)
	buffer = make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("Error reading from master: %s", err)
	}

	buffer = h.parser.WriteArray([][]byte{h.parser.WriteString("REPLCONF"), h.parser.WriteString("capa"), h.parser.WriteString("psync2")})
	_, err = conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("Error writing to master: %s", err)
	}
	fmt.Println("Sent REPLCONF listening-port to master ", listeningPort)
	buffer = make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("Error reading from master: %s", err)
	}

	if conn != nil {
		defer conn.Close()
	}

	return nil
}
