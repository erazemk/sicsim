package asm

import (
	"fmt"
	"os"
)

type Code struct {
	name         string
	startaddr    int
	lc           int
	length       int
	brelative    bool
	pcstartaddr  int
	instructions []Node
	symtab       map[string]interface{}
}

// NewCode returns a new instance of Code
func NewCode() *Code {
	return &Code{
		symtab: make(map[string]interface{}),
	}
}

// ResolveSymbols replaces symbols with operands for nodes that don't already have operands
func (c *Code) ResolveSymbols() {
	// Go through each instruction and replace symbols with proper operands
	for i, node := range c.instructions {
		if node.symbol != "" {
			// Replace symbol with operand
			node.operand = c.symtab[node.symbol].(int)
			c.instructions[i] = node

			if debug {
				fmt.Printf("Resolved symbol '%s' for %s: %d\n", node.symbol, node.mnemonic, node.operand)
			}
		}
	}
}

// CreateObjectFile writes all the necessary records to the specified file
func (c *Code) CreateObjectFile(name string) error {
	lc := c.startaddr

	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	defer file.Close()

	if debug {
		fmt.Printf("Opened future object file '%s': %v\n", name, file)
	}

	// Write header record
	header := fmt.Sprintf("H%-6s%s%s\n", c.name, Word(c.startaddr), Word(c.length))
	nh, err := file.WriteString(header)
	if err != nil {
		return fmt.Errorf("failed to write header record: %w", err)
	}

	if debug {
		fmt.Printf("Wrote header record (%d bytes): %s", nh, header)
	}

	// Write text records
	for _, node := range c.instructions {
		lc += node.length
		bytes := node.Bytes()

		// Skip empty nodes
		if len(bytes) == 0 {
			continue
		}

		text := fmt.Sprintf("T%s%02X%s\n", Word(node.lc), len(bytes)/2, bytes)
		nt, err := file.WriteString(text)
		if err != nil {
			return fmt.Errorf("failed to write text record: %w", err)
		}

		if debug {
			fmt.Printf("Wrote text record (%d bytes): %s", nt, text)
		}
	}

	// Write end record
	end := fmt.Sprintf("E%s\n", Word(c.pcstartaddr))
	ne, err := file.WriteString(end)
	if err != nil {
		return fmt.Errorf("failed to write end record: %w", err)
	}

	if debug {
		fmt.Printf("Wrote end record (%d bytes): %s", ne, end)
	}

	return nil
}

func (c *Code) PrettyPrint() {
	if !prettyPrint {
		return
	}

	fmt.Println("--- Pretty Print ---")
	for _, node := range c.instructions {
		fmt.Println(node.pretty())
	}
	fmt.Println("--------------------")
}
