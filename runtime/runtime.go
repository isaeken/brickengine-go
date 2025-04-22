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
			return formatOutput(ret), nil
		}
		last = val
	}

	return formatOutput(last), nil
}

func formatOutput(output interface{}) string {
	switch v := output.(type) {
	case float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", output)
	}
}
