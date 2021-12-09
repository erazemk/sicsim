package sim

import (
	"encoding/binary"
	"os"
)

// ReadString reads and returns a string of length len from file
func ReadString(file *os.File, len int) (string, error) {
	bytes := make([]byte, len)
	n, err := file.Read(bytes)

	if err != nil {
		return "", err
	}

	return string(bytes[:n]), nil
}

// ReadByte reads and returns a byte from file
func ReadByte(file *os.File) (byte, error) {
	bytes := make([]byte, 1)
	_, err := file.Read(bytes)

	if err != nil {
		return 0, err
	}

	return bytes[0], nil
}

// ReadWord reads and returns a word from file
func ReadWord(file *os.File) (int, error) {
	bytes := make([]byte, 3)
	_, err := file.Read(bytes)

	if err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint32(bytes)), nil
}
