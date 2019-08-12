package lexer

import (
	"github.com/marcel/monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()

	return lex
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	t, ok := token.SingleByteLiteralToType[l.ch]
	if ok {
		defer l.readChar()
		return t.Token(string(l.ch))
	}

	switch l.ch {
	case '=':
		defer l.readChar()
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			return token.EQ.Token(string(ch) + string(l.ch))
		}
		return token.ASSIGN.Token(string(l.ch))
	case '!':
		defer l.readChar()
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			return token.NOT_EQ.Token(string(ch) + string(l.ch))
		}
		return token.BANG.Token(string(l.ch))
	}

	switch {
	case isLetter(l.ch):
		literal := l.readIdentifier()
		return token.LookupIdent(literal).Token(literal)
	case isDigit(l.ch):
		return token.INT.Token(l.readNumber())
	}

	defer l.readChar()
	return token.ILLEGAL.Token(string(l.ch))
}

func (l *Lexer) readWhile(predicate func(byte) bool) string {
	position := l.position

	for predicate(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	return l.readWhile(isDigit)
}

func (l *Lexer) readIdentifier() string {
	return l.readWhile(isLetter)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
