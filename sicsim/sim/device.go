package sim

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Device struct {
	devno   byte
	devfile string
}

// New creates a new device with and identifier and device type
func (d *Device) New(devno byte) {
	if isDevice(devno) {
		switch devno {
		case 0:
			d.devfile = "stdin" // TODO: Handle writing to stdin
		case 1:
			d.devfile = "stdout" // TODO: Handle writing to stdout
		case 2:
			d.devfile = "stderr" // TODO: Handle writing to stderr
		default:
			d.devno = devno
			d.devfile = fmt.Sprintf("0x%x.dev", d.devno)
		}

		if debug {
			log.Println("Added a new device with id", devno)
		}
	}
}

// Read reads a byte from device
func (d *Device) Read() byte {
	file, err := os.Open(d.devfile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	val, err := reader.ReadByte()

	if err != nil {
		log.Fatal(err)
	}

	return val
}

// Write writes a byte to device
func (d *Device) Write(val byte) {
	file, err := os.Open(d.devfile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	err = writer.WriteByte(val)

	if err != nil {
		log.Fatal(err)
	}
}
