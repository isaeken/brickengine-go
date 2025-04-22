package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type PipeExpr struct {
	Left  Expression
	Right Expression
}

func (p *PipeExpr) String() string {
	return fmt.Sprintf("%s | %s", p.Left.String(), p.Right.String())
}

func (p *Parser) parsePipeExpr() (Expression, error) {
	left, err := p.parseBinaryExpr(0)
	if err != nil {
		return nil, fmt.Errorf("invalid left side of pipe expression: %w", err)
	}

	for p.currentToken.Type == lexer.PIPE {
		p.nextToken()

		if p.currentToken.Type == lexer.PIPE || p.currentToken.Type == lexer.EOF {
			return nil, fmt.Errorf("expected expression after pipe ('|'), got '%s'", p.currentToken.Literal)
		}

		right, err := p.parseBinaryExpr(0)
		if err != nil {
			return nil, fmt.Errorf("invalid right side of pipe expression: %w", err)
		}

		left = &PipeExpr{
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}
