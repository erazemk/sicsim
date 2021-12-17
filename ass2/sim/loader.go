package sicsim

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

// ParseByte parses and returns a byte from two characters
func ParseByte(text string) byte {
	bytes, err := hex.DecodeString(text)

	if err != nil {
		panic(err)
	}

	return bytes[0]
}

// ParseWord parses and returs a word from six characters
func ParseWord(text string) int {
	hex, err := hex.DecodeString(text)

	if err != nil {
		panic(err)
	}

	bytes := make([]byte, 1)
	bytes = append(bytes, hex[0], hex[1], hex[2])

	return int(binary.BigEndian.Uint32(bytes))
}

// ParseObj splits an object file into sections and calls appropriate functions
func (m *Machine) ParseObj(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("failed to parse object file: %w", err)
	}

	scanner := bufio.NewScanner(file)

	// Split object file, each line/section is a separate element
	scanner.Split(bufio.ScanLines)
	var sections []string

	for scanner.Scan() {
		sections = append(sections, scanner.Text())
	}

	file.Close()

	var lc int = 0 // Loader Counter - keeps track of current memory location

	if debug {
		fmt.Println("--- Start ParseObj ---")
	}

	// Parse each section
	for _, line := range sections {
		record := line[0:1]

		switch record {
		case "H": // Header
			progName := line[1:7]
			codeAddr := ParseWord(line[7:13])
			codeLen := ParseWord(line[13:19])
			lc = codeAddr // Start writing commands to memory[codeAddr]

			if err != nil {
				panic(err)
			}

			if debug {
				fmt.Println("[Header]")
				fmt.Println("    name: " + progName)
				fmt.Println("    addr: " + printWord(codeAddr))
				fmt.Println("    len: " + printWord(codeLen))
				fmt.Println("    lc: " + strconv.FormatInt(int64(lc), 10))
			}
		case "E": // End
			startAddr := ParseWord(line[1:7])

			if debug {
				fmt.Println("[End]")
				fmt.Println("    start addr: " + printWord(startAddr))
				fmt.Println("    lc: " + strconv.FormatInt(int64(lc), 10))
			}

			m.SetPC(startAddr) // Start executing commands at memory[startAddr]
		case "T": // Text (code)
			codeAddr := ParseWord(line[1:7])
			codeLen := ParseByte(line[7:9])
			code := line[9:]

			if debug {
				fmt.Println("[Text]")
				fmt.Println("    addr: " + printWord(codeAddr))
				fmt.Println("    len: " + printByte(codeLen))
				fmt.Println("    code: " + code)
				fmt.Println("    lc: " + strconv.FormatInt(int64(lc), 10))
			}

			for i := 0; i < int(codeLen)*2; i += 2 {
				bytes := ParseByte(code[i : i+2])

				if debug {
					fmt.Printf("        byte: %02X\n", bytes)
				}

				m.SetByte(lc, bytes)
				lc++
			}
		case "M": // Modification
			offset := line[1:7]
			lineLen := line[7:9]
			longVer := len(line) > 9
			var operator, symbolName string

			if longVer {
				operator = line[9:10]
				symbolName = line[10:16]
			}

			if debug {
				fmt.Println("[Modification]")
				fmt.Println("    offset: " + offset)
				fmt.Println("    len: " + lineLen)
				fmt.Println("    lc: " + strconv.FormatInt(int64(lc), 10))

				if longVer {
					fmt.Println("    operator: " + operator)
					fmt.Println("    symbol name: " + symbolName)
				}
			}
		case "D": // Exported symbol
			name := line[1:7]
			value := line[7:13]
			pairs := line[13:]

			if debug {
				fmt.Println("[Symbol export]")
				fmt.Println("    name: " + name)
				fmt.Println("    value: " + value)
				fmt.Println("    pairs: " + pairs)
				fmt.Println("    lc: " + strconv.FormatInt(int64(lc), 10))
			}
		case "R": // Imported symbol
			name := line[1:7]
			otherNames := line[7:]

			if debug {
				fmt.Println("[Symbol import]")
				fmt.Println("    name: " + name)
				fmt.Println("    other names: " + otherNames)
				fmt.Println("    lc: " + strconv.FormatInt(int64(lc), 10))
			}
		}
	}

	if debug {
		fmt.Println("--- End ParseObj ---")
	}

	return nil
}
