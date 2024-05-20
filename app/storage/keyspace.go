package storage

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type KeySpace struct {
	mu       sync.RWMutex
	keyspace map[string]item
	config   *Config
}

type item struct {
	// mu         sync.RWMutex
	value      []byte
	expiration time.Time
}

func NewKeySpace(config *Config) *KeySpace {

	if config == nil {
		config = DefaultConfig()
	}

	fmt.Println("Using config: ", config)

	ks := &KeySpace{
		keyspace: make(map[string]item),
		config:   config,
	}

	go func() {
		for {
			ks.cleanup()
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// go ks.LogKeySpace()
	ks.LoadSnapshots()

	return ks
}

func (k *KeySpace) LoadSnapshots() {
	// Load snapshots from disk
	k.mu.Lock()
	defer k.mu.Unlock()

	rdbFile := k.config.Dir + "/" + k.config.DBfilename
	data, err := os.ReadFile(rdbFile)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}

	fmt.Println("Loading RDB file...")
	rdbStruct, err := newRdb(data)
	if err != nil {
		fmt.Println("Error reading rdb: ", err)
		return
	}

	k.keyspace = rdbStruct.KeyValues

	fmt.Println("RDB: ", rdbStruct)

}

func (k *KeySpace) cleanup() {
	k.mu.Lock()
	defer k.mu.Unlock()

	for key, item := range k.keyspace {
		if item.expiration != (time.Time{}) && time.Now().After(item.expiration) {
			delete(k.keyspace, key)
		}
	}
}

func (k *KeySpace) GetConfig() *Config {
	return k.config
}

func (k *KeySpace) LogKeySpace() {
	for {
		k.mu.RLock()
		fmt.Println("KeySpace: ")
		for key, item := range k.keyspace {
			fmt.Println(key, ":", string(item.value))
			if item.expiration != (time.Time{}) {
				fmt.Println("  Expiration: ", item.expiration)
			}

		}
		k.mu.RUnlock()

		time.Sleep(5 * time.Second)
	}
}

func (k *KeySpace) Set(key string, value []byte) {
	k.mu.Lock()
	k.keyspace[key] = item{value: value}
	k.mu.Unlock()
}

func (k *KeySpace) SetWithExpiration(key string, value []byte, expiration time.Time) {
	k.mu.Lock()
	k.keyspace[key] = item{value: value, expiration: expiration}
	k.mu.Unlock()
}

func (k *KeySpace) Get(key string) ([]byte, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if item, exist := k.keyspace[key]; exist {
		return item.value, nil
	}

	return nil, fmt.Errorf("Key not found: %s", key)
}

func (k *KeySpace) GetAllKeys() ([]string, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	var keys []string
	for k, _ := range k.keyspace {
		keys = append(keys, k)
	}

	return keys, nil

}

func (k *KeySpace) Delete(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	delete(k.keyspace, key)
}
