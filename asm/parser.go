package asm

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	lexer     Lexer
	mnemonics map[string]Mnemonic
}

func NewParser() Parser {
	p := Parser{}
	p.InitMnemonics()
	return p
}

func (p *Parser) ParseLabel() string {
	if p.lexer.col == 1 && unicode.IsLetter(p.lexer.Peek(0)) {
		return p.lexer.ReadAlphanumeric()
	}

	return ""
}

func (p *Parser) ParseMnemonic() Mnemonic {
	isExtended := p.lexer.AdvanceIf('+')
	name := p.lexer.ReadAlphanumeric()

	var mnem Mnemonic

	if isExtended {
		mnem = p.Mnemonic("+" + name)
	} else {
		mnem = p.Mnemonic(name)
	}

	if mnem == nil {
		// TODO
		fmt.Println("Error")
	}

	return mnem
}

func (p *Parser) ParseSymbol() string {
	return p.lexer.ReadAlphanumeric()
}

func (p *Parser) ParseRegister() int {
	char := p.lexer.Advance()
	reg := strings.Index("AXLBSTF", string(char))

	if reg < 0 {
		// TODO
		fmt.Println("Error")
	}

	return reg
}

func (p *Parser) ParseComma() {
	l := p.lexer
	l.SkipWhitespace()
	l.AdvanceChar(',')
	l.SkipWhitespace()
}

func (p *Parser) ParseIndexed() bool {
	l := p.lexer
	l.SkipWhitespace()

	if l.AdvanceIf(',') {
		l.SkipWhitespace()
		l.AdvanceChar('X')
		return true
	}

	return false
}

func (p *Parser) ParseNumber(low, high int) int {
	var num int
	l := p.lexer

	if l.Peek(0) == '0' {
		r := -1

		switch l.Peek(1) {
		case 'b':
			r = 2
		case 'o':
			r = 8
		case 'x':
			r = 16
		}

		if r != -1 {
			l.Advance()
			l.Advance()

			t, err := strconv.ParseInt(l.ReadDigits(r), r, 32)
			if err != nil {
				panic(err)
			}

			num = int(t)
		} else {
			t, err := strconv.ParseInt(l.ReadDigits(10), 10, 32)
			if err != nil {
				panic(err)
			}

			num = int(t)
		}
	} else if unicode.IsDigit(l.Peek(0)) {
		t, err := strconv.ParseInt(l.ReadDigits(10), 10, 32)
		if err != nil {
			panic(err)
		}

		num = int(t)
	} else {
		// TODO
		fmt.Println("Error: Number expected")
	}

	if unicode.IsLetter(l.Peek(0)) || unicode.IsDigit(l.Peek(0)) {
		// TODO
		fmt.Println("Error: Number must not be followed by letter or digit")
	}

	if num < low || num > high {
		// TODO
		fmt.Println("Error: Number our of range")
	}

	return num
}

func (p *Parser) ParseData() []byte {
	l := p.lexer

	if l.AdvanceIf('C') {
		l.AdvanceChar('\'')
		return []byte(l.ReadTo('\''))
	} else if l.AdvanceIf('X') {
		l.AdvanceChar('\'')
		s := l.ReadTo('\'')
		data := make([]byte, len(s)/2)

		for i := 0; i < len(data); i++ {
			val, err := strconv.ParseInt(s[2*i:2*i+2], 16, 8)
			if err != nil {
				panic(err)
			}

			data[i] = byte(val)
		}

		return data
	} else if unicode.IsDigit(l.Peek(0)) {
		num := p.ParseNumber(0, int(math.Pow(2, 24)-1))
		data := make([]byte, 3)
		data[2] = byte(num)
		data[1] = byte(num >> 8)
		data[0] = byte(num >> 16)

		return data
	}

	// TODO: Throw error
	fmt.Println("Error: Invalid storage specifier")
	return nil
}

func (p *Parser) ParseInstruction() Node {
	l := p.lexer

	if l.col == 1 && l.Peek(0) == '.' {
		return NewComment(l.ReadTo('\n'))
	}

	label := p.ParseLabel()

	if l.SkipWhitespace() && label == "" {
		l.Advance()
		return nil
	}

	mnemonic := p.ParseMnemonic()
	l.SkipWhitespace()
	node := mnemonic.Parse(p)
	// TODO
	// node.SetLabel(label)
	// node.SetComment(l.ReadTo('\n'))
	return node
}

func (p *Parser) ParseCode() Code {
	l := p.lexer
	code := NewCode()

	for l.Peek(0) > 0 {
		for l.Peek(0) > 0 && l.col > 1 {
			l.ReadTo('\n')
		}

		instruction := p.ParseInstruction()

		if instruction != nil {
			code.Append(instruction)
		}
	}

	return code
}

func (p *Parser) Parse(input string) Code {
	p.lexer = NewLexer(input)
	return p.ParseCode()
}

func (p *Parser) Mnemonic(name string) Mnemonic {
	if p.mnemonics[name] != nil {
		return p.mnemonics[name]
	}

	return nil
}

func (p *Parser) PutMnemonic(mnemonic mnemonic) {
	p.mnemonics[mnemonic.name] = &mnemonic
}

func (p *Parser) InitMnemonics() {
	p.mnemonics = make(map[string]Mnemonic)

	// TODO: Add all mnemonics
	p.PutMnemonic(mnemonic{name: "NOBASE", opcode: 1, hint: "hint", description: "description"})
}
