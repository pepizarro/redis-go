package storage

import "fmt"

type Entry struct {
	id     string
	values map[string][]byte
}

type Stream struct {
	entries []Entry
}

func (k *KeySpace) NewEntry(id string, values map[string][]byte) Entry {
	return Entry{
		id:     id,
		values: values,
	}
}

func (k *KeySpace) SetStream(key string, entry Entry) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	i, exists := k.keyspace[key]
	if exists && i.valueType != STREAM {
		return fmt.Errorf("Key %s already exists and is not a stream", key)
	}

	if exists {
		// add entry to the stream
		fmt.Println("Adding entry to stream")
	}

	stream := Stream{
		entries: []Entry{entry},
	}

	k.keyspace[key] = item{value: stream, valueType: STREAM}

	return nil
}
