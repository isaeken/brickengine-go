package parser

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
)

type IfStatement struct {
	Condition   Expression
	ThenBlock   []Expression
	ElseIfParts []ElseIfClause
	ElseBlock   []Expression
}

type ElseIfClause struct {
	Condition Expression
	Block     []Expression
}

func (s *IfStatement) String() string {
	out := fmt.Sprintf("if %s { ... }", s.Condition.String())
	for _, e := range s.ElseIfParts {
		out += fmt.Sprintf(" else if %s { ... }", e.Condition.String())
	}
	if s.ElseBlock != nil {
		out += fmt.Sprintf(" else { ... }")
	}
	return out
}

func (p *Parser) parseIfStatement() (Expression, error) {
	p.nextToken()
	condition, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("invalid condition in if statement: %w", err)
	}

	if p.currentToken.Type != lexer.LBRACE {
		return nil, fmt.Errorf("expected '{' after if condition, got '%s'", p.currentToken.Literal)
	}
	p.nextToken()

	thenBlock, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var elseIfParts []ElseIfClause
	var elseBlock []Expression

	for p.peekToken.Type == lexer.IDENT && p.currentToken.Literal == "else" {
		p.nextToken()

		if p.currentToken.Type == lexer.IDENT && p.currentToken.Literal == "if" {
			p.nextToken()

			cond, err := p.ParseExpression()
			if err != nil {
				return nil, fmt.Errorf("invalid condition in else if: %w", err)
			}
			if p.currentToken.Type != lexer.LBRACE {
				return nil, fmt.Errorf("expected '{' after else if condition, got '%s'", p.currentToken.Literal)
			}
			p.nextToken()

			block, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			elseIfParts = append(elseIfParts, ElseIfClause{Condition: cond, Block: block})
		} else if p.currentToken.Type == lexer.LBRACE {
			p.nextToken()
			block, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			elseBlock = block
			break
		} else {
			return nil, fmt.Errorf("unexpected token after else: %s", p.currentToken.Literal)
		}
	}

	return &IfStatement{
		Condition:   condition,
		ThenBlock:   thenBlock,
		ElseIfParts: elseIfParts,
		ElseBlock:   elseBlock,
	}, nil
}

func (p *Parser) parseBlock() ([]Expression, error) {
	var stmts []Expression

	for p.currentToken.Type != lexer.RBRACE && p.currentToken.Type != lexer.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
		if p.currentToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
	}

	if p.currentToken.Type != lexer.RBRACE {
		return nil, fmt.Errorf("expected closing '}' in block, got '%s'", p.currentToken.Literal)
	}
	p.nextToken()

	return stmts, nil
}
