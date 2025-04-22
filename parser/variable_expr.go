package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
	"strings"
)

type VariableExpr struct {
	Parts []string
}

func (v *VariableExpr) String() string {
	return strings.Join(v.Parts, ".")
}

func (p *Parser) parseVariableExpr() (Expression, error) {
	if p.currentToken.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected identifier at start of variable expression, got '%s'", p.currentToken.Literal)
	}

	parts := []string{p.currentToken.Literal}
	p.nextToken()

	// dot notation
	for p.currentToken.Type == lexer.DOT {
		p.nextToken()
		if p.currentToken.Type != lexer.IDENT {
			return nil, fmt.Errorf("expected identifier after '.'")
		}
		parts = append(parts, p.currentToken.Literal)
		p.nextToken()
	}

	var expr Expression = &VariableExpr{Parts: parts}

	// index access
	for p.currentToken.Type == lexer.LBRACKET {
		if p.peekToken.Type == lexer.RBRACKET {
			return nil, fmt.Errorf("empty index expression is not allowed")
		}

		idxExpr, err := p.parseIndexExpr(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid index on variable expression: %w", err)
		}
		expr = idxExpr
	}

	return expr, nil
}
