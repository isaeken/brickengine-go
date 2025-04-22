package parser

type Node interface{}

type Expression interface {
	String() string
}

type Identifier struct {
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}
