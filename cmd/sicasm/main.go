package main

import (
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~erazemk/sicsim/asm"
	opt "github.com/pborman/getopt/v2"
)

func main() {
	// Flags
	debugFlag := opt.BoolLong("debug", 'd', "Show debug info")
	lstFlag := opt.BoolLong("lst", 'l', "Pretty-print object and assembly code")
	helpFlag := opt.BoolLong("help", 'h', "Show this text")
	outputFlag := opt.StringLong("output", 'o', "", "Generated object file path", "/path/to/file.obj")
	opt.SetParameters("/path/to/file.asm")
	opt.Parse()

	if *helpFlag {
		opt.Usage()
		os.Exit(0)
	}

	asm.SetDebug(*debugFlag)
	asm.SetPrettyPrint(*lstFlag)

	inputFile := opt.Arg(0)
	if inputFile == "" {
		fmt.Printf("No input file provided!\n\n")
		opt.Usage()
		os.Exit(1)
	}

	if *debugFlag {
		fmt.Println("Input file: ", inputFile)
	}

	code := asm.NewCode()

	// First pass: parse code
	if err := code.ParseFile(inputFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Second pass: replace variables with their values
	code.ResolveSymbols()

	// Use same file path as input by default
	outputFile := inputFile[:strings.LastIndex(inputFile, ".")]

	// If output file path is specified, use it
	if *outputFlag != "" {
		outputFile = *outputFlag

		// Remove extension if present
		outputFile = outputFile[:strings.LastIndex(inputFile, ".")]
	}

	outputFile += ".obj"

	if *debugFlag {
		fmt.Println("Output file: ", outputFile)
	}

	if err := code.CreateObjectFile(outputFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *lstFlag {
		code.PrettyPrint()
	}
}
