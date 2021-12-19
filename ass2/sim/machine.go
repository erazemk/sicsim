package sicsim

import (
	"fmt"
	"log"
	"time"
)

const MAX_ADDRESS = 1048576

type Machine struct {
	regs   registers
	mem    [MAX_ADDRESS + 1]byte
	devs   [256](*device)
	tick   time.Duration
	ticker *time.Ticker
	halted bool
}

// New creates a new machine
func (m *Machine) New() {
	m.NewDevice(0)            // stdin
	m.NewDevice(1)            // stdout
	m.NewDevice(2)            // stderr
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

func (m *Machine) TestDevice(id byte) bool {
	return m.devs[id].test()
}

func (m *Machine) ReadDevice(id byte) (byte, error) {
	return m.devs[id].read()
}

func (m *Machine) WriteDevice(id, val byte) error {
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
