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

	m.SetA(10)
	m.SetReg(1, 128)

	fmt.Println(m.Registers())
}
