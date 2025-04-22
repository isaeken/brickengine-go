package parser

import "fmt"

type IndexAssignmentStatement struct {
	Target Expression
	Index  Expression
	Value  Expression
}

func (s *IndexAssignmentStatement) String() string {
	return fmt.Sprintf("%s[%s] = %s", s.Target, s.Index, s.Value)
}
