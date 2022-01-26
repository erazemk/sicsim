all: help

help:
	@echo "Usage: make (sicsim | sicasm)"

sicsim:
	go build github.com/erazemk/sicsim/cmd/sicsim

sicasm:
	go build github.com/erazemk/sicsim/cmd/sicasm
