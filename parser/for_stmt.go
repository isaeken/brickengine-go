package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type ForStatement struct {
	Init      Expression
	Condition Expression
	Update    Expression

	VarName  string
	Iterable Expression

	Body []Expression
}

func (f *ForStatement) String() string {
	return "for (...) { ... }"
}

func (p *Parser) parseForStatement() (Expression, error) {
	p.nextToken()

	if p.currentToken.Type == lexer.LET {
		init, err := p.parseStatement()
		if err != nil {
			return nil, fmt.Errorf("invalid for-loop init: %w", err)
		}
		if p.currentToken.Type != lexer.SEMICOLON {
			return nil, fmt.Errorf("expected ';' after init statement")
		}
		p.nextToken()

		cond, err := p.ParseExpression()
		if err != nil {
			return nil, fmt.Errorf("invalid condition in for-loop: %w", err)
		}
		if p.currentToken.Type != lexer.SEMICOLON {
			return nil, fmt.Errorf("expected ';' after for-loop condition")
		}
		p.nextToken()

		update, err := p.parseStatement()
		if err != nil {
			return nil, fmt.Errorf("invalid update statement in for-loop: %w", err)
		}

		if p.currentToken.Type != lexer.LBRACE {
			return nil, fmt.Errorf("expected '}' after for-loop header")
		}
		p.nextToken()

		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		return &ForStatement{
			Init:      init,
			Condition: cond,
			Update:    update,
			Body:      body,
		}, nil
	}

	if p.currentToken.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected identifier in foreach-style loop")
	}
	varName := p.currentToken.Literal
	p.nextToken()

	if p.currentToken.Type != lexer.IN {
		return nil, fmt.Errorf("expected 'in' after variable name in foreach-style loop")
	}
	p.nextToken()

	iterable, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("invalid iterable expression: %w", err)
	}

	if p.currentToken.Type != lexer.LBRACE {
		return nil, fmt.Errorf("expected '{' to start foreach-loop body")
	}
	p.nextToken()

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ForStatement{
		VarName:  varName,
		Iterable: iterable,
		Body:     body,
	}, nil
}
