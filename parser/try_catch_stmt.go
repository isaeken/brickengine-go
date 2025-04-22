package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type TryCatchStatement struct {
	TryBlock   []Expression
	CatchBlock []Expression
}

func (t *TryCatchStatement) String() string {
	return "try { ... } catch { ... }"
}

func (p *Parser) parseTryCatchStatement() (Expression, error) {
	p.nextToken()

	if p.currentToken.Type != lexer.LBRACE {
		return nil, fmt.Errorf("expected '{' but got '%s'", p.currentToken.Literal)
	}
	p.nextToken()
	tryBlock, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.CATCH || p.currentToken.Literal != "catch" {
		return nil, fmt.Errorf("expected 'catch' but got '%s'", p.currentToken.Literal)
	}
	p.nextToken()

	if p.currentToken.Type != lexer.LBRACE {
		return nil, fmt.Errorf("expected '}' but got %s", p.currentToken.Literal)
	}
	p.nextToken()
	catchBlock, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &TryCatchStatement{
		TryBlock:   tryBlock,
		CatchBlock: catchBlock,
	}, nil
}
