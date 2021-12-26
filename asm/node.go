package asm

type Node interface {
	Enter(code)
	Leave(code)
}

type node struct {
	label    string
	comm     string
	mnemonic Mnemonic
}

type Comment interface {
	Enter(code)
	Leave(code)
}

type comment node

func NewComment(comm string) Comment {
	return &comment{
		comm: comm,
	}
}

func NewNode(mnemonic Mnemonic) Node {
	return &node{
		mnemonic: mnemonic,
	}
}

func (n *node) Enter(code code)    {}
func (n *node) Leave(code code)    {}
func (c *comment) Enter(code code) {}
func (c *comment) Leave(code code) {}

func (n *node) Label() string {
	return n.label
}

func (n *node) SetLabel(label string) {
	n.label = label
}

func (n *node) Comment() string {
	return n.comm
}

func (n *node) SetComment(comment string) {
	n.comm = comment
}

func (c *comment) Comment() string {
	return c.comm
}
