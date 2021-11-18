package sim

import (
	"encoding/binary"
	"fmt"
	"log"
)

const MAX_ADDRESS = 1048576

type Machine struct {
	Regs Registers // TODO: Figure out what to do with registers
	mem  [MAX_ADDRESS + 1]byte
	Devs [256]Device
}

// New creates a new machine
func (m *Machine) New() {
	m.Devs[0].New(0) // stdin
	m.Devs[1].New(1) // stdout
	m.Devs[2].New(2) // stderr

	if debug {
		log.Println("Created a new machine")
	}
}

// Byte returns the byte at m[addr]
func (m Machine) Byte(addr int) byte {
	if isAddr(addr) {
		return m.mem[addr]
	}

	return 0
}

// SetByte sets the byte at the address addr to val
func (m *Machine) SetByte(addr int, val byte) {
	if isAddr(addr) {
		m.mem[addr] = val
	}
}

// Word returns the word at m[addr..addr+2]
func (m Machine) Word(addr int) int {
	if isAddr(addr) {
		val := m.mem[addr : addr+3]
		return int(binary.BigEndian.Uint32(val))
	}

	return -1
}

// SetWord sets the word (3 bytes) at addr to val
func (m *Machine) SetWord(addr, val int) {
	if isAddr(addr) && isWord(val) {
		bytes := make([]byte, 3)
		binary.BigEndian.PutUint32(bytes, uint32(val))

		// TODO: Optimize this shit
		m.mem[addr] = bytes[0]
		m.mem[addr+1] = bytes[1]
		m.mem[addr+2] = bytes[2]
	}
}

// Float returns the float at m[addr..addr+5]
func (m Machine) Float(addr int) float64 {
	if isAddr(addr) {
		val := m.mem[addr : addr+6]
		return float64(binary.BigEndian.Uint64(val))
	}

	return -1
}

// SetFloat sets the float (6 bytes) at addr to val
func (m *Machine) SetFloat(addr int, val float64) {
	if isAddr(addr) && isFloat(val) {
		bytes := make([]byte, 6)
		binary.BigEndian.PutUint64(bytes, uint64(val))

		// TODO: Optimize this shit
		m.mem[addr] = bytes[0]
		m.mem[addr+1] = bytes[1]
		m.mem[addr+2] = bytes[2]
		m.mem[addr+3] = bytes[3]
		m.mem[addr+4] = bytes[4]
		m.mem[addr+5] = bytes[5]
	}
}

// Device returns a device by its number
func (m *Machine) Device(num int) Device {
	if isDevice(num) {
		return m.Devs[num]
	}

	return Device{}
}

// Print outputs the machine's register state
func (m *Machine) Registers() string {
	return fmt.Sprintf(
		"A:  %06X (Dec: %d)\nX:  %06X (Dec: %d)\nL:  %06X (Dec: %d)\nB:  %06X (Dec: %d)\nS:  %06X (Dec: %d)\n"+
			"T:  %06X (Dec: %d)\nF:  %06X (Dec: %d)\nSP: %06X (Dec: %d)\nSW: %06X (Dec: %d)",
		m.Regs.A(), m.Regs.A(), m.Regs.X(), m.Regs.X(), m.Regs.L(), m.Regs.L(), m.Regs.B(), m.Regs.B(),
		m.Regs.S(), m.Regs.S(), m.Regs.T(), m.Regs.T(), m.Regs.F(), m.Regs.F(), m.Regs.PC(), m.Regs.PC(),
		m.Regs.SW(), m.Regs.SW())
}
