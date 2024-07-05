package storage

import (
	"fmt"
	"strings"
)

type Config struct {
	Dir               string
	DBfilename        string
	Role              string
	MasterAddr        string
	MasterPort        string
	ReplicationID     string
	ReplicationOffset int
}

func DefaultConfig() *Config {
	return &Config{
		Dir:        "/tmp",
		DBfilename: "dump.rdb",
	}
}

func NewConfig(dbdir, dbfilename, replica string) *Config {
	newConfig := &Config{
		Dir:        dbdir,
		DBfilename: dbfilename,
	}

	if replica != "" {
		newConfig.Role = "slave"
		masterInfo := strings.Split(replica, " ")
		newConfig.MasterAddr = masterInfo[0]
		newConfig.MasterPort = masterInfo[1]
	} else {
		newConfig.Role = "master"
		newConfig.MasterAddr = ""
		newConfig.MasterPort = ""
		newConfig.ReplicationID = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		newConfig.ReplicationOffset = 0
	}

	return newConfig
}

func (c *Config) GetReplicationInfo() map[string]string {

	replicationInfo := make(map[string]string)
	replicationInfo["role"] = c.Role
	if c.Role == "master" {
		replicationInfo["master_replid"] = c.ReplicationID
		replicationInfo["master_repl_offset"] = fmt.Sprintf("%d", c.ReplicationOffset)
	}

	return replicationInfo

}
