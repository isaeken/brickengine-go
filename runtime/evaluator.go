package runtime

import (
	"fmt"
	"github.com/isaeken/brickengine-go/parser"
	"reflect"
	"runtime"
	"strings"
)

var LoopLimit = 100_000_000
var MaxMemoryBytes = 10 * 1024 * 1024

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
	case *parser.ArrayLiteral:
		var values []interface{}
		for _, el := range node.Elements {
			v, err := Evaluate(el, ctx, funcs)
			if err != nil {
				return nil, err
			}
			values = append(values, v)
		}
		return values, nil
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
		if varExpr, ok := node.Target.(*parser.VariableExpr); ok {
			fnName := strings.Join(varExpr.Parts, ".")

			if fn, ok := funcs[fnName]; ok {
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
		}

		targetVal, err := Evaluate(node.Target, ctx, funcs)
		if err != nil {
			return nil, err
		}

		fn := reflect.ValueOf(targetVal)
		if fn.Kind() != reflect.Func {
			return nil, fmt.Errorf("expression is not callable")
		}

		args := []reflect.Value{}
		for _, arg := range node.Args {
			val, err := Evaluate(arg, ctx, funcs)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(val))
		}

		results := fn.Call(args)
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
	case *parser.ObjectExpr:
		obj := make(map[string]interface{})
		for k, v := range node.Pairs {
			switch fn := v.(type) {
			case *parser.FnStatement:
				fnValue := DeclareFunction(ctx, funcs, fn.Args, fn.Body)
				obj[k] = fnValue
			default:
				val, err := Evaluate(fn, ctx, funcs)
				if err != nil {
					return nil, err
				}
				obj[k] = val
			}
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
		funcs[node.Name] = DeclareFunction(ctx, funcs, node.Args, node.Body)

		return nil, nil
	case *parser.ForStatement:
		count := 0

		if node.Iterable != nil {
			iterVal, err := Evaluate(node.Iterable, ctx, funcs)
			if err != nil {
				return nil, err
			}

			slice, ok := iterVal.([]interface{})
			if !ok {
				return nil, fmt.Errorf("foreach loop target must be an array")
			}

			for _, item := range slice {
				ctx[node.VarName] = item

				for _, stmt := range node.Body {
					val, err := Evaluate(stmt, ctx, funcs)
					if err != nil {
						return nil, err
					}
					if IsReturn(val) {
						return val, nil
					}
				}
			}
			return nil, nil
		}

		_, err := Evaluate(node.Init, ctx, funcs)
		if err != nil {
			return nil, err
		}
		for {
			count++
			if count > LoopLimit {
				return nil, fmt.Errorf("execution limit exceeded (possible infinite loop)")
			}
			if err := checkMemoryLimit(); err != nil {
				return nil, err
			}

			condVal, err := Evaluate(node.Condition, ctx, funcs)
			if err != nil {
				return nil, err
			}
			if !ToBool(condVal) {
				break
			}

			for _, stmt := range node.Body {
				val, err := Evaluate(stmt, ctx, funcs)
				if err != nil {
					return nil, err
				}
				if IsReturn(val) {
					return val, nil
				}
			}

			_, err = Evaluate(node.Update, ctx, funcs)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *parser.WhileStatement:
		count := 0

		for {
			count++
			if count > LoopLimit {
				return nil, fmt.Errorf("execution limit exceeded (possible infinite loop)")
			}
			if count%1000 == 0 {
				if err := checkMemoryLimit(); err != nil {
					return nil, err
				}
			}

			condVal, err := Evaluate(node.Condition, ctx, funcs)
			if err != nil {
				return nil, err
			}
			if !ToBool(condVal) {
				break
			}

			for _, stmt := range node.Body {
				val, err := Evaluate(stmt, ctx, funcs)
				if err != nil {
					return nil, err
				}
				if IsReturn(val) {
					return val, nil
				}
			}
		}
		return nil, nil
	case *parser.TryCatchStatement:
		for _, stmt := range node.TryBlock {
			val, err := Evaluate(stmt, ctx, funcs)
			if err != nil {
				for _, catchStmt := range node.CatchBlock {
					catchVal, cerr := Evaluate(catchStmt, ctx, funcs)
					if cerr != nil {
						return nil, cerr
					}
					if IsReturn(catchVal) {
						return catchVal, nil
					}
				}
				return nil, err
			}
			if IsReturn(val) {
				return val, nil
			}
		}
		return nil, nil
	case *parser.IndexAssignmentStatement:
		target, err := Evaluate(node.Target, ctx, funcs)
		if err != nil {
			return nil, err
		}

		index, err := Evaluate(node.Index, ctx, funcs)
		if err != nil {
			return nil, err
		}

		value, err := Evaluate(node.Value, ctx, funcs)
		if err != nil {
			return nil, err
		}

		slice, ok := target.([]interface{})
		if !ok {
			return nil, fmt.Errorf("assignment target must be an array")
		}

		i := int(reflect.ValueOf(index).Float())
		if i < 0 {
			return nil, fmt.Errorf("negative index not allowed")
		}

		for len(slice) <= i {
			slice = append(slice, nil)
		}
		slice[i] = value

		if varExpr, ok := node.Target.(*parser.VariableExpr); ok {
			err = AssignToContext(ctx, varExpr, slice)
			if err != nil {
				return nil, err
			}
		}

		return value, nil
	default:
		return nil, fmt.Errorf("unknown expression type %T", node)
	}
}

func DeclareFunction(ctx Context, funcs Functions, Args []string, Body []parser.Expression) interface{} {
	return func(args ...interface{}) interface{} {
		localCtx := make(Context)
		for k, v := range ctx {
			localCtx[k] = v
		}
		for i, param := range Args {
			if i < len(args) {
				localCtx[param] = args[i]
			}
		}

		for _, stmt := range Body {
			val, err := Evaluate(stmt, localCtx, funcs)
			if err != nil {
				panic(err)
			}
			if IsReturn(val) {
				return ExtractReturn(val)
			}
		}

		return nil
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

func checkMemoryLimit() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	if memStats.Alloc > uint64(MaxMemoryBytes) {
		return fmt.Errorf("memory limit exceeded (%.2f MB used)", float64(memStats.Alloc))
	}

	return nil
}
