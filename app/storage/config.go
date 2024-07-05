package storage

import "strings"

type Config struct {
	Dir        string
	DBfilename string
	Role       string
	MasterAddr string
	MasterPort string
}

func DefaultConfig() *Config {
	return &Config{
		Dir:        "/tmp",
		DBfilename: "dump.rdb",
	}
}

func NewConfig(dbdir, dbfilename, replica string) *Config {
	var role string
	var masterAddr string
	var masterPort string

	if replica != "" {
		role = "slave"
		masterInfo := strings.Split(replica, " ")
		masterAddr = masterInfo[0]
		masterPort = masterInfo[1]
	} else {
		role = "master"
		masterAddr = ""
		masterPort = ""
	}

	return &Config{
		Dir:        dbdir,
		DBfilename: dbfilename,
		Role:       role,
		MasterAddr: masterAddr,
		MasterPort: masterPort,
	}
}
