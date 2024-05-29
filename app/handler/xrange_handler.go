package handler

import (
	"fmt"
	"net"
	"regexp"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func (h *Handler) XrangeHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}

	if len(params) < 9 {
		fmt.Println("Invalid number of arguments")
		return
	}

	key := string(params[4])

	start := completeEntryID(string(params[6]))
	end := completeEntryID(string(params[8]))

	var entries []storage.Entry
	if start == "-" {
		entries, err = h.store.GetEntriesEndingAt(key, end)
	} else if end == "+" {
		entries, err = h.store.GetEntriesStartingAt(key, start)
	} else if start == "-" && end == "+" {
		entries, err = h.store.GetAllEntries(key)
	} else {
		entries, err = h.store.GetEntriesInRange(key, start, end)
	}

	if err != nil {
		fmt.Println("Error getting stream ending at: ", err)
		_, err = conn.Write(h.parser.WriteError(err.Error()))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
		}
		return
	}

	var response [][]byte
	for _, entry := range entries {
		var entryArray [][]byte
		entryArray = append(entryArray, h.parser.WriteString(entry.Id))
		var valuesArray [][]byte
		for k, v := range entry.Values {
			valuesArray = append(valuesArray, h.parser.WriteString(k))
			valuesArray = append(valuesArray, h.parser.WriteString(string(v)))
		}
		entryArray = append(entryArray, h.parser.WriteArray(valuesArray))
		response = append(response, h.parser.WriteArray(entryArray))
	}

	_, err = conn.Write(h.parser.WriteArray(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
	}
	return

}

func completeEntryID(id string) string {

	onlyMs := regexp.MustCompile(`^\d+$`)

	if onlyMs.MatchString(id) {
		return id + "-0"
	}

	return id
}
