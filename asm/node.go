package asm

import (
	"fmt"
	"strings"
)

type Node struct {
	// Line format: [label] [mnemonic] [operands/symbols]
	label     string
	length    int
	mnemonic  string
	operand   int
	symbol    string
	ni        byte
	indexed   int
	brelative bool
	lc        int
}

func NewNode(command []string, lc int, brelative bool) Node {
	var n Node
	n.brelative = brelative
	n.lc = lc

	// Check if label exits
	if !inSlice(command[0], Mnemonics) {
		n.label, command = command[0], command[1:]
	}

	// Parse mnemonic
	n.mnemonic, command = command[0], command[1:]

	// Parse operand and special bits
	if len(command) > 0 {
		n.ParseOperands(command)
	}

	// Set node length
	if inSlice(n.mnemonic, Directives) {
		n.length = 0
	} else if inSlice(n.mnemonic, StorageDirectives) {
		switch n.mnemonic {
		case "BYTE":
			n.length = 1
		case "WORD":
			n.length = 3
		case "RESB": // Assume that n.operand is a number in this case
			n.length = n.operand
		case "RESW":
			n.length = 3 * n.operand
		}
	} else if inSlice(n.mnemonic, InstructionF1) {
		n.length = 1
	} else if inSlice(n.mnemonic, InstructionsF2) {
		n.length = 2
	} else if inSlice(n.mnemonic, InstructionsF3) {
		n.length = 3
	} else {
		panic(fmt.Sprintf("Invalid mnemonic: %s", n.mnemonic))
	}

	return n
}

// Bytes returns a hexadecimal representation of a node
func (n *Node) Bytes() string {
	var bytes string

	if inSlice(n.mnemonic, StorageDirectives) {
		if n.mnemonic == "BYTE" {
			bytes = fmt.Sprintf("%02X", n.operand)
		}
		if n.mnemonic == "WORD" {
			bytes = Word(n.operand)
		}
	} else if inSlice(n.mnemonic, InstructionF1) {
		bytes = fmt.Sprintf("%02X", Opcodes[n.mnemonic])
	} else if inSlice(n.mnemonic, InstructionsF2) {
		bytes = fmt.Sprintf("%02X%s", Opcodes[n.mnemonic], Word(n.operand))
	} else if inSlice(n.mnemonic, InstructionsF3) {
		var opcode, bp, e, operand int

		// Get proper opcode
		opcode = int(Opcodes[n.mnemonic] | n.ni)

		if debug {
			fmt.Printf("Opcode for '%s' is %s (ni=%d)\n", n.mnemonic, Word(opcode), n.ni)
		}

		// Check if extended
		e = 0 // For now extended will always be zero

		// Check if using PC relative addressing
		if rang := n.operand - n.lc - n.length; rang >= -2048 && rang <= 2047 { // PC-relative
			bp = 0x01
			operand = rang

			if debug {
				fmt.Printf("Rang for '%s' is %06X (%d, %d, %d)\n", n.mnemonic, rang, n.operand, n.lc, n.length)
			}
		} else { // Direct
			operand = n.operand
		}

		if debug {
			fmt.Printf("Operand for '%s' is %s (before)\n", n.mnemonic, Word(operand))
		}

		operand = n.indexed<<15 | bp<<13 | e<<12 | operand&0x0FFF

		if debug {
			fmt.Printf("Operand for '%s' is %s (after)\n", n.mnemonic, Word(operand))
		}

		bytes = fmt.Sprintf("%02X%02X%02X", opcode, operand>>8, operand&0xFF)

		if debug {
			fmt.Printf("Bytes for '%s' are %s\n", n.mnemonic, bytes)
		}
	}

	return bytes
}

