package runtime

import (
	"fmt"
	"github.com/isaeken/brickengine-go/parser"
	"reflect"
)

type Context map[string]interface{}

type Functions map[string]interface{}

func Evaluate(expr parser.Expression, ctx Context, funcs Functions) (interface{}, error) {
	switch node := expr.(type) {
	case *parser.StringLiteral:
		return node.Value, nil
	case *parser.NumberLiteral:
		return node.Value, nil
	case *parser.VariableExpr:
		return ResolveVariable(ctx, node.Parts)
	case *parser.BinaryExpr:
		left, err := Evaluate(node.Left, ctx, funcs)
		if err != nil {
			return nil, err
		}

		right, err := Evaluate(node.Right, ctx, funcs)
		if err != nil {
			return nil, err
		}

		return EvalBinary(left, right, node.Operator)
	case *parser.CallExpr:
		name := node.Name
		fn, ok := funcs[name]
		if !ok {
			return nil, fmt.Errorf("undefined function '%s'", node.Name)
		}
		args := []reflect.Value{}
		for _, arg := range node.Args {
			val, err := Evaluate(arg, ctx, funcs)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(val))
		}
		results := reflect.ValueOf(fn).Call(args)
		return results[0].Interface(), nil
	case *parser.PipeExpr:
		leftVal, err := Evaluate(node.Left, ctx, funcs)
		if err != nil {
			return Evaluate(node.Right, ctx, funcs)
		}

		if IsTruthy(leftVal) {
			return leftVal, nil
		}

		return Evaluate(node.Right, ctx, funcs)
	case *parser.IndexExpr:
		target, err := Evaluate(node.Target, ctx, funcs)
		if err != nil {
			return nil, err
		}
		index, err := Evaluate(node.Index, ctx, funcs)
		if err != nil {
			return nil, err
		}
		arr := reflect.ValueOf(target)
		i := int(reflect.ValueOf(index).Float())
		if arr.Kind() == reflect.Slice && i >= 0 && i < arr.Len() {
			return arr.Index(i).Interface(), nil
		}
		return nil, fmt.Errorf("index out of range")
	default:
		return nil, fmt.Errorf("unknown expression type %T", node)
	}
}
