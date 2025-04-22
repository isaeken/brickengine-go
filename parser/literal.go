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

type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type NullLiteral struct{}

func (n *NullLiteral) String() string {
	return "null"
}
