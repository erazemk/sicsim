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

	if debug {
		fmt.Printf("Opcode [b]: 0x%02X\n", opcode)
	}

	if m.execF1(opcode) == nil {
		return
	}

	operands := m.fetch()
	r1 := int((operands & 0xF0) >> 4)
	r2 := int(operands & 0x0F)

	if m.execF2(opcode, r1, r2) == nil {
		return
	}

	n, i := false, false
	if (opcode & 0x02) == 0x02 {
		n = true
	}
	if (opcode & 0x01) == 0x01 {
		i = true
	}

	// opcode is only 6-bit for SIC/F3/F4 instruction formats
	opcode = opcode & 0xFC

	if debug {
		fmt.Printf("N=%v, I=%v\n", n, i)
		fmt.Printf("Opcode [a]: 0x%02X\n", opcode)
	}

	// operands = xbpe bits + part of address/offset
	if m.execSICF3F4(opcode, operands, n, i) == nil {
		return
	}

	// TODO: Execute() should return error if none of the formats were correct
}

// execF1 tries to execute opcode as format 1
func (m *Machine) execF1(opcode byte) error {
	if debug {
		fmt.Println("[Format 1]")
	}

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
	if debug {
		fmt.Println("[Format 2]")
	}

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

		if err := m.SetReg(op2, val2+val1); err != nil {
			return err
		}
	case CLEAR:
		if err := m.SetReg(op1, 0); err != nil {
			return err
		}
	case COMPR:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		val2, err := m.Reg(op2)

		if err != nil {
			return err
		}

		switch {
		case val1 < val2:
			m.SetSW(LT)
		case val1 == val2:
			m.SetSW(EQ)
		case val1 > val2:
			m.SetSW(GT)
		}
	case DIVR:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		val2, err := m.Reg(op2)

		if err != nil {
			return err
		}

		if err := m.SetReg(op2, val2/val1); err != nil {
			return err
		}
	case MULR:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		val2, err := m.Reg(op2)

		if err != nil {
			return err
		}

		if err := m.SetReg(op2, val2*val1); err != nil {
			return err
		}
	case RMO:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		if err := m.SetReg(op2, val1); err != nil {
			return err
		}
	case SHIFTL:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		if err := m.SetReg(op1, val1<<op2); err != nil {
			return err
		}
	case SHIFTR:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		if err := m.SetReg(op1, val1>>op2); err != nil {
			return err
		}
	case SUBR:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		val2, err := m.Reg(op2)

		if err != nil {
			return err
		}

		if err := m.SetReg(op2, val2-val1); err != nil {
			return err
		}
	case SVC: // Not implemented
	case TIXR:
		val1, err := m.Reg(op1)

		if err != nil {
			return err
		}

		m.SetX(m.X() + 1)
		rX := m.X()

		switch {
		case rX < val1:
			m.SetSW(LT)
		case rX == val1:
			m.SetSW(EQ)
		case rX > val1:
			m.SetSW(GT)
		}
	default:
		return fmt.Errorf("command not implemented: %b", opcode)
	}

	return nil
}

// execSICF3F4 tries to execute opcode either in SIC format, format 3 or format 4
func (m *Machine) execSICF3F4(opcode, operands byte, n, i bool) error {
	if debug {
		fmt.Println("[Format 3]")
	}

	var x, b, p, e bool
	var operand, target_address int

	if (operands & 0x80) == 0x80 {
		x = true
	}

	if (operands & 0x40) == 0x40 {
		b = true
	}

	if (operands & 0x20) == 0x20 {
		p = true
	}

	if (operands & 0x10) == 0x10 {
		e = true
	}

	// Differentiate between formats
	if !(n || i) { // SIC format
		// addr = lower 7 bits from operands + 8 fetched bits
		target_address = int(binary.BigEndian.Uint32([]byte{0, 0, operands & 0x7F, m.fetch()}))
	} else if !e { // Format 3
		// offset = lower 4 bits from operands + 8 fetched bits
		target_address = int(binary.BigEndian.Uint32([]byte{0, 0, operands & 0x0F, m.fetch()}))
	} else { // Format 4
		// addr = lower 4 bits from operands + 16 fetched bits
		target_address = int(binary.BigEndian.Uint32([]byte{0, operands & 0x0F, m.fetch(), m.fetch()}))
	}

	if debug {
		fmt.Printf("Bits: N=%v, I=%v, X=%v, B=%v, P=%v, E=%v\n", n, i, x, b, p, e)
		fmt.Printf("TA[b]: 0x%06X\n", target_address)
	}

	// Addressing type
	if n || i { // F3 / F4
		switch {
		case b && !p: // Base-relative
			target_address += m.B()
		case !b && p: // PC-relative
			target_address += m.PC()
		case !b && !p: // Direct
			// No extra action needed
		case b && p: // Disallowed
			return fmt.Errorf("wrong addressing format")
		}
	}

	if x {
		if n && i { // Simple addressing
			target_address += m.X()
		} else {
			return fmt.Errorf("invalid addressing")
		}
	}

	switch {
	case n && !i: // Indirect
		lvl1, err := m.Word(target_address)

		if err != nil {
			return err
		}

		operand, err = m.Word(lvl1)

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

	if debug {
		fmt.Printf("Operand: 0x%06X\n", operand)
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

		rA := m.A()

		switch {
		case rA < word:
			m.SetSW(LT)
		case rA == word:
			m.SetSW(EQ)
		case rA > word:
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

		char, err := m.ReadDevice(devno)

		if err != nil {
			return err
		}

		m.SetALow(char)
	case RSUB:
		m.SetPC(m.L())
	case SSK: // Not implemented
	case STA:
		if err := m.SetWord(operand, m.A()); err != nil {
			return err
		}
	case STB:
		if err := m.SetWord(operand, m.B()); err != nil {
			return err
		}
	case STCH:
		if err := m.SetByte(operand, m.ALow()); err != nil {
			return err
		}
	case STF:
		if err := m.SetWord(operand, m.F()); err != nil {
			return err
		}
	case STI: // Not implemented
	case STL:
		if err := m.SetWord(operand, m.L()); err != nil {
			return err
		}
	case STS:
		if err := m.SetWord(operand, m.S()); err != nil {
			return err
		}
	case STSW:
		if err := m.SetWord(operand, m.SW()); err != nil {
			return err
		}
	case STT:
		if err := m.SetWord(operand, m.T()); err != nil {
			return err
		}
	case STX:
		if err := m.SetWord(operand, m.X()); err != nil {
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

		if !m.TestDevice(devno) {
			return fmt.Errorf("device test failed")
		}
	case TIX:
		m.SetX(m.X() + 1)

		word, err := m.Word(operand)

		if err != nil {
			return err
		}

		rX := m.X()

		switch {
		case rX < word:
			m.SetSW(LT)
		case rX == word:
			m.SetSW(EQ)
		case rX > word:
			m.SetSW(GT)
		}
	case WD:
		devno, err := m.Byte(operand)

		if err != nil {
			return err
		}

		if err = m.WriteDevice(devno, m.ALow()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("command not implemented: %b", opcode)
	}

	return nil
}
