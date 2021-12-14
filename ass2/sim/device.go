package sicsim

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type device struct {
	num   byte
	name  string
	infd  *os.File
	outfd *os.File
}

// NewDevice creates a new device
func newDevice(id byte) (*device, error) {
	var err error
	dev := device{num: id}

	switch id {
	case 0:
		dev.name = "stdin"
		dev.infd = os.Stdin
	case 1:
		dev.name = "stdout"
		dev.outfd = os.Stdout
	case 2:
		dev.name = "stderr"
		dev.outfd = os.Stderr
	default:
		dev.name = fmt.Sprintf("%02X.dev", dev.num)
		dev.infd, err = os.OpenFile(dev.name, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)

		if err != nil {
			return nil, fmt.Errorf("failed to create input device: %w", err)
		}

		dev.outfd, err = os.OpenFile(dev.name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			return nil, fmt.Errorf("failed to create output device: %w", err)
		}
	}

	if debug {
		log.Println("Added new device:", dev.name)
	}

	return &dev, nil
}

// test checks if a device is available for reading or writing
func (d *device) test() bool {
	return d.infd != nil || d.outfd != nil
}

// read reads a byte from device
func (d *device) read() (byte, error) {
	if d.infd == nil {
		return 0, fmt.Errorf("failed to open device '%s' for reading", d.name)
	}

	r := bufio.NewReader(d.infd)
	val, err := r.ReadByte()

	if err != nil {
		return val, fmt.Errorf("failed to read from device '%s': %w", d.name, err)
	}

	if debug {
		log.Printf("Read byte '%c' from device '%s'\n", val, d.name)
	}

	return val, nil
}

// write writes a byte to device
func (d *device) write(val byte) error {
	if d.outfd == nil {
		return fmt.Errorf("failed to open device '%s' for writing", d.name)
	}

	w := bufio.NewWriter(d.outfd)
	err := w.WriteByte(val)
	w.Flush()

	if err != nil {
		return fmt.Errorf("failed to write to device '%s': %w", d.name, err)
	}

	if debug {
		log.Printf("Wrote byte '%c' to device '%s'\n", val, d.name)
	}

	return nil
}
