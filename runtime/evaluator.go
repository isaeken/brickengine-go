package runtime

import (
	"fmt"
	"github.com/isaeken/brickengine-go/parser"
	"reflect"
)

type Context map[string]interface{}

type Functions map[string]interface{}

type FunctionMap map[string]*parser.FnStatement

func Evaluate(expr parser.Expression, ctx Context, funcs Functions) (interface{}, error) {
	switch node := expr.(type) {
	case *parser.StringLiteral:
		return node.Value, nil
	case *parser.NumberLiteral:
		return node.Value, nil
	case *parser.BoolLiteral:
		return node.Value, nil
	case *parser.NullLiteral:
		return nil, nil
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
		// native fn
		if fn, ok := funcs[node.Name]; ok {
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
		}

		// user-defined fn
		if def, ok := ctx[node.Name].(*parser.FnStatement); ok {
			newCtx := make(Context)
			for i, name := range def.Args {
				argVal, err := Evaluate(node.Args[i], ctx, funcs)
				if err != nil {
					return nil, err
				}
				newCtx[name] = argVal
			}
			for _, stmt := range def.Body {
				val, err := Evaluate(stmt, newCtx, funcs)
				if err != nil {
					return nil, err
				}
				if IsReturn(val) {
					return ExtractReturn(val), nil
				}
			}
			return nil, nil
		}

		return nil, fmt.Errorf("undefined function '%s'", node.Name)
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
	case *parser.ObjectExpr:
		obj := make(map[string]interface{})
		for k, v := range node.Pairs {
			val, err := Evaluate(v, ctx, funcs)
			if err != nil {
				return nil, err
			}
			obj[k] = val
		}
		return obj, nil
	case *parser.AssignmentStmt:
		val, err := Evaluate(node.Value, ctx, funcs)
		if err != nil {
			return nil, err
		}
		err = AssignToContext(ctx, node.Target, val)
		if err != nil {
			return nil, err
		}
		return val, nil
	case *parser.IfStatement:
		condVal, err := Evaluate(node.Condition, ctx, funcs)
		if err != nil {
			return nil, err
		}
		if ToBool(condVal) {
			for _, stmt := range node.ThenBlock {
				val, err := Evaluate(stmt, ctx, funcs)
				if err != nil {
					return nil, err
				}
				if IsReturn(val) {
					return val, nil
				}
			}
			return nil, nil
		}
		for _, elseif := range node.ElseIfParts {
			condVal, err := Evaluate(elseif.Condition, ctx, funcs)
			if err != nil {
				return nil, err
			}
			if ToBool(condVal) {
				for _, stmt := range elseif.Block {
					val, err := Evaluate(stmt, ctx, funcs)
					if err != nil {
						return nil, err
					}
					if IsReturn(val) {
						return val, nil
					}
				}
				return nil, nil
			}
		}
		for _, stmt := range node.ElseBlock {
			val, err := Evaluate(stmt, ctx, funcs)
			if err != nil {
				return nil, err
			}
			if IsReturn(val) {
				return val, nil
			}
		}
		return nil, nil
	case *parser.ReturnStatement:
		val, err := Evaluate(node.Value, ctx, funcs)
		if err != nil {
			return nil, err
		}
		return ReturnedValue{Value: val}, nil
	case *parser.LetStatement:
		val, err := Evaluate(node.Value, ctx, funcs)
		if err != nil {
			return "", err
		}
		ctx[node.Name] = val
		return val, nil
	case *parser.FnStatement:
		ctx[node.Name] = node
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown expression type %T", node)
	}
}

func AssignToContext(ctx Context, target parser.Expression, value interface{}) error {
	varExpr, ok := target.(*parser.VariableExpr)
	if !ok {
		return fmt.Errorf("assignment target must be variable")
	}

	cur := ctx
	for i := 0; i < len(varExpr.Parts)-1; i++ {
		key := varExpr.Parts[i]
		child, ok := cur[key].(map[string]interface{})
		if !ok {
			child = make(map[string]interface{})
			cur[key] = child
		}
		cur = child
	}
	cur[varExpr.Parts[len(varExpr.Parts)-1]] = value
	return nil
}
