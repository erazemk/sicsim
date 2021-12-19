package sicsim

import (
	"encoding/binary"
)

// fetch returns a byte from m[PC] and increments PC
func (m *Machine) fetch() byte {
	addr := m.PC()
	m.SetPC(addr + 1)
	return m.mem[addr]
}

// Execute executes each fetched instruction
func (m *Machine) Execute() {
	opcode := m.fetch()

	if m.execF1(opcode) {
		return
	}

	operands := m.fetch()

	if m.execF2(opcode, operands) {
		return
	}

	ni := opcode & 0x03
	opcode = opcode & 0xFC

	if m.execSICF3F4(opcode, operands, ni) {
		return
	}

	// TODO: Execute() should return error if none of the formats were correct
}

// calcStoreOperand returns the proper operand for store instructions
func (m *Machine) calcStoreOperand(addr int, indirect bool) int {
	if indirect {
		val, _ := m.Word(addr)
		return val
	}

	return addr
}

// calcOperand returns the proper operand for non-store instructions
func (m *Machine) calcOperand(operand int, indirect, immediate bool) int {
	if immediate {
		return operand
	}

	operand, _ = m.Word(operand)

	if indirect {
		operand, _ = m.Word(operand)
	}

	return operand
}

// execF1 tries to execute opcode as format 1
func (m *Machine) execF1(opcode byte) bool {
	switch opcode {
	// case FIX:
	// case FLOAT:
	// case HIO:
	// case NORM:
	// case SIO:
	// case TIO:
	default:
		// Not implemented
		return false
	}

	// Currently unreachable
	//return true
}

// execF2 tries to execute opcode as format 2
func (m *Machine) execF2(opcode, operand byte) bool {
	op1 := int((operand & 0xF0) >> 4)
	op2 := int(operand & 0x0F)

	switch opcode {
	case ADDR:
		r1, _ := m.Reg(op1)
		r2, _ := m.Reg(op2)
		m.SetReg(op2, r2+r1)
	case CLEAR:
		m.SetReg(op1, 0)
	case COMPR:
		r1, _ := m.Reg(op1)
		r2, _ := m.Reg(op2)

		if r1 > r2 {
			m.SetSW(GT)
		} else if r1 == r2 {
			m.SetSW(EQ)
		} else {
			m.SetSW(LT)
		}
	case DIVR:
		r1, _ := m.Reg(op1)
		r2, _ := m.Reg(op2)
		m.SetReg(op2, r2/r1)
	case MULR:
		r1, _ := m.Reg(op1)
		r2, _ := m.Reg(op2)
		m.SetReg(op2, r2*r1)
	case RMO:
		r1, _ := m.Reg(op1)
		m.SetReg(op2, r1)
	case SHIFTL:
		r1, _ := m.Reg(op1)
		m.SetReg(op1, r1<<op2)
	case SHIFTR:
		r1, _ := m.Reg(op1)
		m.SetReg(op1, r1>>op2)
	case SUBR:
		r1, _ := m.Reg(op1)
		r2, _ := m.Reg(op2)
		m.SetReg(op2, r2-r1)
	// case SVC:
	case TIXR:
		r1, _ := m.Reg(op1)
		m.SetX(m.X() + 1)
		rX := m.X()

		if rX > r1 {
			m.SetSW(GT)
		} else if rX == r1 {
			m.SetSW(EQ)
		} else {
			m.SetSW(LT)
		}
	default:
		// Not implemented
		return false
	}

	return true
}

