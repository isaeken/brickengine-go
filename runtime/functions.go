package runtime

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

func UUIDAndFormatFunctions() Functions {
	return Functions{
		"uuid":   func() string { return uuid.NewString() },
		"slug":   func(s string) string { return slug.Make(s) },
		"random": func() float64 { return rand.Float64() },
		"now":    func() string { return time.Now().UTC().Format(time.RFC3339) },
		"format": func(v interface{}) string { return fmt.Sprintf("%v", v) },
		"to_json": func(v interface{}) string {
			data, _ := json.Marshal(v)
			return string(data)
		},
		"parse_json": func(s string) interface{} {
			var out interface{}
			json.Unmarshal([]byte(s), &out)
			return out
		},
	}
}

func StringFunctions() Functions {
	return Functions{
		"strlen":          func(s string) float64 { return float64(len(s)) },
		"str_upper":       strings.ToUpper,
		"str_lower":       strings.ToLower,
		"str_trim":        strings.TrimSpace,
		"str_contains":    strings.Contains,
		"str_starts_with": strings.HasPrefix,
		"str_ends_with":   strings.HasSuffix,
		"str_replace":     strings.ReplaceAll,
		"substr": func(s string, start, length float64) string {
			r := []rune(s)
			end := int(start) + int(length)
			if end > len(r) {
				end = len(r)
			}
			return string(r[int(start):end])
		},
		"split": strings.Split,
		"join": func(arr []interface{}, sep string) string {
			var parts []string
			for _, el := range arr {
				parts = append(parts, fmt.Sprint(el))
			}
			return strings.Join(parts, sep)
		},
		"repeat": func(s string, n float64) string {
			return strings.Repeat(s, int(n))
		},
		"str_reverse": func(s string) string {
			r := []rune(s)
			for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
				r[i], r[j] = r[j], r[i]
			}
			return string(r)
		},
	}
}

func MathFunctions() Functions {
	return Functions{
		"abs":   math.Abs,
		"round": math.Round,
		"floor": math.Floor,
		"ceil":  math.Ceil,
		"min":   math.Min,
		"max":   math.Max,
		"sqrt":  math.Sqrt,
		"pow":   math.Pow,
	}
}

func TypeFunctions() Functions {
	return Functions{
		"type_of": func(v interface{}) string {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				return "array"
			case reflect.Map:
				return "object"
			case reflect.String:
				return "string"
			case reflect.Float64, reflect.Float32, reflect.Int, reflect.Int64:
				return "number"
			case reflect.Bool:
				return "boolean"
			default:
				return "unknown"
			}
		},
	}
}

func ArrayFunctions() Functions {
	return Functions{
		"count": func(arr []interface{}) float64 {
			return float64(len(arr))
		},
		"push": func(arr []interface{}, val interface{}) []interface{} {
			return append(arr, val)
		},
		"pop": func(arr []interface{}) []interface{} {
			if len(arr) == 0 {
				return arr
			}
			return arr[:len(arr)-1]
		},
		"shift": func(arr []interface{}) []interface{} {
			if len(arr) == 0 {
				return arr
			}
			return arr[1:]
		},
		"unshift": func(arr []interface{}, val interface{}) []interface{} {
			return append([]interface{}{val}, arr...)
		},
		"includes": func(arr []interface{}, val interface{}) bool {
			for _, v := range arr {
				if reflect.DeepEqual(v, val) {
					return true
				}
			}
			return false
		},
		"index_of": func(arr []interface{}, val interface{}) float64 {
			for i, v := range arr {
				if reflect.DeepEqual(v, val) {
					return float64(i)
				}
			}
			return -1
		},
		"reverse": func(arr []interface{}) []interface{} {
			n := len(arr)
			res := make([]interface{}, n)
			for i := range arr {
				res[n-1-i] = arr[i]
			}
			return res
		},
		"sort": func(arr []interface{}) []interface{} {
			clone := make([]float64, 0)
			for _, v := range arr {
				if num, ok := v.(float64); ok {
					clone = append(clone, num)
				}
			}
			sort.Float64s(clone)
			res := make([]interface{}, len(clone))
			for i, v := range clone {
				res[i] = v
			}
			return res
		},
		"slice": func(arr []interface{}, start, end float64) []interface{} {
			s := int(start)
			e := int(end)
			if s < 0 || e > len(arr) || s > e {
				return []interface{}{}
			}
			return arr[s:e]
		},
		"concat": func(arr1, arr2 []interface{}) []interface{} {
			return append(arr1, arr2...)
		},
	}
}

func UtilFunctions() Functions {
	return Functions{
		"type_of": func(v interface{}) string {
			return reflect.TypeOf(v).String()
		},
	}
}

func DefaultFunctions() Functions {
	return mergeFunctions([]Functions{
		UUIDAndFormatFunctions(),
		StringFunctions(),
		MathFunctions(),
		TypeFunctions(),
		ArrayFunctions(),
		UtilFunctions(),
	}...)
}

func mergeFunctions(fns ...Functions) Functions {
	merged := Functions{}
	for _, fn := range fns {
		for k, v := range fn {
			merged[k] = v
		}
	}
	return merged
}
