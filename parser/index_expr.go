package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type IndexExpr struct {
	Target Expression
	Index  Expression
}

func (i *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", i.Target.String(), i.Index.String())
}

func (p *Parser) parseIndexExpr(target Expression) (Expression, error) {
	p.nextToken()

	if p.currentToken.Type == lexer.RBRACKET {
		return nil, fmt.Errorf("index cannot be empty")
	}

	indexExr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.RBRACKET {
		return nil, fmt.Errorf("expected closing ']' after index expression")
	}
	p.nextToken()

	return &IndexExpr{Target: target, Index: indexExr}, nil
}