// execSICF3F4 tries to execute opcode either in SIC format, format 3 or format 4
func (m *Machine) execSICF3F4(opcode, operands, ni byte) bool {
	var extended, indexed bool
	var immediate, baserelative, pcrelative bool
	var indirect, sic bool
	var target_address, operand int

	// Addressing modes
	if ni == 0x00 {
		sic = true
	} else { // Can't be combined with SIC mode
		if operands&0x10 == 0x10 {
			extended = true
		}

		if bp := operands & 0x60; bp == 0x00 {
			immediate = true
		} else if bp == 0x40 {
			baserelative = true
		} else if bp == 0x20 {
			pcrelative = true
		} else if bp == 0x60 {
			panic("wrong addressing format")
		}

		if ni == 0x02 {
			indirect = true
		}
	}

	if operands&0x80 == 0x80 {
		indexed = true
	}

	if sic {
		operand = int(binary.BigEndian.Uint32([]byte{0, 0, operands & 0x7F, m.fetch()}))
	} else if extended {
		operand = int(binary.BigEndian.Uint32([]byte{0, operands & 0x0F, m.fetch(), m.fetch()}))
	} else {
		operand = int(binary.BigEndian.Uint32([]byte{0, 0, operands & 0x0F, m.fetch()}))
	}

	if baserelative {
		operand += m.B()
	} else if pcrelative {
		if operand >= 2048 {
			operand = ^operand + 1
		}

		operand += m.PC()
	}

	if indexed {
		operand += m.X()
	}

	switch opcode {
	case ADD:
		m.SetA(m.A() + m.calcOperand(operand, indirect, immediate))
	// case ADDF:
	case AND:
		m.SetA(m.A() & m.calcOperand(operand, indirect, immediate))
	case COMP:
		rA := m.A()
		val := m.calcOperand(operand, indirect, immediate)

		if rA > val {
			m.SetSW(GT)
		} else if rA == val {
			m.SetSW(EQ)
		} else {
			m.SetSW(LT)
		}
	// case COMPF:
	case DIV:
		m.SetA(m.A() / m.calcOperand(operand, indirect, immediate))
	// case DIVF:
	case J:
		m.SetPC(m.calcOperand(operand, indirect, immediate))
	case JEQ:
		if m.SW() == EQ {
			m.SetPC(m.calcOperand(operand, indirect, immediate))
		}
	case JGT:
		if m.SW() == GT {
			m.SetPC(m.calcOperand(operand, indirect, immediate))
		}
	case JLT:
		if m.SW() == LT {
			m.SetPC(m.calcOperand(operand, indirect, immediate))
		}
	case JSUB:
		m.SetL(m.PC())
		m.SetPC(m.calcOperand(operand, indirect, immediate))
	case LDA:
		m.SetA(m.calcOperand(operand, indirect, immediate))
	case LDB:
		m.SetB(m.calcOperand(operand, indirect, immediate))
	case LDCH:
		m.SetALow(byte(m.calcOperand(operand, indirect, immediate)))
	case LDF:
		m.SetF(m.calcOperand(operand, indirect, immediate))
	case LDL:
		m.SetL(m.calcOperand(operand, indirect, immediate))
	case LDS:
		m.SetS(m.calcOperand(operand, indirect, immediate))
	case LDT:
		m.SetT(m.calcOperand(operand, indirect, immediate))
	case LDX:
		m.SetX(m.calcOperand(operand, indirect, immediate))
	// case LPS:
	case MUL:
		m.SetA(m.A() * m.calcOperand(operand, indirect, immediate))
	// case MULF:
	case OR:
		m.SetA(m.A() | m.calcOperand(operand, indirect, immediate))
	case RD:
		char, _ := m.ReadDevice(byte(m.calcOperand(operand, indirect, immediate)))
		m.SetALow(char)
	case RSUB:
		m.SetPC(m.L())
	// case SSK:
	case STA:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.A())
	case STB:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.B())
	case STCH:
		m.SetByte(m.calcStoreOperand(operand, indirect), m.ALow())
	case STF:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.F())
	// case STI:
	case STL:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.L())
	case STS:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.S())
	case STSW:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.SW())
	case STT:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.T())
	case STX:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.X())
	case SUB:
		m.SetA(m.A() - m.calcOperand(operand, indirect, immediate))
	// case SUBF:
	case TD:
		m.TestDevice(byte(m.calcOperand(operand, indirect, immediate)))
	case TIX:
		m.SetX(m.X() + 1)
		rX := m.X()
		val := m.calcOperand(operand, indirect, immediate)

		if rX > val {
			m.SetSW(GT)
		} else if rX == val {
			m.SetSW(EQ)
		} else {
			m.SetSW(LT)
		}
	case WD:
		m.WriteDevice(byte(m.calcOperand(operand, indirect, immediate)), m.ALow())
	default:
		// Not implemented
		return false
	}

	return true
}
