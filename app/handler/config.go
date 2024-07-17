package handler

import (
	"fmt"
	"strings"
)

type HandlerConfig struct {
	address           string
	port              string
	role              string
	masterAddr        string
	masterPort        string
	replicationID     string
	replicationOffset int
}

func NewHandlerConfig(address, port, replicaInfo string) *HandlerConfig {
	newConfig := &HandlerConfig{}

	if replicaInfo != "" {
		newConfig.role = "slave"
		masterInfo := strings.Split(replicaInfo, " ")
		newConfig.masterAddr = masterInfo[0]
		newConfig.masterPort = masterInfo[1]
	} else {
		newConfig.role = "master"
		newConfig.masterAddr = ""
		newConfig.masterPort = ""
		newConfig.replicationID = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		newConfig.replicationOffset = 0
	}

	newConfig.address = address
	newConfig.port = port

	return newConfig
}

func (c *HandlerConfig) getReplicationInfo() map[string]string {

	replicationInfo := make(map[string]string)
	replicationInfo["role"] = c.role
	replicationInfo["master_replid"] = c.replicationID
	replicationInfo["master_repl_offset"] = fmt.Sprintf("%d", c.replicationOffset)

	return replicationInfo

}

func (c *HandlerConfig) IsReplica() bool {
	return c.role == "slave"
}

func (c *HandlerConfig) isMaster() bool {
	return c.role == "master"
}

func (c *HandlerConfig) getSocket() string {
	return c.address + ":" + c.port
}

func (c *HandlerConfig) getListeningPort() string {
	return c.port
}

func (c *HandlerConfig) getMasterSocket() string {
	return c.masterAddr + ":" + c.masterPort
}
