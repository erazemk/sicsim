package sim

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type device struct {
	num    byte
	name   string
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewDevice creates a new device
func newDevice(id byte) (*device, error) {
	dev := device{num: id}

	switch id {
	case 0:
		dev.name = "stdin"
		dev.reader = bufio.NewReader(os.Stdin)
	case 1:
		dev.name = "stdout"
		dev.writer = bufio.NewWriter(os.Stdout)
	case 2:
		dev.name = "stderr"
		dev.writer = bufio.NewWriter(os.Stderr)
	default:
		dev.name = fmt.Sprintf("%02X.dev", dev.num)

		infd, err := os.OpenFile(dev.name, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to create input device: %w", err)
		}

		dev.reader = bufio.NewReader(infd)

		outfd, err := os.OpenFile(dev.name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to create output device: %w", err)
		}

		dev.writer = bufio.NewWriter(outfd)
	}

	if debug {
		log.Println("Added new device:", dev.name)
	}

	return &dev, nil
}

// test checks if a device is available for reading or writing
func (d *device) test() bool {
	return d.reader != nil && d.writer != nil
}

// read reads a byte from device
func (d *device) read() (byte, error) {
	if d.reader == nil {
		return 0, fmt.Errorf("failed to open device '%s' for reading", d.name)
	}

	val, err := d.reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			val = 0
		} else {
			return val, fmt.Errorf("failed to read from device '%s': %w", d.name, err)
		}
	}

	if debug {
		log.Printf("Read byte '%c' from device '%s'\n", val, d.name)
	}

	return val, nil
}

// write writes a byte to device
func (d *device) write(val byte) error {
	if d.writer == nil {
		return fmt.Errorf("failed to open device '%s' for writing", d.name)
	}

	err := d.writer.WriteByte(val)
	if err != nil {
		return fmt.Errorf("failed to write to device '%s': %w", d.name, err)
	}

	d.writer.Flush()

	if debug {
		log.Printf("Wrote byte '%c' to device '%s'\n", val, d.name)
	}

	return nil
}
