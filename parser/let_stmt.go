package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type LetStatement struct {
	Name  string
	Value Expression
}

func (s *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s", s.Name, s.Value.String())
}

func (p *Parser) parseLetStatement() (Expression, error) {
	p.nextToken()
	if p.currentToken.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected identifier after 'let', got %s", p.currentToken.Type)
	}
	name := p.currentToken.Literal

	p.nextToken() // '='
	if p.currentToken.Type != lexer.ASSIGN {
		return nil, fmt.Errorf("expected '=' after identifier in let statement, got %s", p.currentToken.Type)
	}

	p.nextToken()
	value, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &LetStatement{name, value}, nil
}
