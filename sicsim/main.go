package main

/*
Global TODO:
	- Add error reporting / handling
	- Add unit tests
*/

import (
	"fmt"
	"sicsim/sim"
)

func main() {
	// Temporary testing
	m := sim.Machine{}
	m.New()

	m.Regs.SetA(10)
	m.Regs.SetReg(1, 128)

	fmt.Println(m.Registers())
}
