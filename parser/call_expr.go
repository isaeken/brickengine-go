package parser

import (
	"fmt"
	"strings"
)

type CallExpr struct {
	Target Expression
	Args   []Expression
}

func (c *CallExpr) String() string {
	args := []string{}
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", c.Target.String(), strings.Join(args, ","))
}

func (p *Parser) parseCallExpr(target Expression) (Expression, error) {
	p.nextToken()
	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}
	return &CallExpr{
		Target: target,
		Args:   args,
	}, nil
}
