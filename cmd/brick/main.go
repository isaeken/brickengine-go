package main

import (
	"fmt"
	"github.com/isaeken/brickengine-go/runtime"
)

func main() {
	code := `slug(ips[0]) | "a\t\"sd\""`

	ctx := runtime.Context{
		"var": map[string]interface{}{
			"extra": 128,
			"name":  "İsa",
		},
		"items": []string{"a", "b", "c"},
	}

	funcs := runtime.Functions{
		"gb": func(x float64) float64 { return x * 1024 },
		"mb": func(x float64) float64 { return x },
	}

	input := fmt.Sprintf("size: {{ %s }}", code)
	output, err := runtime.EvalTemplate(input, ctx, funcs)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
