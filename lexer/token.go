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
	LBRACE     = "LBRACE"
	RBRACE     = "RBRACE"
	COMMA      = "COMMA"
	EXPR_OPEN  = "EXPR_OPEN"
	EXPR_CLOSE = "EXPR_CLOSE"
	SEMICOLON  = "SEMICOLON"
	COLON      = "COLON"

	RETURN = "RETURN"
	LET    = "LET"
	ASSIGN = "="
	FUNC   = "FUNC"

	EQL = "=="
	NEQ = "!="
	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	TRUE  = "TRUE"
	FALSE = "FALSE"
	NULL  = "NULL"
)
