package storage

type StorageConfig struct {
	Dir        string
	DBfilename string
}

func DefaultConfig() *StorageConfig {
	return &StorageConfig{
		Dir:        "/tmp",
		DBfilename: "dump.rdb",
	}
}

func NewConfig(dbdir, dbfilename string) *StorageConfig {
	newConfig := &StorageConfig{
		Dir:        dbdir,
		DBfilename: dbfilename,
	}

	return newConfig
}

func (c *StorageConfig) GetReplicationInfo() map[string]string {

	replicationInfo := make(map[string]string)

	return replicationInfo

}
