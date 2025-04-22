package runtime

import (
	"fmt"
	"strconv"
)

type ReturnedValue struct {
	Value interface{}
}

func ResolveVariable(ctx Context, parts []string) (interface{}, error) {
	var val interface{} = map[string]interface{}(ctx)
	for _, p := range parts {
		if m, ok := val.(map[string]interface{}); ok {
			val = m[p]
		} else {
			return nil, fmt.Errorf("cannot access '%s' in non-object (%T)", p, val)
		}
	}
	return val, nil
}

func EvalBinary(left interface{}, right interface{}, op string) (interface{}, error) {
	lf, lok := ToFloat(left)
	rf, rok := ToFloat(right)

	if IsComparisonOperator(op) {
		if !lok || !rok {
			return nil, fmt.Errorf("comparison operations require numeric operands")
		}
		switch op {
		case "==":
			return lf == rf, nil
		case "!=":
			return lf != rf, nil
		case "<":
			return lf < rf, nil
		case "<=":
			return lf <= rf, nil
		case ">":
			return lf > rf, nil
		case ">=":
			return lf >= rf, nil
		}
	}

	if !lok || !rok {
		return nil, fmt.Errorf("binary operations only supported on numeric types")
	}

	switch op {
	case "+":
		return lf + rf, nil
	case "-":
		return lf - rf, nil
	case "*":
		return lf * rf, nil
	case "/":
		if rf == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return lf / rf, nil
	default:
		return nil, fmt.Errorf("unsupported operator '%s'", op)
	}
}

func ToFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case float64:
		return val, true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func ToBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	if f, ok := v.(float64); ok {
		return f != 0
	}
	if s, ok := v.(string); ok {
		return s != ""
	}
	return v != nil
}

func IsReturn(v interface{}) bool {
	_, ok := v.(ReturnedValue)
	return ok
}

func ExtractReturn(val interface{}) interface{} {
	if ret, ok := val.(ReturnedValue); ok {
		return ret.Value
	}
	return val
}

func IsTruthy(val interface{}) bool {
	switch v := val.(type) {
	case nil:
		return false
	case string:
		return v != ""
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return true
	}
}

func IsComparisonOperator(operator string) bool {
	switch operator {
	case "==", "!=", "<", ">", "<=", ">=":
		return true
	default:
		return false
	}
}
