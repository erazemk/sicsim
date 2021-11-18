package sim

// TODO: Maybe convert these into constants so they can be referenced as types
type Registers struct {
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

// TODO: do something with SW values
const (
	LT = 0x00
	EQ = 0x40
	GT = 0x80
)

// Reg returns the value of register reg
func (r Registers) Reg(reg int) int {
	switch reg {
	case 0:
		return r.a
	case 1:
		return r.x
	case 2:
		return r.l
	case 3:
		return r.b
	case 4:
		return r.s
	case 5:
		return r.t
	case 6:
		return r.f
	case 8:
		return r.pc
	case 9:
		return r.sw
	}

	return -1
}

// SetReg sets the value of register reg
func (r *Registers) SetReg(reg, val int) {
	if reg >= 0 && reg <= 6 && isWord(val) {
		switch reg {
		case 0:
			r.a = val
		case 1:
			r.x = val
		case 2:
			r.l = val
		case 3:
			r.b = val
		case 4:
			r.s = val
		case 5:
			r.t = val
		case 6:
			r.f = val
		case 8:
			r.pc = val
		case 9:
			r.sw = val
		}
	}
}

// A returns the value of
func (r Registers) A() int {
	return r.a
}

// X returns the value of the X register
func (r Registers) X() int {
	return r.x
}

// L returns the value of the L register
func (r Registers) L() int {
	return r.l
}

// B returns the value of the B register
func (r Registers) B() int {
	return r.b
}

// S returns the value of the S register
func (r Registers) S() int {
	return r.s
}

// T returns the value of the T register
func (r Registers) T() int {
	return r.t
}

// F returns the value of the F register
func (r Registers) F() int {
	return r.f
}

// PC returns the value of the PC register
func (r Registers) PC() int {
	return r.pc
}

// SW returns the value of the SW register
func (r Registers) SW() int {
	return r.sw
}

// SetA sets the value of the A register
func (r *Registers) SetA(val int) {
	if isWord(val) {
		r.a = val
	}
}

// SetX sets the value of the X register
func (r *Registers) SetX(val int) {
	if isWord(val) {
		r.x = val
	}
}

// SetL sets the value of the L register
func (r *Registers) SetL(val int) {
	if isWord(val) {
		r.l = val
	}
}

// SetB sets the value of the B register
func (r *Registers) SetB(val int) {
	if isWord(val) {
		r.b = val
	}
}

// SetS sets the value of the S register
func (r *Registers) SetS(val int) {
	if isWord(val) {
		r.s = val
	}
}

// SetT sets the value of the T register
func (r *Registers) SetT(val int) {
	if isWord(val) {
		r.t = val
	}
}

// SetF sets the value of the F register
func (r *Registers) SetF(val int) {
	if isWord(val) {
		r.f = val
	}
}

// SetPC sets the value of the PC register
func (r *Registers) SetPC(val int) {
	if isWord(val) {
		r.pc = val
	}
}

// SetSW sets the value of the SW register
func (r *Registers) SetSW(val int) {
	if isWord(val) {
		r.sw = val
	}
}
