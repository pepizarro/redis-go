package storage

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Entry struct {
	Id     string
	Values map[string][]byte
}

type Stream struct {
	entries []Entry
}

func (k *KeySpace) NewEntry(id string, values map[string][]byte) Entry {
	return Entry{
		Id:     id,
		Values: values,
	}
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

// func (k *KeySpace) GetStreamStartingAt(key string, start string) (Stream, error) {
//
//
//
// 	return Stream{}, nil
// }
//
// func (k *KeySpace) GetStreamEndingAt(key string, start string) (Stream, error) {
//
// 	return Stream{}, nil
// }

func (k *KeySpace) GetEntriesInRange(key string, start string, end string) ([]Entry, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	i, exists := k.keyspace[key]
	if !exists {
		return nil, fmt.Errorf("Key not found: %s", key)
	}
	if i.valueType != STREAM {
		return nil, fmt.Errorf("Invalid type: %s", i.valueType)
	}

	startMs, startSeq, err := parseEntryId(start)
	if err != nil {
		return nil, err
	}
	endMs, endSeq, err := parseEntryId(end)
	if err != nil {
		return nil, err
	}

	stream := i.value.(Stream)
	entries := stream.entries
	var newEntries []Entry
	for _, entry := range entries {
		ms, seq, err := parseEntryId(entry.Id)
		if err != nil {
			return nil, err
		}

		if ms >= startMs && ms <= endMs {
			if seq >= startSeq && seq <= endSeq {
				newEntries = append(newEntries, entry)
			}
		}
	}

	return newEntries, nil

}

func (k *KeySpace) SetEntryWithAutoGeneratedID(key string, entry Entry) (string, error) {
	fmt.Println("Setting entry with WithAutoGeneratedID")

	k.mu.Lock()
	defer k.mu.Unlock()
	i, exists := k.keyspace[key]
	if exists && i.valueType != STREAM {
		return "", fmt.Errorf("Key %s already exists and is not a stream", key)
	}

	if exists {
		// add entry to the stream
		lastEntry := i.value.(Stream).entries[len(i.value.(Stream).entries)-1]
		lastMs, lastSeq, err := parseEntryId(lastEntry.Id)
		if err != nil {
			return "", err
		}

		entryID := fmt.Sprintf("%d-%d", lastMs+1, lastSeq+1)
		entry.Id = entryID
		streamEntries := i.value.(Stream).entries
		streamEntries = append(streamEntries, entry)
		k.keyspace[key] = item{value: Stream{entries: streamEntries}, valueType: STREAM}
	}

	ms := time.Now().UnixMilli()
	entryID := fmt.Sprintf("%d-%d", ms, 0)
	entry.Id = entryID
	k.keyspace[key] = item{value: Stream{entries: []Entry{entry}}, valueType: STREAM}

	return entry.Id, nil
}

func (k *KeySpace) SetEntryWithAutoGeneratedSequence(key string, ms string, entry Entry) (string, error) {
	fmt.Println("Setting entry with WithAutoGeneratedSequence")

	k.mu.Lock()
	defer k.mu.Unlock()
	i, exists := k.keyspace[key]
	if exists && i.valueType != STREAM {
		return "", fmt.Errorf("Key %s already exists and is not a stream", key)
	}

	if exists {
		// add entry to the stream
		lastSeq, err := validateMilliseconds(i.value.(Stream), ms)
		if err != nil {
			return "", err
		}
		fmt.Println("Adding entry to stream")

		entryID := fmt.Sprintf("%s-%d", ms, lastSeq+1)
		entry.Id = entryID
		streamEntries := i.value.(Stream).entries
		streamEntries = append(streamEntries, entry)
		k.keyspace[key] = item{value: Stream{entries: streamEntries}, valueType: STREAM}

		return entryID, nil
	}

	entryID := fmt.Sprintf("%s-%d", ms, 1)
	entry.Id = entryID
	k.keyspace[key] = item{value: Stream{entries: []Entry{entry}}, valueType: STREAM}

	return entry.Id, nil
}

func (k *KeySpace) SetEntryWithID(key string, id string, entry Entry) (string, error) {
	fmt.Println("Setting entry with WithID")

	i, exists := k.keyspace[key]
	if exists && i.valueType != STREAM {
		return "", fmt.Errorf("Key %s already exists and is not a stream", key)
	}

	if exists {
		// add entry to the stream
		ms, lastSeq, err := validateId(i.value.(Stream), id)
		if err != nil {
			return "", err
		}
		fmt.Println("Adding entry to stream")

		entryID := fmt.Sprintf("%d-%d", ms, lastSeq+1)
		entry.Id = entryID
		streamEntries := i.value.(Stream).entries
		streamEntries = append(streamEntries, entry)
		k.keyspace[key] = item{value: Stream{entries: streamEntries}, valueType: STREAM}

		return entryID, nil
	}

	entry.Id = id
	k.keyspace[key] = item{value: Stream{entries: []Entry{entry}}, valueType: STREAM}

	return entry.Id, nil
}

func validateId(stream Stream, id string) (int, int, error) {

	// get the last entry
	lastEntry := stream.entries[len(stream.entries)-1]
	lastMs, lastSeq, err := parseEntryId(lastEntry.Id)
	if err != nil {
		return 0, 0, err
	}

	ms, seq, err := parseEntryId(id)
	if err != nil {
		return 0, 0, err
	}

	if ms < lastMs {
		return 0, 0, fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if ms == lastMs && seq <= lastSeq {
		return 0, 0, fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return ms, lastSeq, nil
}

func validateMilliseconds(stream Stream, ms string) (int, error) {

	// get the last entry
	lastEntry := stream.entries[len(stream.entries)-1]
	lastMs, lastSeq, err := parseEntryId(lastEntry.Id)
	if err != nil {
		return 0, err
	}

	msInt, err := strconv.Atoi(ms)
	if err != nil {
		return 0, err
	}

	if msInt < lastMs {
		return 0, fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if msInt == lastMs {
		return lastSeq, nil
	}

	return -1, nil
}

func parseEntryId(id string) (int, int, error) {

	r, _ := regexp.Compile(`(\d+)-(\d+)`)
	matches := r.FindStringSubmatch(id)
	if matches == nil {
		return 0, 0, fmt.Errorf("Invalid entry ID format")
	}

	ms, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, err
	}

	seq, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, err
	}

	return ms, seq, nil
}

func parseEntryIdWithRange(id string) (error, string, string) {

	r, _ := regexp.Compile(`(\d+)-(\d+|\*)`)
	matches := r.FindStringSubmatch(id)
	if matches == nil {
		return fmt.Errorf("Invalid entry ID format"), "", ""
	}

	ms := matches[1]
	seq := matches[2]

	return nil, ms, seq
}
