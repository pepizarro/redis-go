package protocol

type Parser interface {
	GetCommand(buffer []byte) (string, error)
	GetSubCommand(buffer []byte) (string, error)
	GetParams(buffer []byte) ([][]byte, error)
	GetArgs(buffer []byte, start int) (map[string][]byte, error)
	GetType(buffer []byte) (string, error)
	GetArrays(buffer []byte) ([][]byte, error)
	WriteSimpleString(string) []byte
	WriteString(string) []byte
	WriteOk() []byte
	WriteError(string) []byte
	WriteNull() []byte
	WriteArray([][]byte) []byte
	WriteFile([]byte) []byte
	WriteInteger(int) []byte
}
