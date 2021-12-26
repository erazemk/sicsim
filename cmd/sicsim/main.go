package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"git.sr.ht/~erazemk/sicsim/sim"
	"github.com/pborman/getopt/v2"
)

func main() {
	// Flags
	debugFlag := getopt.BoolLong("debug", 'd', "Enable debug output")
	helpFlag := getopt.BoolLong("help", 'h', "Show this text")
	interactiveFlag := getopt.BoolLong("non-repl", 'n', "Automatically run programs (non-REPL mode)")
	getopt.Parse()

	if *helpFlag {
		help()
		os.Exit(0)
	}

	sim.SetDebug(*debugFlag)

	objFile := getopt.Arg(0)
	if objFile == "" {
		fmt.Printf("No object file provided!\n\n")
		help()
		os.Exit(1)
	}

	// Clear screen if running in REPL mode (overwritten by debug mode)
	if !*interactiveFlag {
		scr := exec.Command("clear")
		scr.Stdout = os.Stdout
		scr.Run()
	}

	// Create a new machine
	var m sim.Machine
	m.New()
	m.SetInteractive(*interactiveFlag)

	if err := m.ParseObjFile(objFile); err != nil {
		fmt.Println(err)
	}

	if !*interactiveFlag {
		header()
		fmt.Println("(REPL mode)")
		replHelp()
		repl(m)
	} else {
		m.Start()
	}
}

// Runs the simulator in REPL mode
func repl(m sim.Machine) {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for sc.Scan() {
		switch text := strings.Split(sc.Text(), " "); text[0] {
		case "regs", "r":
			fmt.Println(m.Regs())
		case "mem", "m":
			low, err := strconv.Atoi(text[1])

			if err != nil {
				panic(err)
			}

			high, err := strconv.Atoi(text[2])

			if err != nil {
				panic(err)
			}

			fmt.Println(m.Mem(low, high))
		case "exec", "e":
			if !m.Halted() {
				m.Execute()
			} else {
				fmt.Println("Finished executing program, stop trying to break things")
			}
		case "step", "s":
			if !m.Halted() {
				m.Execute()
				fmt.Println(m.Regs())
			} else {
				fmt.Println("Finished executing program, stop trying to break things")
			}
		case "word", "w":
			addr, err := strconv.Atoi(text[1])

			if err != nil {
				panic(err)
			}

			word, err := m.Word(addr)

			if err != nil {
				panic(err)
			}

			fmt.Printf("%02X\n", word)
		case "byte", "b":
			addr, err := strconv.Atoi(text[1])

			if err != nil {
				panic(err)
			}

			byt, err := m.Byte(addr)

			if err != nil {
				panic(err)
			}

			fmt.Printf("%02X\n", byt)
		case "setreg", "sr":
			no, err := strconv.Atoi(text[1])

			if err != nil {
				switch text[1] {
				case "a", "A":
					no = 0
				case "x", "X":
					no = 1
				case "l", "L":
					no = 2
				case "b", "B":
					no = 3
				case "s", "S":
					no = 4
				case "t", "T":
					no = 5
				case "f", "F":
					no = 6
				case "pc", "PC":
					no = 8
				case "sw", "SW":
					no = 9
				default:
					fmt.Printf("Invalid register: %s\n", text[2])
					continue
				}
			}

			val, err := strconv.Atoi(text[2])

			if err != nil {
				panic(err)
			}

			if err := m.SetReg(no, val); err != nil {
				panic(err)
			}
		case "begin", "bt":
			if !m.Halted() {
				fmt.Println("Started automatic execution")
				m.Start()
			} else {
				fmt.Println("Finished executing program, stop trying to break things")
			}
		case "end", "et":
			if !m.Halted() {
				fmt.Println("Stopped automatic execution")
				m.Stop()
			} else {
				fmt.Println("Finished executing program, stop trying to break things")
			}
		default:
			replHelp()
		}

		fmt.Print("> ")
	}
}

func help() {
	fmt.Println("Usage: sicsim (-dhn) /path/to/file.obj")
	fmt.Println()
	fmt.Println("  -d, --debug       Print debug info during execution")
	fmt.Println("  -h, --help        Print this text")
	fmt.Println("  -n, --non-repl    Automatically run programs (non-REPL mode)")
	fmt.Println()
}

func replHelp() {
	fmt.Println("Usage: [command] (options)")
	fmt.Println()
	fmt.Println("  Memory and registers:")
	fmt.Println("    b, byte [addr]           Returns the byte at memory[addr]")
	fmt.Println("    w, word [addr]           Returns the word at memory[addr]")
	fmt.Println("    m, mem [low] [high]      Prints memory contents from low to high address")
	fmt.Println("    r, regs                  Prints register values")
	fmt.Println("    sr, setreg [no] [val]    Sets the register [no] to [val]")
	fmt.Println()
	fmt.Println("  Instructions:")
	fmt.Println("    e, exec                  Executes the next instruction")
	fmt.Println("    s, step                  Executes the next instruction and prints register values")
	fmt.Println("    bt, begin                Starts automatically executing instructions")
	fmt.Println("    et, end                  Stops automatically executing instructions")
}

func header() {
	fmt.Printf(
		"███████ ██  ██████ ███████ ██ ███    ███\n" +
			"██      ██ ██      ██      ██ ████  ████\n" +
			"███████ ██ ██      ███████ ██ ██ ████ ██\n" +
			"     ██ ██ ██           ██ ██ ██  ██  ██\n" +
			"███████ ██  ██████ ███████ ██ ██      ██\n\n")
}
