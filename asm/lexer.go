package asm

import (
	"fmt"
	"strconv"
	"unicode"
)

type Lexer struct {
	input string
	pos   int
	row   int
	col   int
	start int
}

func NewLexer(input string) Lexer {
	return Lexer{
		input: input,
		pos:   0,
		row:   1,
		col:   1,
		start: 0,
	}
}

func (l *Lexer) Mark() {
	l.start = l.pos
}

func (l *Lexer) Extract(offset int) string {
	return l.input[l.start : l.pos+offset]
}

func (l *Lexer) Peek(ahead int) rune {
	if l.pos+ahead < len(l.input) {
		return rune(l.input[l.pos+ahead])
	}

	return 0
}

func (l *Lexer) Advance() rune {
	char := l.Peek(0)
	l.pos++

	if char == '\n' {
		l.row++
		l.col = 1
	} else if char == '\t' {
		l.col = ((l.col-1)/4)*4 + 5
	} else {
		l.col++
	}

	return char
}

func (l *Lexer) AdvanceIf(char rune) bool {
	if l.Peek(0) != char {
		return false
	}

	l.Advance()
	return true
}

func (l *Lexer) AdvanceChar(char rune) {
	if !l.AdvanceIf(char) {
		// TODO
		fmt.Println("Error")
	}
}

func (l *Lexer) SkipWhitespace() bool {
	for l.Peek(0) == ' ' || l.Peek(0) == '\t' {
		l.Advance()
	}

	return l.Peek(0) != '\n' || l.Peek(0) == 0
}

func (l *Lexer) ReadTo(delim rune) string {
	l.Mark()

	for l.Peek(0) > 0 && l.Peek(0) != delim {
		l.Advance()
	}

	l.Advance()
	return l.Extract(-1)
}

func (l *Lexer) ReadAlphanumeric() string {
	l.Mark()

	for unicode.IsDigit(l.Peek(0)) || unicode.IsLetter(l.Peek(0)) || l.Peek(0) == '_' {
		l.Advance()
	}

	return l.Extract(0)
}

func (l *Lexer) ReadDigits(radix int) string {
	l.Mark()

	for val, _ := strconv.ParseInt(string(l.Peek(0)), radix, 32); val != -1; {
		l.Advance()
	}

	return l.Extract(0)
}
