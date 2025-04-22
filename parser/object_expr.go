package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type ObjectExpr struct {
	Pairs map[string]Expression
}

func (o *ObjectExpr) String() string {
	out := "{ "
	first := true
	for k, v := range o.Pairs {
		if !first {
			out += ", "
		}
		first = false
		out += fmt.Sprintf("%s: %s", k, v.String())
	}
	out += " }"
	return out
}

func (p *Parser) parseObjectExpr() (Expression, error) {
	pairs := make(map[string]Expression)
	p.nextToken()

	for p.currentToken.Type != lexer.RBRACE && p.currentToken.Type != lexer.EOF {
		if p.currentToken.Type != lexer.IDENT {
			return nil, fmt.Errorf("expected identifier in object key, got '%s'", p.currentToken.Literal)
		}
		key := p.currentToken.Literal
		p.nextToken()

		if p.currentToken.Type != lexer.COLON {
			return nil, fmt.Errorf("expected ':' after key '%s', got '%s'", key, p.currentToken.Literal)
		}
		p.nextToken()

		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		pairs[key] = expr

		if p.currentToken.Type == lexer.COMMA {
			p.nextToken()
		} else if p.currentToken.Type != lexer.RBRACE {
			return nil, fmt.Errorf("expected ',' or '}' in object, got '%s'", p.currentToken.Literal)
		}
	}

	if p.currentToken.Type != lexer.RBRACE {
		return nil, fmt.Errorf("expected closing '}' in object, got '%s'", p.currentToken.Literal)
	}
	p.nextToken()

	return &ObjectExpr{Pairs: pairs}, nil
}
