package main

import (
	sim "ass2/sim"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pborman/getopt/v2"
)

func main() {
	// Flags
	objFileFlag := getopt.StringLong("object", 'o', "", "Object file to run")
	interactiveFlag := getopt.BoolLong("interactive", 'i', "Run interactive programs")
	debugFlag := getopt.BoolLong("debug", 'd', "Enable debug output")
	helpFlag := getopt.BoolLong("help", 'h', "Show this text")
	getopt.Parse()

	if *helpFlag {
		help()
		os.Exit(0)
	}

	sim.SetDebug(*debugFlag)

	// Create a new machine
	var m sim.Machine
	m.New()

	if *objFileFlag == "" {
		fmt.Println("No object file provided!")
		help()
		os.Exit(1)
	} else {
		if err := m.ParseObj(*objFileFlag); err != nil {
			fmt.Println(err)
		}
	}

	if !*interactiveFlag {
		// Clear screen
		scr := exec.Command("clear")
		scr.Stdout = os.Stdout
		scr.Run()

		header()
		fmt.Println("(REPL mode)")
		replHelp()
		repl(m)
	} else {
		m.Start()
		sc := bufio.NewScanner(os.Stdin)

		for sc.Scan() {
			if !m.Halted() {
				for _, char := range sc.Bytes() {
					m.WriteDevice(1, char)
				}
			} else {
				os.Exit(0)
			}
		}
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
				m.Start()
				fmt.Println("Started automatic execution")
			} else {
				fmt.Println("Finished executing program, stop trying to break things")
			}
		case "end", "et":
			if !m.Halted() {
				m.Stop()
				fmt.Println("Stopped automatic execution")
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
	fmt.Println()
	fmt.Println("Usage: sicsim (-h) (-d) [-i /path/to/file.obj]")
	fmt.Println()
	fmt.Println("  -d, --debug                        Print debug info during execution")
	fmt.Println("  -h, --help                         Print this text")
	fmt.Println("  -i, --input [/path/to/file.obj]    Object file to execute")
	fmt.Println()
}

func replHelp() {
	fmt.Println()
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
	fmt.Println()
}

func header() {
	fmt.Printf(
		"███████ ██  ██████ ███████ ██ ███    ███\n" +
			"██      ██ ██      ██      ██ ████  ████\n" +
			"███████ ██ ██      ███████ ██ ██ ████ ██\n" +
			"     ██ ██ ██           ██ ██ ██  ██  ██\n" +
			"███████ ██  ██████ ███████ ██ ██      ██\n\n")
}
