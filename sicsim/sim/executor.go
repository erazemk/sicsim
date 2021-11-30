package sim

import "encoding/binary"

// Fetch returns a byte from m[PC] and increments PC
func (m *Machine) Fetch() byte {
	addr := m.PC()
	val := m.mem[addr]
	m.SetPC(addr + 1)
	return val
}

// Execute tries to execute fetched operation
func (m *Machine) Execute() {
	opcode := m.Fetch()

	if m.execF1(opcode) {
		return
	}

	operands := m.Fetch()
	op1 := int(operands & 0x0F)
	op2 := int((operands & 0xF0) >> 4)

	if m.execF2(opcode, op1, op2) {
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

	if m.execSICF3F4(opcode, operands, n, i) {
		return
	}
}

// execF1 tries to execute opcode as format 1
func (m *Machine) execF1(opcode byte) bool {
	switch opcode {
	case FIX:
		m.SetA(int(m.F()))
	case FLOAT:
		m.SetF(float64(m.A()))
	case HIO: // TODO: ustavi I/O - haltio(A)
	case NORM: // TODO: normalizacija floata (popravljanje mantise) - F <- norm(F)
	case SIO: // TODO: štartaj I/O - startio(A, S)
	case TIO: // TODO: testiraj I/O - testio(A)
	default:
		return false
	}

	return false
}

// execF2 tries to execute opcode as format 2
func (m *Machine) execF2(opcode byte, op1, op2 int) bool {
	switch opcode {
	case ADDR:
		m.SetReg(op2, m.Reg(op1)+m.Reg(op2))
	case CLEAR:
		m.SetReg(op1, 0)
	case COMPR:
		r1 := m.Reg(op1)
		r2 := m.Reg(op2)

		switch {
		case r1 < r2:
			m.SetSW(LT)
		case r1 == r2:
			m.SetSW(EQ)
		case r1 > r2:
			m.SetSW(GT)
		}
	case DIVR:
		m.SetReg(op2, op2/op1)
	case MULR:
		m.SetReg(op2, op2*op1)
	case RMO:
		m.SetReg(op2, op1)
	case SHIFTL:
		m.SetReg(op1, op1<<op2)
	case SHIFTR:
		m.SetReg(op1, op1>>op2)
	case SUBR:
		m.SetReg(op2, op2-op1)
	case SVC:
		// TODO: throwi error not implemented
	case TIXR:
		m.SetX(m.X() + 1)

		switch {
		case m.X() < op1:
			m.SetSW(LT)
		case m.X() == op2:
			m.SetSW(EQ)
		case m.X() > op2:
			m.SetSW(GT)
		}
	default:
		return false
	}

	return false
}

func (m *Machine) execSICF3F4(opcode, operands byte, n, i bool) bool {
	x := false
	b, p, e := false, false, false
	var address, offset uint32 = 0, 0

	if ((operands & 0x80) >> 7) == 1 {
		x = true
	}

	if !(n || i) {
		// Address is 0 + lower 7 bits from operands + 8 new bits, which are fetched
		address = binary.BigEndian.Uint32([]byte{operands & 0x7F, m.Fetch()})
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

	// Differentiate between Format 3 and 4
	if !e { // Format 3
		// Offset is lower 4 bits from operands + 8 new bits, which are fetched
		offset = binary.BigEndian.Uint32([]byte{operands & 0x0F, m.Fetch()})
	} else { // Format 4
		// Address is lower 4 bits from operands + 16 new bits, which are fetched
		address = binary.BigEndian.Uint32([]byte{operands & 0x0F, m.Fetch(), m.Fetch()})
	}

	switch opcode {
	case ADD:
		m.SetA(m.A() + m.Word(address))
	case ADDF:
		m.SetF(m.F() + m.Float(address))
	case AND:
		m.SetA(m.A() & m.Word(address))
	case COMP:
		switch {
		case m.A() < m.Word(address):
			m.SetSW(LT)
		case m.A() == m.Word(address):
			m.SetSW(EQ)
		case m.A() > m.Word(address):
			m.SetSW(GT)
		}
	case COMPF:
		switch {
		case m.F() < m.Float(address):
			m.SetSW(LT)
		case m.F() == m.Float(address):
			m.SetSW(EQ)
		case m.F() > m.Float(address):
			m.SetSW(GT)
		}
	case DIV:
		m.SetA(m.A() / m.Word(address))
	case DIVF:
		m.SetF(m.F() / m.Float(address))
	case J:
		m.SetPC(m.Word(address))
	case JEQ:
		if m.SW() == EQ {
			m.SetPC(m.Word(address))
		}
	case JGT:
		if m.SW() == GT {
			m.SetPC(m.Word(address))
		}
	case JLT:
		if m.SW() == LT {
			m.SetPC(m.Word(address))
		}
	case JSUB:
		m.SetL(m.Word(uint32(m.PC()))) // TODO: Conversion shouldn't need to be used
		m.SetPC(m.Word(address))
	case LDA:
		m.SetA(m.Word(address))
	case LDB:
		m.SetB(m.Word(address))
	case LDCH:
		// TODO: Fixi da bo dejansko spremenilo samo spodnjih 8 bitov
		ch := binary.BigEndian.Uint32([]byte{0x00, 0x00, 0x00, m.Byte(address)})
		m.SetA(int(ch))
	case LDF:
		m.SetF(m.Float(address))
	case LDL:
		m.SetL(m.Word(address))
	case LDS:
		m.SetS(m.Word(address))
	case LDT:
		m.SetT(m.Word(address))
	case LDX:
		m.SetX(m.Word(address))
	case LPS:
		// TODO: throwi error not implemented
	case MUL:
		m.SetA(m.A() * m.Word(address))
	case MULF:
		m.SetF(m.F() * m.Float(address))
	case OR:
		m.SetA(m.A() | m.Word(address))
	case RD:
		// TODO: Dodaj branje iz naprave
	case RSUB:
		m.SetPC(m.L())
	case SSK:
		// TODO: throwi error not implemented
	case STA:
		m.SetWord(address, m.A())
	case STB:
		m.SetWord(address, m.B())
	case STCH:
		m.SetWord(address, m.A()&0x000F)
	case STF:
		m.SetFloat(address, m.F())
	case STI:
		// TODO: throwi error not implemented
	case STL:
		m.SetWord(address, m.L())
	case STS:
		m.SetWord(address, m.S())
	case STSW:
		m.SetWord(address, m.SW())
	case STT:
		m.SetWord(address, m.T())
	case STX:
		m.SetWord(address, m.X())
	case SUB:
		m.SetA(m.A() - m.Word(address))
	case SUBF:
		m.SetF(m.F() - m.Float(address))
	case TD:
		// TODO: Implementiri test device
	case TIX:
		m.SetX(m.X() + 1)
		switch {
		case m.X() < m.Word(address):
			m.SetSW(LT)
		case m.X() == m.Word(address):
			m.SetSW(EQ)
		case m.X() > m.Word(address):
			m.SetSW(GT)
		}
	case WD:
		// TODO: Implementiri write device
		m.Devs[m.Word(address)].Write(byte(m.A()))
	default:
		return false
	}

	return false
}
