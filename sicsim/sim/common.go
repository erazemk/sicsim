package sim

import (
	"math"
	"os"
)

var debug bool

// Opcodes
const (
	ADD    int = 0x18
	ADDF   int = 0x58
	AND    int = 0x40
	CLEAR  int = 0xB4
	COMP   int = 0x28
	COMPF  int = 0x88
	COMPR  int = 0xA0
	DIV    int = 0x24
	DIVF   int = 0x64
	DIVR   int = 0x9C
	FIX    int = 0xC4
	FLOAT  int = 0xC0
	HIO    int = 0xF4
	J      int = 0x3C
	JEQ    int = 0x30
	JGT    int = 0x34
	JLT    int = 0x38
	JSUB   int = 0x48
	LDA    int = 0x00
	LDB    int = 0x68
	LDCH   int = 0x50
	LDF    int = 0x70
	LDL    int = 0x08
	LDS    int = 0x6C
	LDT    int = 0x74
	LDX    int = 0x04
	LPS    int = 0xD0
	MUL    int = 0x20
	MULF   int = 0x60
	MULR   int = 0x98
	NORM   int = 0xC8
	OR     int = 0x44
	RD     int = 0xD8
	RMO    int = 0xAC
	RSUB   int = 0x4C
	SHIFTL int = 0xA4
	SHIFTR int = 0xA8
	SIO    int = 0xF0
	SSK    int = 0xEC
	STA    int = 0x0C
	STB    int = 0x78
	STCH   int = 0x54
	STF    int = 0x80
	STI    int = 0xD4
	STL    int = 0x14
	STS    int = 0x7C
	STSW   int = 0xE8
	STT    int = 0x84
	STX    int = 0x10
	SUB    int = 0x1C
	SUBF   int = 0x5C
	SUBR   int = 0x94
	SVC    int = 0xB0
	TD     int = 0xE0
	TIO    int = 0xF8
	TIX    int = 0x2C
	TIXR   int = 0xB8
	WD     int = 0xDC
)

func init() {
	// Functions print logs if debug is true
	_, debug = os.LookupEnv("SICSIM_DEBUG")
}

// isWord checks if val is a valid SIC word (24 bits)
func isWord(val int) bool {
	return val > 0 && float64(val) < math.Pow(2, 24)
}

// isFloat check if val is a valid SIC float (48 bits)
func isFloat(val float64) bool {
	return val > 0 && val < math.Pow(2, 48)
}

// isAddr checks if addr is a valid SIC address
func isAddr(addr int) bool {
	return addr >= 0 && addr <= MAX_ADDRESS
}

// isDevice check if dev is a valid SIC device
func isDevice(dev int) bool {
	return dev >= 0 && dev <= 255
}