// ParseOperands sets all of node's attributes based on received operands
func (n *Node) ParseOperands(operands []string) {
	if inSlice(n.mnemonic, Directives) || inSlice(n.mnemonic, StorageDirectives) { // One parameter
		if num, ok := isNumber(operands[0]); ok {
			n.operand = num
		} else {
			n.symbol = operands[0]
		}

		n.ni = 0
		n.indexed = 0
	} else if inSlice(n.mnemonic, InstructionF2r) { // One register
		reg := strings.IndexRune("AXLBSTF", []rune(operands[0])[0])

		n.operand = reg << 4
		n.ni = 0
		n.indexed = 0
	} else if inSlice(n.mnemonic, InstructionF2rn) { // One register, one number
		var num, reg int
		var ok bool

		switch len(operands) {
		case 1: // ["r,n"]
			arr := strings.Split(operands[0], ",")
			reg = strings.IndexRune("AXLBSTF", []rune(arr[0])[0])

			num, ok = isNumber(arr[1])
			if !ok {
				panic(fmt.Sprintf("not a number: '%s'", arr[1]))
			}
		case 2: // ["r," "n"] or ["r" ",n"]
			regRune := []rune(strings.Split(operands[0], ",")[0])[0] // Register char
			numStr := strings.Split(operands[1], ",")[0]             // Number without ','
			reg = strings.IndexRune("AXLBSTF", regRune)

			num, ok = isNumber(numStr)
			if !ok {
				panic(fmt.Sprintf("not a number: '%s'", numStr))
			}
		case 3: // ["r" "," "n"]
			regRune := []rune(operands[0])[0] // Register char
			reg = strings.IndexRune("AXLBSTF", regRune)

			num, ok = isNumber(operands[2])
			if !ok {
				panic(fmt.Sprintf("not a number: '%s'", operands[2]))
			}
		}

		n.operand = reg<<4 | num
		n.ni = 0
		n.indexed = 0
	} else if inSlice(n.mnemonic, InstructionF2rr) { // Two registers
		var r1, r2 int

		switch len(operands) {
		case 1: // ["r1,r2"]
			arr := strings.Split(operands[0], ",")
			regRune1 := []rune(arr[0])[0]
			regRune2 := []rune(arr[1])[0]

			r1 = strings.IndexRune("AXLBSTF", regRune1)
			r2 = strings.IndexRune("AXLBSTF", regRune2)
		case 2: // ["r1," "r2"] or ["r1" ",r2"]
			regRune1 := []rune(strings.Split(operands[0], ",")[0])[0] // Register char
			regRune2 := []rune(strings.Split(operands[1], ",")[0])[0] // Register char

			r1 = strings.IndexRune("AXLBSTF", regRune1)
			r2 = strings.IndexRune("AXLBSTF", regRune2)
		case 3: // ["r" "," "n"]
			regRune1 := []rune(operands[0])[0] // Register char
			regRune2 := []rune(operands[2])[0] // Register char

			r1 = strings.IndexRune("AXLBSTF", regRune1)
			r2 = strings.IndexRune("AXLBSTF", regRune2)
		}

		n.operand = r1<<4 | r2
		n.ni = 0
		n.indexed = 0
	} else if inSlice(n.mnemonic, InstructionF3) { // No operands
		// No need to do anything
	} else if inSlice(n.mnemonic, InstructionF3m) { // One number or one number, X
		var num int
		var ok bool

		// Check if using indexed addressing (need ,X because variable could end in X - e.g. 'ADD SIX')
		if strings.HasSuffix(operands[0], ",X") { // ["num,X"]
			n.indexed = 1
		}
		if len(operands) == 2 && (strings.HasSuffix(operands[0], ",") ||
			strings.HasSuffix(operands[1], "X")) { // ["num," "X"] or ["num" ",X"]
			n.indexed = 1
		}

		m := strings.Split(operands[0], ",")[0]

		// Check if using direct or indirect addressing
		if strings.HasPrefix(m, "#") {
			n.ni = 0x01

			num, ok = isNumber(m[1:])
			if ok {
				n.operand = num
			} else {
				panic(fmt.Sprintf("not a number: '%s'", m[1:]))
			}
		} else if strings.HasPrefix(m, "@") {
			n.ni = 0x02

			num, ok = isNumber(m[1:])
			if ok {
				n.operand = num
			} else {
				panic(fmt.Sprintf("not a number: '%s'", m[1:]))
			}
		} else {
			n.ni = 0x03

			num, ok = isNumber(m)
			if ok {
				n.operand = num
			} else {
				n.symbol = m
			}
		}
	}
}

func (n *Node) pretty() string {
	deb := debug
	SetDebug(false)
	str := fmt.Sprintf("%s\t%s\t%-6s\t%-6s\t%-6s", Word(n.lc), n.Bytes(), n.label, n.mnemonic, n.symbol)
	SetDebug(deb)
	return str
}
