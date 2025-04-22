package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type AssignmentStmt struct {
	Target Expression
	Value  Expression
}

func (s *AssignmentStmt) String() string {
	return fmt.Sprintf("%s = %s", s.Target.String(), s.Value.String())
}

func (p *Parser) parseAssignmentStmt() (Expression, error) {
	target, err := p.parseVariableExpr()
	if err != nil {
		return nil, fmt.Errorf("invalid assignment target: %w", err)
	}

	if p.currentToken.Type != lexer.ASSIGN {
		return nil, fmt.Errorf("expected '=' in assignment statement, got '%s'", p.currentToken.Literal)
	}
	p.nextToken()

	value, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("invalid assignment value: %w", err)
	}

	return &AssignmentStmt{
		Target: target,
		Value:  value,
	}, nil
}

func (p *Parser) tryAssignmentOrExpression() (Expression, error) {
	savedToken := p.currentToken

	target, err := p.parseVariableExpr()
	if err == nil && p.currentToken.Type == lexer.ASSIGN {
		p.nextToken()
		value, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		return &AssignmentStmt{
			Target: target,
			Value:  value,
		}, nil
	}

	p.currentToken = savedToken
	p.peekToken = p.lexer.NextToken()
	return p.ParseExpression()
}
