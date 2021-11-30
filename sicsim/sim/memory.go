package sim

import "encoding/binary"

// Byte returns the byte at m[addr]
func (m Machine) Byte(addr uint32) byte {
	if isAddr(addr) {
		return m.mem[addr]
	}

	return 0
}

// SetByte sets the byte at the address addr to val
func (m *Machine) SetByte(addr uint32, val byte) {
	if isAddr(addr) {
		m.mem[addr] = val
	}
}

// Word returns the word at m[addr..addr+2]
func (m Machine) Word(addr uint32) int {
	if isAddr(addr) {
		val := m.mem[addr : addr+3]
		return int(binary.BigEndian.Uint32(val))
	}

	return -1
}

// SetWord sets the word (3 bytes) at addr to val
func (m *Machine) SetWord(addr uint32, val int) {
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
func (m Machine) Float(addr uint32) float64 {
	if isAddr(addr) {
		val := m.mem[addr : addr+6]
		return float64(binary.BigEndian.Uint64(val))
	}

	return -1
}

// SetFloat sets the float (6 bytes) at addr to val
func (m *Machine) SetFloat(addr uint32, val float64) {
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
