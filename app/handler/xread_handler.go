package handler

import (
	"fmt"
	"net"
)

func (h *Handler) XreadHandler(conn net.Conn, buffer []byte) {
	fmt.Println("XREAD Handler")
	h.store.LogOnce()
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

	keysAndIDs := params[6:]
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

	fmt.Println("Keys: ", keys)

	var response [][]byte
	h.store.LogOnce()

	for key, start := range keys {
		// 	fmt.Println("Key: ", key, "Start: ", start)
		//
		start = completeEntryID(string(start))
		//
		entries, err := h.store.GetEntriesAfter(key, start)
		fmt.Println("entries: ", entries)
		if err != nil {
			fmt.Println("Error getting stream ending at: ", err)
			_, err = conn.Write(h.parser.WriteError(err.Error()))
			if err != nil {
				fmt.Println("Error writing to client: ", err)
			}
			continue
		}
		//
		// 	fmt.Println("Received Entry: ", entries)
		//
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
		response = append(response, h.parser.WriteArray(streamArray))
	}

	fmt.Println("Response: \n", string(h.parser.WriteArray(response)))
	_, err = conn.Write(h.parser.WriteArray(response))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
	}

	return
}

// func (h *Handler) XreadHandler(conn net.Conn, buffer []byte) {
// 	fmt.Println("XREAD Handler")
//
// 	h.store.LogOnce()
//
// }
