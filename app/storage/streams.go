package storage

import (
	"fmt"
	"regexp"
	"strconv"
)

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
		err := k.validateEntryID(i.value.(Stream), entry.id)
		if err != nil {
			return err
		}
		fmt.Println("Adding entry to stream")
	}

	err, _, _ := parseEntryId(entry.id)
	if err != nil {
		return err
	}

	stream := Stream{
		entries: []Entry{entry},
	}

	k.keyspace[key] = item{value: stream, valueType: STREAM}

	return nil
}

func (k *KeySpace) GetStream(key string) (Stream, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	i, exists := k.keyspace[key]
	if !exists {
		return Stream{}, fmt.Errorf("Key not found: %s", key)
	}

	if i.valueType != STREAM {
		return Stream{}, fmt.Errorf("Invalid type: %s", i.valueType)
	}

	return i.value.(Stream), nil
}

func (k *KeySpace) validateEntryID(stream Stream, id string) error {
	// check id format n-n
	err, milliseconds, sequence := parseEntryId(id)
	if err != nil {
		return err
	}

	if milliseconds == 0 && sequence == 0 {
		return fmt.Errorf("ERR The ID specified in XADD must be greater than 0-0")
	}

	// compare with last entry
	if len(stream.entries) == 0 {
		return nil
	}

	lastEntry := stream.entries[len(stream.entries)-1]
	err, lastMilliseconds, lastSequence := parseEntryId(lastEntry.id)
	if err != nil {
		return err
	}

	if milliseconds < lastMilliseconds {
		return fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if milliseconds == lastMilliseconds && sequence <= lastSequence {
		return fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return nil
}

func parseEntryId(id string) (error, int, int) {
	r, _ := regexp.Compile(`(\d+)-(\d+)`)
	matches := r.FindStringSubmatch(id)

	if matches == nil {
		return fmt.Errorf("Invalid entry ID format"), 0, 0
	}

	milliseconds, _ := strconv.Atoi(matches[1])
	sequence, _ := strconv.Atoi(matches[2])

	return nil, milliseconds, sequence
}
