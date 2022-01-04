package asm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var isPCset bool

// ParseFile reads the contents of the provided file and sends each line to ParseLine
func (c *Code) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	for sc := bufio.NewScanner(file); sc.Scan(); {
		if err := c.ParseLine(sc.Text()); err != nil {
			return fmt.Errorf("failed to parse line: %w", err)
		}
	}

	c.length = c.lc
	return nil
}

// ParseLine parses a provided line and sets the code's attributes accordingly
func (c *Code) ParseLine(line string) error {
	// Remove comments (only parse line until comment)
	if strings.ContainsRune(line, '.') {
		line = line[:strings.IndexRune(line, '.')]
	}

	// Empty line (after removing comments)
	if len(line) == 0 {
		return nil
	}

	// Split command into parts
	command := strings.Fields(line)
	if len(command) == 0 { // Only spaces in line
		return nil
	}

	node := NewNode(command, c.lc, c.brelative)

	// Check if label already exists in symtab
	if node.label != "" {
		if _, exists := c.symtab[node.label]; !exists {
			c.symtab[node.label] = c.lc

			if debug {
				fmt.Printf("Added '%s' to symtab at %d\n", node.label, c.symtab[node.label])
			}
		} else {
			return fmt.Errorf("label '%s' already declared", node.label)
		}
	}

	// Set program name and start address
	if node.mnemonic == "START" {
		c.startaddr = node.operand
		c.name = node.label

		if len(c.name) > 6 {
			return fmt.Errorf("program name must not be longer than 6 characters")
		}

		if debug {
			fmt.Printf("Set start address to '%[1]d (%06[1]X)'\n", c.startaddr)
		}
	}

	// Set PC start address based on where the first instruction is
	if !isPCset && inSlice(node.mnemonic, Instructions) {
		c.pcstartaddr = c.lc
		isPCset = true

		if debug {
			fmt.Printf("Set PC start address to '%[1]d (%06[1]X)' at instruction '%[2]s'\n", c.pcstartaddr, node.mnemonic)
		}
	}

	// Set base relative attributes
	switch node.mnemonic {
	case "ORG":
		c.lc = node.operand
	case "BASE":
		c.brelative = true
	case "NOBASE":
		c.brelative = false
	}

	c.instructions = append(c.instructions, node)
	c.lc += node.length
	return nil
}
