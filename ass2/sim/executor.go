package sicsim

import (
	"encoding/binary"
	"fmt"
)

// fetch returns a byte from m[PC] and increments PC
func (m *Machine) fetch() byte {
	addr := m.PC()
	m.SetPC(addr + 1)
	return m.mem[addr]
}

// Execute tries to execute fetched operation
func (m *Machine) Execute() {
	opcode := m.fetch()

	if m.execF1(opcode) == nil {
		return
	}

	operands := m.fetch()
	op1 := int(operands & 0x0F)
	op2 := int((operands & 0xF0) >> 4)

	if m.execF2(opcode, op1, op2) == nil {
		return
	}

	n, i := false, false
	if (opcode & 0x02) == 1 {
		n = true
	}
	if (opcode & 0x01) == 1 {
		i = true
	}

	opcode = opcode >> 2

	if m.execSICF3F4(opcode, operands, n, i) == nil {
		return
	}
}

// execF1 tries to execute opcode as format 1
func (m *Machine) execF1(opcode byte) error {
	switch opcode {
	case FIX: // Not implemented
	case FLOAT: // Not implemented
	case HIO: // Not implemented
	case NORM: // Not implemented
	case SIO: // Not implemented
	case TIO: // Not implemented
	default:
		return fmt.Errorf("command not implemented: %b", opcode)
	}

	return nil
}

// execF2 tries to execute opcode as format 2
func (m *Machine) execF2(opcode byte, op1, op2 int) error {
	switch opcode {
	case ADDR:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		val2, err := m.Reg(op2)

		if err != nil {
			return err
		}

		m.SetReg(op2, val1+val2)
	case CLEAR:
		err := m.SetReg(op1, 0)

		if err != nil {
			return err
		}
	case COMPR:
		r1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		r2, err := m.Reg(op2)

		if err != nil {
			return err
		}

		switch {
		case r1 < r2:
			m.SetSW(LT)
		case r1 == r2:
			m.SetSW(EQ)
		case r1 > r2:
			m.SetSW(GT)
		}
	case DIVR:
		err := m.SetReg(op2, op2/op1)

		if err != nil {
			return err
		}
	case MULR:
		err := m.SetReg(op2, op2*op1)

		if err != nil {
			return err
		}
	case RMO:
		err := m.SetReg(op2, op1)

		if err != nil {
			return err
		}
	case SHIFTL:
		err := m.SetReg(op1, op1<<op2)

		if err != nil {
			return err
		}
	case SHIFTR:
		err := m.SetReg(op1, op1>>op2)

		if err != nil {
			return err
		}
	case SUBR:
		err := m.SetReg(op2, op2-op1)

		if err != nil {
			return err
		}
	case SVC: // Not implemented
	case TIXR:
		if !isRegister(op1) {
			return fmt.Errorf("not a valid register: %d", op1)
		}

		m.SetX(m.X() + 1)

		switch {
		case m.X() < op1:
			m.SetSW(LT)
		case m.X() == op1:
			m.SetSW(EQ)
		case m.X() > op1:
			m.SetSW(GT)
		}
	default:
		return fmt.Errorf("command not implemented: %b", opcode)
	}

	return nil
}

