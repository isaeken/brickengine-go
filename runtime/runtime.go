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

	for _, stmt := range statements {
		switch node := stmt.(type) {
		case *parser.ReturnStatement:
			val, err := Evaluate(node.Value, ctx, funcs)
			if err != nil {
				return "", err
			}
			return fmt.Sprint(val), nil
		case *parser.LetStatement:
			val, err := Evaluate(node.Value, ctx, funcs)
			if err != nil {
				return "", err
			}
			ctx[node.Name] = val
		default:
			_, err := Evaluate(node, ctx, funcs)
			if err != nil {
				return "", err
			}
		}
	}

	return "", nil
}
