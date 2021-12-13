package sicsim

import "time"

// Start starts executing commands from memory
func (m *Machine) Start() {
	m.ticker = time.NewTicker(m.tick) // Always reset the ticker

	go func() {
		for range m.ticker.C {
			m.Execute()
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
func (m *Machine) Speed() time.Duration {
	return m.tick
}

func (m *Machine) SetSpeed(kHz int) {
	// 1 kHz == 1 ms
	m.tick = time.Duration(kHz * int(time.Millisecond))
}
