all: help

help:
	@echo "Usage: make (sicsim | sicasm)"

sicsim:
	go build git.sr.ht/~erazemk/sicsim/cmd/sicsim

sicasm:
	go build git.sr.ht/~erazemk/sicsim/cmd/sicasm
