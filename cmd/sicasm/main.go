package main

import (
	"fmt"

	"git.sr.ht/~erazemk/sicsim/asm"
)

func main() {
	input := "    START 42\n    END zacetek"
	parser := asm.NewParser()
	code := parser.Parse(input)
	fmt.Println(code)
}
