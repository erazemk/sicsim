package sim

import (
	"encoding/binary"
	"fmt"
)

// Keeps track of jump addresses to detect loops - halt execution
var JMPADDR int
var LASTINST byte

// fetch returns a byte from m[PC] and increments PC
func (m *Machine) fetch() byte {
	addr := m.PC()
	m.SetPC(addr + 1)
	val, _ := m.Byte(addr)
	return val
}

// Execute executes each fetched instruction
func (m *Machine) Execute() error {
	var success bool
	var err error

	opcode := m.fetch()

	if success, err = m.execF1(opcode); err == nil {
		if success {
			LASTINST = opcode
			return nil
		}
	} else {
		return fmt.Errorf("failed to execute command (format 1): %w", err)
	}

	operands := m.fetch()

	if success, err = m.execF2(opcode, operands); err == nil {
		if success {
			LASTINST = opcode
			return nil
		}
	} else {
		return fmt.Errorf("failed to execute command (format 2): %w", err)
	}

	ni := opcode & 0x03
	opcode = opcode & 0xFC

	if success, err = m.execSICF3F4(opcode, operands, ni); err == nil {
		if success {
			LASTINST = opcode
			return nil
		}
	} else {
		return fmt.Errorf("failed to execute command (format 3): %w", err)
	}

	return fmt.Errorf("failed to execute command: not a valid SIC command (any format)")
}

