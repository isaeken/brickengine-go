package runtime

import (
	"fmt"
	"strconv"
)

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
