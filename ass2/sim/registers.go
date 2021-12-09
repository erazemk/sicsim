package sim

import (
	"encoding/binary"
	"fmt"
)

type registers struct {
	a  int
	x  int
	l  int
	b  int
	s  int
	t  int
	f  int
	pc int
	sw int
}

// SW register values
const (
	LT = 0x00
	EQ = 0x40
	GT = 0x80
)

// Reg returns the value of register reg
func (m Machine) Reg(reg int) (int, error) {
	switch reg {
	case 0:
		return m.regs.a, nil
	case 1:
		return m.regs.x, nil
	case 2:
		return m.regs.l, nil
	case 3:
		return m.regs.b, nil
	case 4:
		return m.regs.s, nil
	case 5:
		return m.regs.t, nil
	case 6:
		return m.regs.f, nil
	case 8:
		return m.regs.pc, nil
	case 9:
		return m.regs.sw, nil
	}

	return -1, fmt.Errorf("not a valid register: %d", reg)
}

// SetReg sets the value of register reg
func (m *Machine) SetReg(reg int, val int) error {
	if !isRegister(reg) {
		return fmt.Errorf("not a valid register: %d", reg)
	}

	switch reg {
	case 0:
		m.regs.a = val
	case 1:
		m.regs.x = val
	case 2:
		m.regs.l = val
	case 3:
		m.regs.b = val
	case 4:
		m.regs.s = val
	case 5:
		m.regs.t = val
	case 6:
		m.regs.f = val
	case 8:
		m.regs.pc = val
	case 9:
		m.regs.sw = val
	}

	return nil
}

// A returns the value of the A register
func (m Machine) A() int {
	return m.regs.a
}

// ALow returns the lowest byte of the A register
func (m Machine) ALow() byte {
	low := make([]byte, 4)
	binary.LittleEndian.PutUint32(low, uint32(m.regs.a))
	return low[0]
}

// X returns the value of the X register
func (m Machine) X() int {
	return m.regs.x
}

// L returns the value of the L register
func (m Machine) L() int {
	return m.regs.l
}

// B returns the value of the B register
func (m Machine) B() int {
	return m.regs.b
}

// S returns the value of the S register
func (m Machine) S() int {
	return m.regs.s
}

// T returns the value of the T register
func (m Machine) T() int {
	return m.regs.t
}

// F returns the value of the F register
func (m Machine) F() int {
	return m.regs.f
}

// PC returns the value of the PC register
func (m Machine) PC() int {
	return m.regs.pc
}

// SW returns the value of the SW register
func (m Machine) SW() int {
	return m.regs.sw
}

// SetA sets the value of the A register
func (m *Machine) SetA(val int) {
	if isWord(val) {
		m.regs.a = val
	}
}

// SetALow sets the value of lowest byte of the A register
func (m *Machine) SetALow(val byte) {
	word := make([]byte, 4)
	word[0] = val
	m.regs.a = int(binary.LittleEndian.Uint32(word))
}

// SetX sets the value of the X register
func (m *Machine) SetX(val int) {
	if isWord(val) {
		m.regs.x = val
	}
}

// SetL sets the value of the L register
func (m *Machine) SetL(val int) {
	if isWord(val) {
		m.regs.l = val
	}
}

// SetB sets the value of the B register
func (m *Machine) SetB(val int) {
	if isWord(val) {
		m.regs.b = val
	}
}

// SetS sets the value of the S register
func (m *Machine) SetS(val int) {
	if isWord(val) {
		m.regs.s = val
	}
}

// SetT sets the value of the T register
func (m *Machine) SetT(val int) {
	if isWord(val) {
		m.regs.t = val
	}
}

// SetF sets the value of the F register
func (m *Machine) SetF(val int) {
	if isWord(val) {
		m.regs.f = val
	}
}

// SetPC sets the value of the PC register
func (m *Machine) SetPC(val int) {
	if isWord(val) {
		m.regs.pc = val
	}
}

// SetSW sets the value of the SW register
func (m *Machine) SetSW(val int) {
	if isWord(val) {
		m.regs.sw = val
	}
}

// Print outputs the machine's register state
func (m *Machine) Registers() string {
	return fmt.Sprintf(
		"A:  %06X (Dec: %d)\nX:  %06X (Dec: %d)\nL:  %06X (Dec: %d)\nB:  %06X (Dec: %d)\nS:  %06X (Dec: %d)\n"+
			"T:  %06X (Dec: %d)\nF:  %06X (Dec: %d)\nSP: %06X (Dec: %d)\nSW: %06X (Dec: %d)",
		m.regs.a, m.regs.a, m.regs.x, m.regs.x, m.regs.l, m.regs.l, m.regs.b, m.regs.b,
		m.regs.s, m.regs.s, m.regs.t, m.regs.t, m.regs.f, m.regs.f, m.regs.pc, m.regs.pc,
		m.regs.sw, m.regs.sw)
}
