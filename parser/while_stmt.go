package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type WhileStatement struct {
	Condition Expression
	Body      []Expression
}

func (w *WhileStatement) String() string {
	return "while (...) { ... }"
}

func (p *Parser) parseWhileStatement() (Expression, error) {
	p.nextToken()

	cond, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("invalid while condition: %w", err)
	}

	if p.currentToken.Type != lexer.LBRACE {
		return nil, fmt.Errorf("expected '{' to start while block")
	}
	p.nextToken()

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &WhileStatement{
		Condition: cond,
		Body:      body,
	}, nil
}
