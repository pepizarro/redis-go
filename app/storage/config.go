package storage

type Config struct {
	Dir        string
	DBfilename string
}

func DefaultConfig() *Config {
	return &Config{
		Dir:        "/tmp",
		DBfilename: "dump.rdb",
	}
}

func NewConfig(dbdir, dbfilename string) *Config {
	return &Config{
		Dir:        dbdir,
		DBfilename: dbfilename,
	}
}
