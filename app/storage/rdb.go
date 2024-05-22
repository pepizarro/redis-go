package storage

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	EOF            = 0xFF
	SELECT_DB      = 0xFE
	RESIZE_DB      = 0xFB
	EXPIRE_TIME    = 0xFD
	EXPIRE_TIME_MS = 0xFC
	AUX            = 0xFA
	MAGIC_NUMBER   = "REDIS"
	// RDB Types
	STRING_ENCODING              = 0
	LIST_ENCODING                = 1
	SET_ENCODING                 = 2
	SET_SORTED_ENCODING          = 3
	HASH_ENCODING                = 4
	ZIP_MAP_ENCODING             = 9
	ZIP_LIST_ENCODING            = 10
	INTSET_ENCODING              = 11
	SORTED_SET_ZIP_LIST_ENCODING = 12
	HASH_ZIP_LIST_ENCODING       = 13
	LIST_QUICKLIST_LIST_ENCODING = 14
)

type rdb struct {
	version         string
	auxiliaryFields map[string][]byte
	KeyValues       map[string]item
}

func isRDB() bool {
	return true
}

func newRdb(data []byte) (*rdb, error) {

	if len(data) < 9 {
		return nil, fmt.Errorf("Invalid RDB file")
	}

	newRDB := &rdb{
		version:         "",
		auxiliaryFields: make(map[string][]byte),
		KeyValues:       make(map[string]item),
	}

	// get redis magic string
	if string(data[:5]) != "REDIS" {
		return nil, fmt.Errorf("Invalid RDB file")
	}

	// get version
	newRDB.version = string(data[5:9])
	buf := data[9:]

	opCodes := map[byte]func(*rdb, []byte) int{
		AUX: readAux,
		// SELECT_DB:      readSelectDB,
		RESIZE_DB: readResizeDB,
		// EXPIRE_TIME:    readExpireTime,
		// EXPIRE_TIME_MS: readExpireTimeMS,
	}

	for i := 0; i < len(buf); i++ {
		b := buf[i]
		fmt.Printf("buf[%d]: %x\n", i, buf[i])
		if b == EOF {
			fmt.Println("EOF")
			break
		}

		if op, ok := opCodes[b]; ok {
			i += op(newRDB, buf[i+1:])
			fmt.Println("i: ", i)
		}
	}

	return newRDB, nil
}

// readLengthEncoded reads the ammount of bytes that a is using
// it then returns if it's a special format, the length of the value,
// and the slice of bytes from where the value starts
func readLengthEncoded(data []byte) (bool, int, []byte) {
	firstByte := data[0]
	twoMSB := firstByte & 0b11000000

	switch twoMSB {
	// case MSB == 00
	case 0:

		return false, int(firstByte), data[1:]

	// case MSB == 01
	case 64:
		if len(data) < 2 {
			return false, 0, data
		}
		bit16 := uint16(firstByte)<<8 | uint16(data[1])
		return false, int(bit16 & 0b0011111111111111), data[2:]

	// case MSB == 10
	case 128:
		if len(data) < 5 {
			return false, 0, data
		}
		num := int32(binary.BigEndian.Uint32(data[1:5]))
		return false, int(num), data[5:]

	// case MSB == 11
	case 192:
		return true, int(firstByte & 0b00111111), data[1:]
	}

	return false, 0, data
}

// readRedisString reads a string from the data using Length Encoding
// it then returns the value as a slice of bytes and the ammount of bytes read
func readRedisString(data []byte) ([]byte, []byte) {
	isSpecialFormat, length, buf := readLengthEncoded(data)
	if !isSpecialFormat {
		if len(buf) < length {
			return nil, nil
		}
		return buf[:length], buf[length:]
	}
	return nil, nil
}

func readAux(file *rdb, data []byte) int {
	fmt.Println("\n\nReading AUX")

	buf := data

	key, buf := readRedisString(buf)
	value, buf := readRedisString(buf)

	fmt.Println("Key: ", string(key), "value: ")
	file.auxiliaryFields[string(key)] = value

	if len(buf) == 0 {
		return 0
	}
	return len(data) - len(buf)
}

func readResizeDB(file *rdb, data []byte) int {
	fmt.Println("\n\nReading RESIZE_DB")
	buf := data[2:]
	ammount := int(data[0])
	// the first byte is the ammount of keys

	// isSpecialFormat, length, buf := readLengthEncoded(data)

	i := readKeyValue(ammount, file, buf)

	return i
}

func readKeyValue(ammount int, file *rdb, data []byte) int {

	buf := data
	for i := 0; i < ammount; i++ {
		fmt.Println("Buf in readKeyValue")
		printHex(buf)
		firstByte := buf[0]
		expireTime := time.Time{}
		if firstByte == EXPIRE_TIME {
			expireTime, buf = readExpireTime(buf[1:])
		} else if firstByte == EXPIRE_TIME_MS {
			expireTime, buf = readExpireTimeMs(buf[1:])
		}

		fmt.Println("Buf in readKeyValue AFTER")
		printHex(buf)
		valueType := buf[0]
		if valueType != STRING_ENCODING {
			fmt.Println("Not a string encoding")
			return 0
		}

		key, newBuf := readRedisString(buf[1:])

		// read value
		value, secondBuf := readRedisString(newBuf)

		buf = secondBuf
		fmt.Println("Key: ", key, "\nValue: ", value)
		file.KeyValues[string(key)] = item{value: value, expiration: expireTime}
	}

	return len(data) - len(buf)
}

func readExpireTime(data []byte) (time.Time, []byte) {
	// read the expire time
	timeBytes := data[:5]

	// convert the bytes to a time.Time
	timeInt := binary.LittleEndian.Uint32(timeBytes)
	expireTime := time.Unix(int64(timeInt), 0)

	return expireTime, data[4:]
}

func readExpireTimeMs(data []byte) (time.Time, []byte) {
	// read the expire time in milliseconds
	timeBytes := data[:9]

	timeInt := binary.LittleEndian.Uint64(timeBytes)
	expireTime := time.UnixMilli(int64(timeInt))

	return expireTime, data[8:]
}

func printHex(data []byte) {
	for _, b := range data {
		fmt.Printf("%x ", b)
	}
	fmt.Println()
	fmt.Println()
}
