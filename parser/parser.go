package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
	"strconv"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseExpression() (Expression, error) {
	expr, err := p.parsePipeExpr()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.EOF && p.currentToken.Type != lexer.EXPR_CLOSE {
		return nil, fmt.Errorf("unexpected token '%s' at end of expression", p.currentToken.Literal)
	}

	return expr, nil
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

func (p *Parser) parsePrimary() (Expression, error) {
	switch p.currentToken.Type {
	case lexer.IDENT:
		if p.peekToken.Type == lexer.LPAREN {
			return p.parseCallExpr()
		}

		return p.parseVariableExpr()
	case lexer.STRING:
		str := &StringLiteral{Value: p.currentToken.Literal}
		p.nextToken()
		return str, nil
	case lexer.NUMBER:
		num, err := strconv.ParseFloat(p.currentToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s'", p.currentToken.Literal)
		}
		p.nextToken()
		return &NumberLiteral{Value: num}, nil
	case lexer.LPAREN:
		p.nextToken()
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, fmt.Errorf("invalid expression in parenthesis: %w", err)
		}
		if p.currentToken.Type != lexer.RPAREN {
			return nil, fmt.Errorf("expected ')' after expression")
		}
		p.nextToken()
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token %s", p.currentToken.Literal)
	}
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

func (p *Parser) parseArguments() ([]Expression, error) {
	var args []Expression

	if p.currentToken.Type == lexer.RPAREN {
		p.nextToken()
		return args, nil
	}

	for {
		arg, err := p.parsePrimary()
		if err != nil {
			return nil, fmt.Errorf("invalid function argument: %w", err)
		}
		args = append(args, arg)

		if p.currentToken.Type == lexer.COMMA {
			p.nextToken()
			continue
		}

		if p.currentToken.Type == lexer.RPAREN {
			p.nextToken()
			break
		}

		return nil, fmt.Errorf("expected ',' or ')' in argument list, got '%s'", p.currentToken.Literal)
	}

	return args, nil
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

func (p *Parser) parseIndexExpr(target Expression) (Expression, error) {
	p.nextToken()

	if p.currentToken.Type == lexer.RBRACKET {
		return nil, fmt.Errorf("index cannot be empty")
	}

	indexExr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.RBRACKET {
		return nil, fmt.Errorf("expected closing ']' after index expression")
	}
	p.nextToken()

	return &IndexExpr{Target: target, Index: indexExr}, nil
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

func isOperator(t lexer.TokenType) bool {
	return t == lexer.OPERATOR
}
