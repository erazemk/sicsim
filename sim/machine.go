package sim

import (
	"fmt"
	"log"
	"time"
)

const MAX_ADDRESS = 1048576

// Stack for JSUB and RSUB instructions
type stack []int

type Machine struct {
	regs        registers
	mem         [MAX_ADDRESS + 1]byte
	devs        [256](*device)
	stack       stack
	tick        time.Duration
	ticker      *time.Ticker
	halted      bool
	interactive bool
}

// New creates a new machine
func (m *Machine) New() {
	m.NewDevice(0) // stdin
	m.NewDevice(1) // stdout
	m.NewDevice(2) // stderr
	m.stack = stack{}
	m.tick = time.Millisecond // Default clock duration
	m.ticker = nil

	if debug {
		log.Println("Created a new machine")
	}
}

// Returns true if execution has halted
func (m *Machine) Halted() bool {
	return m.halted
}

func (m *Machine) SetInteractive(interactive bool) {
	m.interactive = interactive
}

func (m *Machine) TestDevice(id byte) bool {
	if m.devs[id] != nil {
		return m.devs[id].test()
	}

	m.NewDevice(id)
	return m.devs[id].test()
}

func (m *Machine) ReadDevice(id byte) (byte, error) {
	if m.devs[id] != nil {
		return m.devs[id].read()
	}

	m.NewDevice(id)
	return m.devs[id].read()
}

func (m *Machine) WriteDevice(id, val byte) error {
	if m.devs[id] != nil {
		return m.devs[id].write(val)
	}

	m.NewDevice(id)
	return m.devs[id].write(val)
}

func (m *Machine) NewDevice(id byte) error {
	if m.devs[id] != nil {
		return fmt.Errorf("device '%s' already exists", m.devs[id].name)
	}

	dev, err := newDevice(id)

	if err != nil {
		return err
	}

	m.devs[id] = dev
	return nil
}

func (m *Machine) push(val int) {
	m.stack = append(m.stack, val)
}

func (m *Machine) pop() int {
	st := m.stack[len(m.stack)-1]
	m.stack = m.stack[:len(m.stack)-1]
	return st
}
