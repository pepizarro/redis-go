package handler

import (
	// "encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func (h *Handler) XreadHandler(conn net.Conn, buffer []byte) {
	fmt.Println("XREAD Handler")
	// h.store.LogOnce()
	//
	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}
	if len(params) < 6 {
		fmt.Println("Invalid number of arguments")
		return
	}

	var keysAndIDs [][]byte
	var block bool
	var timeout int // in ms

	// fmt.Println("Params: ", params)

	if strings.ToLower(string(params[4])) == "block" {
		block = true
		// timeout = int(binary.BigEndian.Uint32(params[5]))
		strTimeout := string(params[6])
		timeout, err = strconv.Atoi(strTimeout)
		if err != nil {
			fmt.Println("Error converting timeout to int: ", err)
			return
		}
		keysAndIDs = params[10:]
		fmt.Println("params[10:]", params[10:])
	} else if strings.ToLower(string(params[4])) == "streams" {
		block = false
		keysAndIDs = params[6:]
	}
	// fmt.Println("block: ", block)
	// fmt.Println("timeout: ", timeout)
	// fmt.Println("keysAndIDs: ", len(keysAndIDs))

	if len(keysAndIDs)%4 != 0 {
		fmt.Println("Invalid number of arguments")
		_, _ = conn.Write(h.parser.WriteError("Invalid number of arguments"))
		return
	}

	keys := make(map[string]string)

	for i := 0; i < len(keysAndIDs)/2; i += 2 {
		key := string(keysAndIDs[i])
		keyID := string(keysAndIDs[i+len(keysAndIDs)/2])
		keys[key] = keyID
	}

	var response [][]byte

	startTime := time.Now()

	lastIDs := make(map[string]string)
	for key, start := range keys {
		if start == "$" {
			lastID, err := h.store.GetLastEntryID(key)
			if err != nil {
				_, err = conn.Write(h.parser.WriteNull())
				continue
			}
			lastIDs[key] = lastID
		}
	}

	for {
		for key, start := range keys {
			start = completeEntryID(string(start))

			var entries []storage.Entry

			if start == "$" {
				start = lastIDs[key]
				entries, err = h.store.GetEntriesAfter(key, start)
			} else {
				entries, err = h.store.GetEntriesAfter(key, start)
			}

			if err != nil {
				if block {
					continue
				}
				_, err = conn.Write(h.parser.WriteError(err.Error()))
				if err != nil {
					continue
				}
				continue
			}

			if len(entries) == 0 {
				continue
			}

			var streamArray [][]byte
			var entriesArray [][]byte

			streamArray = append(streamArray, h.parser.WriteString(key))
			for _, entry := range entries {
				var entryArray [][]byte
				entryArray = append(entryArray, h.parser.WriteString(entry.Id))
				var valuesArray [][]byte
				for k, v := range entry.Values {
					valuesArray = append(valuesArray, h.parser.WriteString(k))
					valuesArray = append(valuesArray, h.parser.WriteString(string(v)))
				}
				entryArray = append(entryArray, h.parser.WriteArray(valuesArray))
				entriesArray = append(entriesArray, h.parser.WriteArray(entryArray))
			}
			streamArray = append(streamArray, h.parser.WriteArray(entriesArray))
			response = append(response, h.parser.WriteArray(streamArray))
		}
		// fmt.Println("Response: \n", response)
		if !block {
			fmt.Println("Not blocking")
			break
		} else {
			if response != nil {
				break
			} else {
				if time.Since(startTime).Milliseconds() > int64(timeout) && timeout != 0 {
					_, _ = conn.Write(h.parser.WriteNull())
					return
				}
			}
		}
	}

	fmt.Println("Response: \n", response)
	if response == nil {
		_, _ = conn.Write(h.parser.WriteNull())
		return
	}
	_, err = conn.Write(h.parser.WriteArray(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
	}

	return
}
