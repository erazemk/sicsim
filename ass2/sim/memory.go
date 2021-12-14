package sicsim

import (
	"encoding/binary"
	"fmt"
)

// Byte returns the byte at m[addr]
func (m Machine) Byte(addr int) (byte, error) {
	if isAddr(addr) {
		return m.mem[addr], nil
	}

	return 0, fmt.Errorf("not a valid address: %d", addr)
}

// SetByte sets the byte at the address addr to val
func (m *Machine) SetByte(addr int, val byte) error {
	if isAddr(addr) {
		m.mem[addr] = val
		return nil
	}

	return fmt.Errorf("not a valid address: %d", addr)
}

// Word returns the word at m[addr..addr+2]
func (m Machine) Word(addr int) (int, error) {
	if isAddr(addr) {
		buf := make([]byte, 4)
		buf[0] = 0
		buf[1] = m.mem[addr]
		buf[2] = m.mem[addr+1]
		buf[3] = m.mem[addr+2]
		word := int(binary.BigEndian.Uint32(buf))
		return word, nil
	}

	return 0, fmt.Errorf("not a valid address: %d", addr)
}

// SetWord sets the word (3 bytes) at addr to val
func (m *Machine) SetWord(addr, val int) error {
	if isAddr(addr) && isWord(val) {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(val))

		fmt.Printf("Buf: %[1]v\n", buf)

		// buf[0] == MSB, which is too big for SIC words, so it isn't used
		m.mem[addr] = buf[1]
		m.mem[addr+1] = buf[2]
		m.mem[addr+2] = buf[3]
		return nil
	}

	return fmt.Errorf("not a valid address or value: %d, %d", addr, val)
}

// Mem prints the content of the memory from startAddr to endAddr
func (m *Machine) Mem(startAddr, endAddr int) string {
	return fmt.Sprintf("%v", m.mem[startAddr:endAddr])
}
