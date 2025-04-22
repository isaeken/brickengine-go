package parser

import "fmt"

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%v", n.Value)
}
