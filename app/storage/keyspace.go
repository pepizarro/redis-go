package storage

import (
	"fmt"
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

	go ks.LogKeySpace()

	return ks
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

func (k *KeySpace) Delete(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	delete(k.keyspace, key)
}