// execSICF3F4 tries to execute opcode either in SIC format, format 3 or format 4
func (m *Machine) execSICF3F4(opcode, operands byte, n, i bool) error {
	var x, b, p, e bool
	var addr, offset, operand, target_address int

	if ((operands & 0x80) >> 7) == 1 {
		x = true
	}

	if ((operands & 0x40) >> 6) == 1 {
		b = true
	}
	if ((operands & 0x20) >> 5) == 1 {
		p = true
	}
	if ((operands & 0x10) >> 4) == 1 {
		e = true
	}

	// Differentiate between formats
	if !(n || i) { // SIC format
		// Address is 0 + lower 7 bits from operands + 8 new bits, which are fetched
		addr = int(binary.BigEndian.Uint32([]byte{operands & 0x7F, m.fetch()}))
	} else if !e { // Format 3
		// Offset is lower 4 bits from operands + 8 new bits, which are fetched
		offset = int(binary.BigEndian.Uint32([]byte{operands & 0x0F, m.fetch()}))
	} else { // Format 4
		// Address is lower 4 bits from operands + 16 new bits, which are fetched
		addr = int(binary.BigEndian.Uint32([]byte{operands & 0x0F, m.fetch(), m.fetch()}))
	}

	// Addressing type
	switch {
	case b && !p: // Base-relative
		target_address = m.B() + offset
	case !b && p: // PC-relative
		target_address = m.PC() + offset
	case !b && !p: // Direct
		target_address = addr
	case b && p: // Disallowed
		return fmt.Errorf("wrong addressing format")
	}

	if x {
		target_address = target_address + m.X()
	}

	switch {
	case n && !i: // Indirect
		l1, err := m.Word(target_address)

		if err != nil {
			return err
		}

		operand, err = m.Word(l1)

		if err != nil {
			return err
		}
	case !n && i: // Immediate
		operand = target_address
	case n && i: // Simple
		var err error
		operand, err = m.Word(target_address)

		if err != nil {
			return err
		}
	case !n && !i: // SIC format
		operand = target_address
	}

	switch opcode {
	case ADD:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetA(m.A() + word)
	case ADDF: // Not implemented
	case AND:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetA(m.A() & word)
	case COMP:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		switch {
		case m.A() < word:
			m.SetSW(LT)
		case m.A() == word:
			m.SetSW(EQ)
		case m.A() > word:
			m.SetSW(GT)
		}
	case COMPF: // Not implemented
	case DIV:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetA(m.A() / word)
	case DIVF: // Not implemented
	case J:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetPC(word)
	case JEQ:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		if m.SW() == EQ {
			m.SetPC(word)
		}
	case JGT:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		if m.SW() == GT {
			m.SetPC(word)
		}
	case JLT:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		if m.SW() == LT {
			m.SetPC(word)
		}
	case JSUB:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetL(m.PC())
		m.SetPC(word)
	case LDA:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetA(word)
	case LDB:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetB(word)
	case LDCH:
		low, err := m.Byte(operand)

		if err != nil {
			return err
		}

		m.SetALow(low)
	case LDF:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetF(word)
	case LDL:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetL(word)
	case LDS:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetS(word)
	case LDT:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetT(word)
	case LDX:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetX(word)
	case LPS: // Not implemented
	case MUL:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetA(m.A() * word)
	case MULF: // Not implemented
	case OR:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetA(m.A() | word)
	case RD:
		devno, err := m.Byte(operand)

		if err != nil {
			return err
		}

		dev := m.Dev(devno)
		text, err := dev.Read()

		if err != nil {
			return err
		}

		m.SetALow(text)
	case RSUB:
		m.SetPC(m.L())
	case SSK: // Not implemented
	case STA:
		err := m.SetWord(operand, m.A())

		if err != nil {
			return err
		}
	case STB:
		err := m.SetWord(operand, m.B())

		if err != nil {
			return err
		}
	case STCH:
		err := m.SetByte(operand, m.ALow())

		if err != nil {
			return err
		}
	case STF:
		err := m.SetWord(operand, m.F())

		if err != nil {
			return err
		}
	case STI: // Not implemented
	case STL:
		err := m.SetWord(operand, m.L())

		if err != nil {
			return err
		}
	case STS:
		err := m.SetWord(operand, m.S())

		if err != nil {
			return err
		}
	case STSW:
		err := m.SetWord(operand, m.SW())

		if err != nil {
			return err
		}
	case STT:
		err := m.SetWord(operand, m.T())

		if err != nil {
			return err
		}
	case STX:
		err := m.SetWord(operand, m.X())

		if err != nil {
			return err
		}
	case SUB:
		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		m.SetA(m.A() - word)
	case SUBF: // Not implemented
	case TD:
		devno, err := m.Byte(operand)

		if err != nil {
			return err
		}

		dev := m.Dev(devno)
		err = dev.Test()

		if err != nil {
			return err
		}
	case TIX:
		m.SetX(m.X() + 1)

		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		switch {
		case m.X() < word:
			m.SetSW(LT)
		case m.X() == word:
			m.SetSW(EQ)
		case m.X() > word:
			m.SetSW(GT)
		}
	case WD:
		devno, err := m.Byte(operand)

		if err != nil {
			return err
		}

		dev := m.Dev(devno)
		err = dev.Write(m.ALow())

		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("command not implemented: %b", opcode)
	}

	return nil
}
