package handler

import (
	"fmt"
	"net"
)

func (h *Handler) XreadHandler(conn net.Conn, buffer []byte) {

	params, err := h.parser.GetParams(buffer)
	if err != nil {
		fmt.Println("Error getting params: ", err)
		return
	}

	// if len(params) < 6 {
	// 	fmt.Println("Invalid number of arguments")
	// 	return
	// }

	key := string(params[6])

	start := completeEntryID(string(params[8]))

	entries, err := h.store.GetEntriesAfter(key, start)
	fmt.Println("Entries: ", entries)
	if err != nil {
		fmt.Println("Error getting stream ending at: ", err)
		_, err = conn.Write(h.parser.WriteError(err.Error()))
		if err != nil {
			fmt.Println("Error writing to client: ", err)
		}
		return
	}

	var streamArray [][]byte
	streamArray = append(streamArray, h.parser.WriteString(key))
	var entriesArray [][]byte
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

	var response [][]byte
	response = append(response, h.parser.WriteArray(streamArray))

	fmt.Println("Response: \n", string(h.parser.WriteArray(response)))
	_, err = conn.Write(h.parser.WriteArray(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
	}
	return

}
