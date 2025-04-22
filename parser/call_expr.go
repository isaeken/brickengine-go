package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
	"strings"
)

type CallExpr struct {
	Name string
	Args []Expression
}

func (c *CallExpr) String() string {
	args := []string{}
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(args, ","))
}

func (p *Parser) parseCallExpr() (Expression, error) {
	funcName := p.currentToken.Literal
	p.nextToken()

	if p.currentToken.Type != lexer.LPAREN {
		return nil, fmt.Errorf("expected '(' after function name")
	}
	p.nextToken()

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	return &CallExpr{Name: funcName, Args: args}, nil
}
