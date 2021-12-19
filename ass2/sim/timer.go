package sicsim

import (
	"fmt"
	"time"
)

// Start starts executing commands from memory
func (m *Machine) Start() {
	m.ticker = time.NewTicker(m.tick) // Always reset the ticker

	go func() {
		for range m.ticker.C {
			if !m.Halted() {
				m.Execute()
			} else {
				m.Stop()
				fmt.Printf("\n-- Done (executed all instructions) --\n")
				fmt.Print("> ")
				return
			}
		}
	}()
}

// Stop stops executing commands and stops the machine's ticker
func (m *Machine) Stop() {
	m.ticker.Stop()
	m.ticker = nil
}

// IsRunning returns the status of the machine's ticker
func (m *Machine) IsRunning() bool {
	return m.ticker != nil
}

// Speed returns the current ticker speed
func (m *Machine) Speed() string {
	return m.tick.String()
}

func (m *Machine) SetSpeed(kHz int) {
	// 1 kHz == 1 ms
	m.tick = time.Duration(kHz * int(time.Millisecond))
}
