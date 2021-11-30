package sim

import (
	"fmt"
	"log"
	"os"
)

type device struct {
	devno   byte
	devname string
	devfile *os.File
}

// New creates a new device
func (d *device) New(devno byte) {
	d.devno = devno

	switch devno {
	case 0:
		d.devname = "stdin"
		d.devfile = os.Stdin
	case 1:
		d.devname = "stdout"
		d.devfile = os.Stdout
	case 2:
		d.devname = "stderr"
		d.devfile = os.Stderr
	default:
		d.devname = fmt.Sprintf("%x.dev", d.devno)
		d.devfile = nil
	}

	if debug {
		log.Println("Created a new device:", d.devname)
	}
}

// Dev returns a device with the id devno
func (m *Machine) Dev(devno byte) *device {
	dev := m.devs[devno]

	if dev == nil {
		dev.New(devno)
		m.devs[devno] = dev
	}

	return dev
}

// Read reads a byte from device
func (d *device) Read() (byte, error) {
	var file *os.File

	if d.devfile == nil {
		file, err := os.Open(d.devname)
		defer file.Close()

		if err != nil {
			return 0, fmt.Errorf("Failed to open device for reading: %w", err)
		}
	}

	val := make([]byte, 1)

	_, err := file.Read(val)

	if err != nil {
		return 0, fmt.Errorf("Failed to read from device: %w", err)
	}

	if debug {
		log.Printf("Read byte from device %s: %b [%c]\n", d.devname, val[0], val[0])
	}

	return val[0], nil
}

// Write writes a byte to device
func (d *device) Write(val byte) error {
	var file *os.File

	if d.devfile == nil {
		file, err := os.Open(fmt.Sprintf("%x.dev", d.devno))
		defer file.Close()

		if err != nil {
			return fmt.Errorf("Failed to open device for writing: %w", err)
		}
	}

	_, err := file.Write([]byte{val})

	if err != nil {
		return fmt.Errorf("Failed to write to device: %w", err)
	}

	if debug {
		log.Printf("Wrote byte to device %s: %b [%c]\n", d.devname, val, val)
	}

	return nil
}
