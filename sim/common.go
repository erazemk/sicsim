package sim

import (
	"fmt"
	"math"
	"os"
)

var debug bool

// Opcodes
const (
	ADD    byte = 0x18
	ADDF   byte = 0x58
	ADDR   byte = 0x90
	AND    byte = 0x40
	CLEAR  byte = 0xB4
	COMP   byte = 0x28
	COMPF  byte = 0x88
	COMPR  byte = 0xA0
	DIV    byte = 0x24
	DIVF   byte = 0x64
	DIVR   byte = 0x9C
	FIX    byte = 0xC4
	FLOAT  byte = 0xC0
	HIO    byte = 0xF4
	J      byte = 0x3C
	JEQ    byte = 0x30
	JGT    byte = 0x34
	JLT    byte = 0x38
	JSUB   byte = 0x48
	LDA    byte = 0x00
	LDB    byte = 0x68
	LDCH   byte = 0x50
	LDF    byte = 0x70
	LDL    byte = 0x08
	LDS    byte = 0x6C
	LDT    byte = 0x74
	LDX    byte = 0x04
	LPS    byte = 0xD0
	MUL    byte = 0x20
	MULF   byte = 0x60
	MULR   byte = 0x98
	NORM   byte = 0xC8
	OR     byte = 0x44
	RD     byte = 0xD8
	RMO    byte = 0xAC
	RSUB   byte = 0x4C
	SHIFTL byte = 0xA4
	SHIFTR byte = 0xA8
	SIO    byte = 0xF0
	SSK    byte = 0xEC
	STA    byte = 0x0C
	STB    byte = 0x78
	STCH   byte = 0x54
	STF    byte = 0x80
	STI    byte = 0xD4
	STL    byte = 0x14
	STS    byte = 0x7C
	STSW   byte = 0xE8
	STT    byte = 0x84
	STX    byte = 0x10
	SUB    byte = 0x1C
	SUBF   byte = 0x5C
	SUBR   byte = 0x94
	SVC    byte = 0xB0
	TD     byte = 0xE0
	TIO    byte = 0xF8
	TIX    byte = 0x2C
	TIXR   byte = 0xB8
	WD     byte = 0xDC
)

func init() {
	// Functions print logs if debug is true
	_, debug = os.LookupEnv("SICSIM_DEBUG")
}

func SetDebug(debugFlag bool) {
	debug = debugFlag
}

// isWord checks if val is a valid SIC word (24 bits)
func isWord(word int) bool {
	return word >= -int(math.Pow(2, 24)) && word < int(math.Pow(2, 24))
}

// isAddr checks if addr is a valid SIC address
func isAddr(addr int) bool {
	return addr >= 0 && addr <= MAX_ADDRESS
}

// isRegister checks if reg is a valid SIC register
func isRegister(reg int) bool {
	return (reg >= 0 && reg <= 6) || (reg == 8 || reg == 9)
}

func printWord(val int) string {
	return fmt.Sprintf("0x%06[1]X - %[1]d (%[1]b)", val)
}

func printByte(val byte) string {
	return fmt.Sprintf("0x%02[1]X - %[1]d (%[1]b)", val)
}