// calcStoreOperand returns the proper operand for store instructions
func (m *Machine) calcStoreOperand(addr int, indirect bool) int {
	if indirect {
		addr, _ = m.Word(addr)
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

// calcByteOperand returns the proper byte for non-store instructions
func (m *Machine) calcByteOperand(operand int, indirect, immediate bool) byte {
	if immediate {
		return byte(operand & 0xFF)
	}

	if indirect {
		val, _ := m.Word(operand)
		op, _ := m.Byte(val)
		return op
	}

	op, _ := m.Byte(operand)
	return op
}

// execF1 tries to execute opcode as format 1
func (m *Machine) execF1(opcode byte) (bool, error) {
	switch opcode {
	case FIX:
		return false, fmt.Errorf("instruction not implemented: %s", "FIX")
	case FLOAT:
		return false, fmt.Errorf("instruction not implemented: %s", "FLOAT")
	case HIO:
		return false, fmt.Errorf("instruction not implemented: %s", "HIO")
	case NORM:
		return false, fmt.Errorf("instruction not implemented: %s", "NORM")
	case SIO:
		return false, fmt.Errorf("instruction not implemented: %s", "SIO")
	case TIO:
		return false, fmt.Errorf("instruction not implemented: %s", "TIO")
	default:
		// Not a format 1 instruction
		return false, nil
	}

	// Currently unreachable
	// return true, nil
}

// execF2 tries to execute opcode as format 2
func (m *Machine) execF2(opcode, operand byte) (bool, error) {
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
	case SVC:
		return false, fmt.Errorf("instruction not implemented: %s", "SVC")
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
		// Not a format 2 instruction
		return false, nil
	}

	return true, nil
}

// execSICF3F4 tries to execute opcode either in SIC format, format 3 or format 4
func (m *Machine) execSICF3F4(opcode, operands, ni byte) (bool, error) {
	var extended, indexed bool
	var direct, baserelative, pcrelative bool // BP bits
	var immediate, indirect, sic bool         // NI bits
	var operand int

	// Addressing modes
	if ni == 0x00 {
		sic = true
	} else { // Can't be combined with SIC mode
		if operands&0x10 == 0x10 {
			extended = true
		}

		if bp := operands & 0x60; bp == 0x00 {
			direct = true
		} else if bp == 0x40 {
			baserelative = true
		} else if bp == 0x20 {
			pcrelative = true
		} else if bp == 0x60 {
			return false, fmt.Errorf("wrong addressing format")
		}

		if ni == 0x02 {
			indirect = true
		} else if ni == 0x01 {
			immediate = true
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
			operand -= 4096
		}

		operand += m.PC()
	}

	if indexed {
		operand += m.X()
	}

	if debug {
		fmt.Printf("Instruction: 0x%02X (0x%04X)\n", opcode, operand)
		fmt.Println("Addressing:")
		fmt.Printf("  SIC: %v\n", sic)
		fmt.Printf("  indirect: %v\n", indirect)
		fmt.Printf("  direct: %v\n", direct)
		fmt.Printf("  extended: %v\n", extended)
		fmt.Printf("  indexed: %v\n", indexed)
		fmt.Printf("  immediate: %v\n", immediate)
		fmt.Printf("  base relative: %v\n", baserelative)
		fmt.Printf("  pc relative: %v\n", pcrelative)
	}

	switch opcode {
	case ADD:
		m.SetA(m.A() + m.calcOperand(operand, indirect, immediate))
	case ADDF:
		return false, fmt.Errorf("instruction not implemented: %s", "ADDF")
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
	case COMPF:
		return false, fmt.Errorf("instruction not implemented: %s", "COMPF")
	case DIV:
		m.SetA(m.A() / m.calcOperand(operand, indirect, immediate))
	case DIVF:
		return false, fmt.Errorf("instruction not implemented: %s", "DIVF")
	case J:
		addr := m.calcStoreOperand(operand, indirect)

		if debug {
			fmt.Printf("Jump addr: 0x%02X\n", addr)
		}

		// Halt processor
		if addr == JMPADDR && LASTINST == J {
			m.halted = true
			return true, nil
		} else {
			JMPADDR = addr
		}

		m.SetPC(addr)
	case JEQ:
		if m.SW() == EQ {
			m.SetPC(m.calcStoreOperand(operand, indirect))
		}
	case JGT:
		if m.SW() == GT {
			m.SetPC(m.calcStoreOperand(operand, indirect))
		}
	case JLT:
		if m.SW() == LT {
			m.SetPC(m.calcStoreOperand(operand, indirect))
		}
	case JSUB:
		m.SetL(m.PC())
		m.push(m.PC())
		m.SetPC(m.calcStoreOperand(operand, indirect))
	case LDA:
		m.SetA(m.calcOperand(operand, indirect, immediate))
	case LDB:
		m.SetB(m.calcOperand(operand, indirect, immediate))
	case LDCH:
		m.SetALow(m.calcByteOperand(operand, indirect, immediate))
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
	case LPS:
		return false, fmt.Errorf("instruction not implemented: %s", "LPS")
	case MUL:
		m.SetA(m.A() * m.calcOperand(operand, indirect, immediate))
	case MULF:
		return false, fmt.Errorf("instruction not implemented: %s", "MULF")
	case OR:
		m.SetA(m.A() | m.calcOperand(operand, indirect, immediate))
	case RD:
		char, err := m.ReadDevice(m.calcByteOperand(operand, indirect, immediate))
		if err != nil {
			fmt.Println(err)
			return false, err
		}

		m.SetALow(char)
	case RSUB:
		m.SetPC(m.L())
		m.SetL(m.pop())
	case SSK:
		return false, fmt.Errorf("instruction not implemented: %s", "SSK")
	case STA:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.A())
	case STB:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.B())
	case STCH:
		m.SetByte(m.calcStoreOperand(operand, indirect), m.ALow())
	case STF:
		m.SetWord(m.calcStoreOperand(operand, indirect), m.F())
	case STI:
		return false, fmt.Errorf("instruction not implemented: %s", "STI")
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
	case SUBF:
		return false, fmt.Errorf("instruction not implemented: %s", "SUBF")
	case TD:
		m.TestDevice(m.calcByteOperand(operand, indirect, immediate))
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
		err := m.WriteDevice(m.calcByteOperand(operand, indirect, immediate), m.ALow())
		if err != nil {
			fmt.Println(err)
			return false, err
		}
	default:
		// Not a format 3, 4, SIC instruction
		return false, nil
	}

	return true, nil
}
