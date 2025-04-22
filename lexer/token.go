package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"

	IDENT      = "IDENT"
	NUMBER     = "NUMBER"
	STRING     = "STRING"
	OPERATOR   = "OPERATOR"
	PIPE       = "PIPE"
	DOT        = "DOT"
	LPAREN     = "LPAREN"
	RPAREN     = "RPAREN"
	LBRACKET   = "LBRACKET"
	RBRACKET   = "RBRACKET"
	COMMA      = "COMMA"
	EXPR_OPEN  = "EXPR_OPEN"
	EXPR_CLOSE = "EXPR_CLOSE"
)
