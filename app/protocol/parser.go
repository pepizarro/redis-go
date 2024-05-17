package protocol

type Parser interface {
	GetCommand(buffer []byte) (string, error)
	GetSubCommand(buffer []byte) (string, error)
	GetParams(buffer []byte) ([][]byte, error)
	GetArgs(buffer []byte, start int) (map[string][]byte, error)
	WriteString(string) []byte
	WriteOk() []byte
	WriteNull() []byte
	WriteArray([][]byte) []byte
}
