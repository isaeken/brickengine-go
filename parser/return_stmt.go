package parser

type ReturnStatement struct {
	Value Expression
}

func (r *ReturnStatement) String() string {
	return "return " + r.Value.String()
}

func (p *Parser) parseReturnStatement() (*ReturnStatement, error) {
	p.nextToken()
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	return &ReturnStatement{Value: expr}, nil
}
