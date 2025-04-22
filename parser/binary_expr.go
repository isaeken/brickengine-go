package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator, b.Right.String())
}

func (p *Parser) parseBinaryExpr(precedence int) (Expression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, fmt.Errorf("invalid left-hand side of binary expression: %w", err)
	}

	for isOperator(p.currentToken.Type) {
		op := p.currentToken.Literal
		p.nextToken()

		if p.currentToken.Type == lexer.EOF {
			return nil, fmt.Errorf("expected right-hand side after operator '%s'", op)
		}

		right, err := p.parsePrimary()
		if err != nil {
			return nil, fmt.Errorf("invalid right-hand side of binary expression: %w", err)
		}

		left = &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}
