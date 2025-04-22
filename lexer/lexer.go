package lexer

import "unicode"

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // null byte
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case 0:
		return Token{Type: EOF, Literal: ""}
	case '{':
		if l.peekChar() == '{' {
			l.readChar()
			l.readChar()
			return Token{Type: EXPR_OPEN, Literal: "{{"}
		}
	case '}':
		if l.peekChar() == '}' {
			l.readChar()
			l.readChar()
			return Token{Type: EXPR_CLOSE, Literal: "}}"}
		}
	case '+', '-', '*', '/':
		ch := l.ch
		l.readChar()
		return Token{Type: OPERATOR, Literal: string(ch)}
	case '|':
		l.readChar()
		return Token{Type: PIPE, Literal: "|"}
	case '.':
		l.readChar()
		return Token{Type: DOT, Literal: "."}
	case '(':
		l.readChar()
		return Token{Type: LPAREN, Literal: "("}
	case ')':
		l.readChar()
		return Token{Type: RPAREN, Literal: ")"}
	case '[':
		l.readChar()
		return Token{Type: LBRACKET, Literal: "["}
	case ']':
		l.readChar()
		return Token{Type: RBRACKET, Literal: "]"}
	case ',':
		l.readChar()
		return Token{Type: COMMA, Literal: ","}
	case '"':
		return l.readString()
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			return Token{Type: IDENT, Literal: ident}
		} else if isDigit(l.ch) {
			number := l.readNumber()
			return Token{Type: NUMBER, Literal: number}
		} else {
			ch := l.ch
			l.readChar()
			return Token{Type: ILLEGAL, Literal: string(ch)}
		}
	}

	return Token{Type: ILLEGAL, Literal: ""}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() Token {
	l.readChar()
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[start:l.position]
	l.readChar()
	return Token{Type: STRING, Literal: str}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
