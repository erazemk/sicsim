package asm

type Mnemonic interface {
	Parse(*Parser) Node
}

type mnemonic struct {
	name        string
	opcode      int
	hint        string
	description string
}

type mnemonicD struct{}
type mnemonicDn struct{}
type mnemonicF1 struct{}
type mnemonicF2n struct{}
type mnemonicF2r struct{}
type mnemonicF2rn struct{}
type mnemonicF2rr struct{}
type mnemonicF3 struct{}
type mnemonicF3m struct{}
type mnemonicF4m struct{}
type mnemonicSd struct{}
type mnemonicSn struct{}

func NewMnemonic(name, hint, desc string, opcode int) Mnemonic {
	return &mnemonic{
		name:        name,
		opcode:      opcode,
		hint:        hint,
		description: desc,
	}
}

func (m *mnemonic) Parse(p *Parser) Node {
	return NewNode(m)
}
