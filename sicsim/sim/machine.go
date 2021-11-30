package sim

import (
	"log"
)

const MAX_ADDRESS = 1048576

type Machine struct {
	regs registers
	mem  [MAX_ADDRESS + 1]byte
	devs [256](*device)
}

// New creates a new machine
func (m *Machine) New() {
	m.devs[0].New(0) // stdin
	m.devs[1].New(1) // stdout
	m.devs[2].New(2) // stderr

	if debug {
		log.Println("Created a new machine")
	}
}
