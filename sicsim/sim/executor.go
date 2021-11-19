package sim

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
	op := 0
	n, i := false, false

	if m.execF1(opcode) {
		return
	}

	if m.execF2(opcode, op) {
		return
	}

	if m.execSICF3F4(opcode, n, i, op) {
		return
	}
}

// execF1 tries to execute opcode as format 1
func (m *Machine) execF1(opcode byte) bool {
	switch opcode {
	case FIX:
	case FLOAT:
	case HIO:
	case NORM:
	case SIO:
	case TIO:
	default:
		return false
	}

	return false
}

// execF2 tries to execute opcode as format 2
func (m *Machine) execF2(opcode byte, op int) bool {
	switch opcode {
	case ADDR:
	case CLEAR:
	case COMPR:
	case DIVR:
	case MULR:
	case RMO:
	case SHIFTL:
	case SHIFTR:
	case SUBR:
	case SVC:
	case TIXR:
	default:
		return false
	}

	return false
}

func (m *Machine) execSICF3F4(opcode byte, n, i bool, op int) bool {
	switch opcode {
	case ADD:
	case ADDF:
	case AND:
	case COMP:
	case COMPF:
	case DIV:
	case DIVF:
	case J:
	case JEQ:
	case JGT:
	case JLT:
	case JSUB:
	case LDA:
	case LDB:
	case LDCH:
	case LDF:
	case LDL:
	case LDS:
	case LDT:
	case LDX:
	case LPS:
	case MUL:
	case MULF:
	case OR:
	case RD:
	case RSUB:
	case SSK:
	case STA:
	case STB:
	case STCH:
	case STF:
	case STI:
	case STL:
	case STS:
	case STSW:
	case STT:
	case STX:
	case SUB:
	case SUBF:
	case TD:
	case TIX:
	case WD:
	default:
		return false
	}

	return false
}
