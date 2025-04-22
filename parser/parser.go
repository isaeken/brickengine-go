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

	return expr, nil
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
	case lexer.TRUE:
		p.nextToken()
		return &BoolLiteral{Value: true}, nil
	case lexer.FALSE:
		p.nextToken()
		return &BoolLiteral{Value: false}, nil
	case lexer.NULL:
		p.nextToken()
		return &NullLiteral{}, nil
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
	case lexer.LBRACE:
		return p.parseObjectExpr()
	case lexer.LBRACKET:
		return p.parseArrayLiteral()
	default:
		return nil, fmt.Errorf("unexpected token %s", p.currentToken.Literal)
	}
}

func (p *Parser) parseArrayLiteral() (Expression, error) {
	var elements []Expression
	p.nextToken()

	for p.currentToken.Type != lexer.RBRACKET && p.currentToken.Type != lexer.EOF {
		elem, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)

		if p.currentToken.Type == lexer.COMMA {
			p.nextToken()
		} else if p.currentToken.Type != lexer.RBRACKET {
			return nil, fmt.Errorf("expected ',' or ']' in array, got '%s'", p.currentToken.Literal)
		}
	}

	if p.currentToken.Type != lexer.RBRACKET {
		return nil, fmt.Errorf("expected closing ']' for array, got '%s'", p.currentToken.Literal)
	}
	p.nextToken()

	return &ArrayLiteral{Elements: elements}, nil
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

func (p *Parser) parseStatement() (Expression, error) {
	switch p.currentToken.Type {
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.IDENT:
		if p.currentToken.Literal == "if" {
			return p.parseIfStatement()
		}
		fallthrough
	case lexer.FUNC:
		if p.currentToken.Literal == "fn" {
			return p.parseFnStatement()
		}
		fallthrough
	default:
		return p.tryAssignmentOrExpression()
	}
}

func (p *Parser) Parse() ([]Expression, error) {
	var stmts []Expression

	for p.currentToken.Type != lexer.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)

		if p.currentToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
	}

	return stmts, nil
}

func isOperator(t lexer.TokenType) bool {
	return t == lexer.OPERATOR ||
		t == lexer.EQL || t == lexer.NEQ ||
		t == lexer.LT || t == lexer.GT ||
		t == lexer.LTE || t == lexer.GTE
}
