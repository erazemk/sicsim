package asm

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var Opcodes map[string]byte
var debug, prettyPrint bool

// Mnemonics
var Directive = []string{"NOBASE", "LTORG"}
var DirectiveN = []string{"START", "END", "BASE", "ORG"}
var StorageDirective = []string{"BYTE", "WORD"}
var StorageDirectiveN = []string{"RESB", "RESW"}

var InstructionF1 = []string{"FIX", "FLOAT", "HIO", "NORM", "SIO", "TIO"}
var InstructionF2n = []string{"SVC"}
var InstructionF2r = []string{"CLEAR", "TIXR"}
var InstructionF2rn = []string{"SHIFTL", "SHIFTR"}
var InstructionF2rr = []string{"ADDR", "SUBR", "DIVR", "MULR", "COMPR", "RMO"}
var InstructionF3 = []string{"RSUB", "STI"}
var InstructionF3m = []string{
	"LDA", "LDB", "LDCH", "LDF", "LDL", "LDS", "LDT", "LDX", // Load
	"STA", "STB", "STCH", "STF", "STL", "STS", "STT", "STX", // Store
	"ADD", "AND", "COMP", "DIV", "SUB", "MUL", "OR", // Math
	"ADDF", "COMPF", "DIVF", "SUBF", "MULF", // Float
	"J", "JEQ", "JGT", "GLT", "JSUB", // Jump
	"RD", "TD", "WD", // Device I/O
	"SSK", // Other
}

var Directives []string
var StorageDirectives []string
var InstructionsF2 []string
var InstructionsF3 []string
var Instructions []string
var Mnemonics []string

func init() {
	Directives = append(Directives, Directive...)
	Directives = append(Directives, DirectiveN...)

	StorageDirectives = append(StorageDirectives, StorageDirective...)
	StorageDirectives = append(StorageDirectives, StorageDirectiveN...)

	InstructionsF2 = append(InstructionsF2, InstructionF2n...)
	InstructionsF2 = append(InstructionsF2, InstructionF2r...)
	InstructionsF2 = append(InstructionsF2, InstructionF2rn...)
	InstructionsF2 = append(InstructionsF2, InstructionF2rr...)

	InstructionsF3 = append(InstructionsF3, InstructionF3...)
	InstructionsF3 = append(InstructionsF3, InstructionF3m...)

	Instructions = append(Instructions, InstructionF1...)
	Instructions = append(Instructions, InstructionsF2...)
	Instructions = append(Instructions, InstructionsF3...)

	Mnemonics = append(Mnemonics, Directives...)
	Mnemonics = append(Mnemonics, StorageDirectives...)
	Mnemonics = append(Mnemonics, Instructions...)

	Opcodes = map[string]byte{
		"ADD":    0x18,
		"ADDF":   0x58,
		"ADDR":   0x90,
		"AND":    0x40,
		"CLEAR":  0xB4,
		"COMP":   0x28,
		"COMPF":  0x88,
		"COMPR":  0xA0,
		"DIV":    0x24,
		"DIVF":   0x64,
		"DIVR":   0x9C,
		"FIX":    0xC4,
		"FLOAT":  0xC0,
		"HIO":    0xF4,
		"J":      0x3C,
		"JEQ":    0x30,
		"JGT":    0x34,
		"JLT":    0x38,
		"JSUB":   0x48,
		"LDA":    0x00,
		"LDB":    0x68,
		"LDCH":   0x50,
		"LDF":    0x70,
		"LDL":    0x08,
		"LDS":    0x6C,
		"LDT":    0x74,
		"LDX":    0x04,
		"LPS":    0xD0,
		"MUL":    0x20,
		"MULF":   0x60,
		"MULR":   0x98,
		"NORM":   0xC8,
		"OR":     0x44,
		"RD":     0xD8,
		"RMO":    0xAC,
		"RSUB":   0x4C,
		"SHIFTL": 0xA4,
		"SHIFTR": 0xA8,
		"SIO":    0xF0,
		"SSK":    0xEC,
		"STA":    0x0C,
		"STB":    0x78,
		"STCH":   0x54,
		"STF":    0x80,
		"STI":    0xD4,
		"STL":    0x14,
		"STS":    0x7C,
		"STSW":   0xE8,
		"STT":    0x84,
		"STX":    0x10,
		"SUB":    0x1C,
		"SUBF":   0x5C,
		"SUBR":   0x94,
		"SVC":    0xB0,
		"TD":     0xE0,
		"TIO":    0xF8,
		"TIX":    0x2C,
		"TIXR":   0xB8,
		"WD":     0xDC,
	}
}

func SetDebug(debugFlag bool) {
	debug = debugFlag
}

func SetPrettyPrint(lstFlag bool) {
	prettyPrint = lstFlag
}

func inSlice(elem string, slice []string) bool {
	for _, el := range slice {
		if el == elem {
			return true
		}
	}

	return false
}

func isNumber(number string) (int, bool) {
	var num int64
	var err error

	if strings.HasPrefix(number, "X'") { // Hex format
		num, err = strconv.ParseInt(number[2:len(number)-1], 16, 8)
		if err != nil {
			panic(fmt.Errorf("not a number '%s': %w", number, err))
		}
	} else { // Either number or variable
		num, err = strconv.ParseInt(number, 10, 8)
		if err != nil {
			return 0, false
		}
	}

	return int(num), true
}

func Word(number int) string {
	if number < 0 { // twos compliment negative operation
		number = int(math.Pow(2, 24)) + number
	}

	return fmt.Sprintf("%06X", number)
}
