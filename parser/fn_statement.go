package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type FnStatement struct {
	Name string
	Args []string
	Body []Expression
}

func (s *FnStatement) String() string {
	return fmt.Sprintf("declare[%s(%v)]", s.Name, s.Args)
}

func (p *Parser) parseFnStatement() (Expression, error) {
	p.nextToken()

	if p.currentToken.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected function name")
	}
	name := p.currentToken.Literal

	p.nextToken()
	if p.currentToken.Type != lexer.LPAREN {
		return nil, fmt.Errorf("expected '(' after function name")
	}

	args := []string{}
	p.nextToken()
	for p.currentToken.Type != lexer.RPAREN {
		if p.currentToken.Type != lexer.IDENT {
			return nil, fmt.Errorf("expected identifier in argument list")
		}
		args = append(args, p.currentToken.Literal)
		p.nextToken()
		if p.currentToken.Type == lexer.COMMA {
			p.nextToken()
		}
	}
	p.nextToken()

	if p.currentToken.Type != lexer.LBRACE {
		return nil, fmt.Errorf("expected '{' to start function body")
	}
	p.nextToken()

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &FnStatement{
		Name: name,
		Args: args,
		Body: body,
	}, nil
}
