package sim

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func parseString(r *bufio.Reader, len int) string {
	buf := make([]rune, len)

	for i := 0; i < len; i++ {
		char, _, _ := r.ReadRune()
		buf[i] = char
	}

	if debug {
		fmt.Printf("String: %q\n", string(buf))
	}

	return string(buf)
}

func parseRune(r *bufio.Reader) rune {
	char, _, _ := r.ReadRune()

	if debug {
		fmt.Printf("Char: %q\n", char)
	}

	return char
}

func parseWord(r *bufio.Reader) int {
	buf := make([]rune, 6)

	for i := 0; i < 6; i++ {
		char, _, _ := r.ReadRune()
		buf[i] = char
	}

	word, err := strconv.ParseInt(string(buf), 16, 32)
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Printf("Word (before): %q\n", string(buf))
		fmt.Printf("Word (after): %06X\n", word)
	}

	return int(word)
}

func parseByte(r *bufio.Reader) byte {
	buf := make([]rune, 2)

	for i := 0; i < 2; i++ {
		char, _, _ := r.ReadRune()
		buf[i] = char
	}

	bytes, err := strconv.ParseInt(string(buf), 16, 32)
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Printf("Byte (before): %q\n", string(buf))
		fmt.Printf("Byte (after): %02X\n", bytes)
	}

	return byte(bytes)
}

func (m *Machine) ParseObjFile(objFile string) error {
	file, err := os.Open(objFile)
	if err != nil {
		return fmt.Errorf("failed to parse object file: %w", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	if debug {
		fmt.Println("--- Start ParseObj ---")
	}

	rec := parseRune(reader)

	// Header record
	if rec != 'H' {
		return fmt.Errorf("failed to parse object file: no header record")
	}

	progName := parseString(reader, 6)
	startAddr := parseWord(reader)
	codeLen := parseWord(reader)

	if debug {
		fmt.Println("[Header]")
		fmt.Println("    name: " + progName)
		fmt.Println("    addr: " + printWord(startAddr))
		fmt.Println("    len: " + printWord(codeLen))
	}

	// Seek to a new line and parse the record type
	reader.ReadLine()
	rec = parseRune(reader)

	// Text records
	for rec == 'T' {
		addr := parseWord(reader)
		len := parseByte(reader)

		if debug {
			fmt.Println("[Text]")
			fmt.Println("    addr: " + printWord(addr))
			fmt.Println("    len: " + printByte(len))
		}

		for i := 0; i < int(len); i++ {
			val := parseByte(reader)
			m.SetByte(addr, val)
			addr++
		}

		reader.ReadLine()
		rec = parseRune(reader)
	}

	// Modification records
	for rec == 'M' {
		offset := parseWord(reader)
		len := parseByte(reader)

		// TODO: Implement reading long version
		// operator := ParseStringReader(reader, 1)
		// name := ParseStringReader(reader, 6)

		if debug {
			fmt.Println("[Modification]")
			fmt.Println("    offset: " + printWord(offset))
			fmt.Println("    len: " + printByte(len))

			// if long {
			// 	fmt.Println("    operator: " + operator)
			// 	fmt.Println("    symbol name: " + name)
			// }
		}

		reader.ReadLine()
		rec = parseRune(reader)
	}

	// End record
	if rec != 'E' {
		return fmt.Errorf("failed to parse object file: no end record")
	}

	m.SetPC(parseWord(reader))

	if debug {
		fmt.Println("--- End ParseObj ---")
	}

	return nil
}
