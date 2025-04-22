package runtime

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
	"github.com/isaeken/brickengine-go/parser"
)

func RunTemplate(code string, ctx Context, funcs Functions) (string, error) {
	return EvalTemplate(code, ctx, funcs)
}

func RunScript(code string, ctx Context, funcs Functions) (string, error) {
	l := lexer.New(code)
	p := parser.New(l)

	statements, err := p.Parse()
	if err != nil {
		return "", err
	}

	var last interface{} = ""

	for _, stmt := range statements {
		val, err := Evaluate(stmt, ctx, funcs)
		if err != nil {
			return "", err
		}
		if IsReturn(val) {
			ret := ExtractReturn(val)
			return fmt.Sprintf("%v", ret), nil
		}
		last = val
	}

	return fmt.Sprintf("%v", last), nil
}
