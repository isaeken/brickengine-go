package parser

type Node interface{}

type Expression Node

type Identifier struct {
	Value string
}

type StringLiteral struct {
	Value string
}

type NumberLiteral struct {
	Value float64
}

type VariableExpr struct {
	Parts []string
}

type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

type CallExpr struct {
	Name string
	Args []Expression
}

type PipeExpr struct {
	Left  Expression
	Right Expression
}

type IndexExpr struct {
	Target Expression
	Index  Expression
}
