package sim

import (
	"log"
)

const MAX_ADDRESS = 1048576

type Machine struct {
	regs registers
	mem  [MAX_ADDRESS + 1]byte
	Devs [256]Device
}

// New creates a new machine
func (m *Machine) New() {
	m.Devs[0].New(0) // stdin
	m.Devs[1].New(1) // stdout
	m.Devs[2].New(2) // stderr

	if debug {
		log.Println("Created a new machine")
	}
}

// Device returns a device by its number
func (m *Machine) Device(num byte) Device {
	if isDevice(num) {
		return m.Devs[num]
	}

	return Device{}
}
