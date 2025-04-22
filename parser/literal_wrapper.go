package parser

import (
	"fmt"
)

type LiteralWrapper struct {
	Value interface{}
}

func (l *LiteralWrapper) Pos() int {
	return 0
}

func (l *LiteralWrapper) End() int {
	return 0
}

func (l *LiteralWrapper) String() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l *LiteralWrapper) NodeType() Node {
	return l
}

func (l *LiteralWrapper) Evaluate() (interface{}, error) {
	return l.Value, nil
}
