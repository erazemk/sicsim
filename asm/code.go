package asm

type Code interface {
	Append(Node)
}

type code struct {
	name         string
	instructions []Node
	lc           int
}

func NewCode() Code {
	return &code{}
}

func (c *code) Append(instruction Node) {
	c.instructions = append(c.instructions, instruction)
}
