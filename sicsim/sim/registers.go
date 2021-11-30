package sim

import "fmt"

type registers struct {
	a  [3]byte
	x  [3]byte
	l  [3]byte
	b  [3]byte
	s  [3]byte
	t  [3]byte
	f  [6]byte
	pc [3]byte
	sw [1]byte
}

// SW register values
const (
	LT = 0x00
	EQ = 0x40
	GT = 0x80
)

// Reg returns the value of register reg
func (m Machine) Reg(reg int) ([]byte, error) {
	switch reg {
	case 0:
		return m.regs.a[:], nil
	case 1:
		return m.regs.x[:], nil
	case 2:
		return m.regs.l[:], nil
	case 3:
		return m.regs.b[:], nil
	case 4:
		return m.regs.s[:], nil
	case 5:
		return m.regs.t[:], nil
	case 6:
		return m.regs.f[:], nil
	case 8:
		return m.regs.pc[:], nil
	case 9:
		return m.regs.sw[:], nil
	}

	return nil, fmt.Errorf("Not a valid register: %d", reg)
}

// SetReg sets the value of register reg
func (m *Machine) SetReg(reg int, val []byte) error {
	if !isRegister(reg) {
		return fmt.Errorf("Not a valid register: %d", reg)
	}

	switch reg {
	case 0:
		copy(m.regs.a[:], val)
	case 1:
		copy(m.regs.x[:], val)
	case 2:
		copy(m.regs.l[:], val)
	case 3:
		copy(m.regs.b[:], val)
	case 4:
		copy(m.regs.s[:], val)
	case 5:
		copy(m.regs.t[:], val)
	case 6:
		copy(m.regs.f[:], val)
	case 8:
		copy(m.regs.pc[:], val)
	case 9:
		copy(m.regs.sw[:], val)
	}

	return nil
}

// A returns the value of
func (m Machine) A() [3]byte {
	return m.regs.a
}

// X returns the value of the X register
func (m Machine) X() [3]byte {
	return m.regs.x
}

// L returns the value of the L register
func (m Machine) L() [3]byte {
	return m.regs.l
}

// B returns the value of the B register
func (m Machine) B() [3]byte {
	return m.regs.b
}

// S returns the value of the S register
func (m Machine) S() [3]byte {
	return m.regs.s
}

// T returns the value of the T register
func (m Machine) T() [3]byte {
	return m.regs.t
}

// F returns the value of the F register
func (m Machine) F() [6]byte {
	return m.regs.f
}

// PC returns the value of the PC register
func (m Machine) PC() [3]byte {
	return m.regs.pc
}

// SW returns the value of the SW register
func (m Machine) SW() [1]byte {
	return m.regs.sw
}

// SetA sets the value of the A register
func (m *Machine) SetA(val []byte) {
	if isWord(val) {
		copy(m.regs.a[:], val)
	}
}

// SetX sets the value of the X register
func (m *Machine) SetX(val []byte) {
	if isWord(val) {
		copy(m.regs.x[:], val)
	}
}

// SetL sets the value of the L register
func (m *Machine) SetL(val []byte) {
	if isWord(val) {
		copy(m.regs.l[:], val)
	}
}

// SetB sets the value of the B register
func (m *Machine) SetB(val []byte) {
	if isWord(val) {
		copy(m.regs.b[:], val)
	}
}

// SetS sets the value of the S register
func (m *Machine) SetS(val []byte) {
	if isWord(val) {
		copy(m.regs.s[:], val)
	}
}

// SetT sets the value of the T register
func (m *Machine) SetT(val []byte) {
	if isWord(val) {
		copy(m.regs.t[:], val)
	}
}

// SetF sets the value of the F register
func (m *Machine) SetF(val []byte) {
	if isFloat(val) {
		copy(m.regs.f[:], val)
	}
}

// SetPC sets the value of the PC register
func (m *Machine) SetPC(val []byte) {
	if isWord(val) {
		copy(m.regs.pc[:], val)
	}
}

// SetSW sets the value of the SW register
func (m *Machine) SetSW(val []byte) {
	if isWord(val) {
		copy(m.regs.sw[:], val)
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
