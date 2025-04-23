package lexer

import (
	"strings"
	"unicode"
)

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
	l.skipComment()

	switch l.ch {
	case 0:
		return Token{Type: EOF, Literal: ""}
	case '{':
		if l.peekChar() == '{' {
			l.readChar()
			l.readChar()
			return Token{Type: EXPR_OPEN, Literal: "{{"}
		}

		l.readChar()
		return Token{Type: LBRACE, Literal: "{"}
	case '}':
		if l.peekChar() == '}' {
			l.readChar()
			l.readChar()
			return Token{Type: EXPR_CLOSE, Literal: "}}"}
		}

		l.readChar()
		return Token{Type: RBRACE, Literal: "}"}
	case '+', '-', '*', '/':
		if l.ch == '-' && isDigit(l.peekChar()) {
			number := l.readNumber()
			return Token{Type: NUMBER, Literal: number}
		}

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
		return l.readString('"')
	case '\'':
		return l.readString('\'')
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return Token{Type: EQL, Literal: "=="}
		}
		l.readChar()
		return Token{Type: ASSIGN, Literal: "="}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return Token{Type: NEQ, Literal: "!="}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return Token{Type: LTE, Literal: "<="}
		}
		l.readChar()
		return Token{Type: LT, Literal: "<"}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return Token{Type: GTE, Literal: ">="}
		}
		l.readChar()
		return Token{Type: GT, Literal: ">"}
	case ';':
		l.readChar()
		return Token{Type: SEMICOLON, Literal: ";"}
	case ':':
		l.readChar()
		return Token{Type: COLON, Literal: ":"}
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()

			switch strings.ToLower(ident) {
			case "return":
				return Token{Type: RETURN, Literal: ident}
			case "true":
				return Token{Type: TRUE, Literal: ident}
			case "false":
				return Token{Type: FALSE, Literal: ident}
			case "null":
				return Token{Type: NULL, Literal: ident}
			case "let":
				return Token{Type: LET, Literal: ident}
			case "fn":
				return Token{Type: FUNC, Literal: ident}
			case "for":
				return Token{Type: FOR, Literal: ident}
			case "in":
				return Token{Type: IN, Literal: ident}
			case "while":
				return Token{Type: WHILE, Literal: ident}
			case "try":
				return Token{Type: TRY, Literal: ident}
			case "catch":
				return Token{Type: CATCH, Literal: ident}
			default:
				return Token{Type: IDENT, Literal: ident}
			}
		} else if isDigit(l.ch) {
			number := l.readNumber()
			return Token{Type: NUMBER, Literal: number}
		} else {
			ch := l.ch
			l.readChar()
			return Token{Type: ILLEGAL, Literal: string(ch)}
		}
	}

	return Token{Type: ILLEGAL, Literal: string(l.ch)}
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
	for isDigit(l.ch) || l.ch == '.' || l.ch == 'e' || l.ch == '+' || l.ch == '-' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString(quote byte) Token {
	l.readChar()
	var str []rune

	for l.ch != 0 && l.ch != quote {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				str = append(str, '\n')
			case 'r':
				str = append(str, '\r')
			case 't':
				str = append(str, '\t')
			case '\\':
				str = append(str, '\\')
			case quote:
				str = append(str, rune(quote))
			default:
				str = append(str, '\\', rune(l.ch))
			}
		} else {
			str = append(str, rune(l.ch))
		}

		l.readChar()
	}

	l.readChar()
	return Token{Type: STRING, Literal: string(str)}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for {
		if l.ch == '#' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		}

		if l.ch == '/' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		}

		if unicode.IsSpace(rune(l.ch)) {
			l.skipWhitespace()
			continue
		}

		if l.ch != '#' && !(l.ch == '/' && l.peekChar() == '/') {
			break
		}
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
