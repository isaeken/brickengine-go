# brickengine-go

a lightweight, go-powered expression parser and template engine.

features:
- arithmetic: `5 + 5`, `gb(1) - mb(512) + 128`
- pipe fallback: `var.vlaue | "default"`
- nested calls `slug(var.id)`
- variable resolution: `var.user.name`
- built-in `{{ .. }}` template rendering

## installation

```bash
go get github.com/isaeken/brickengine-go
```

## usage

```go
input := `
disk_size: {{ gb(1) + var.extra | 512 }}
memory: {{ mb(256) + 128 }}
text: Hello {{ var.name | "Guest" }}
`

ctx := runtime.Context{
    "var": map[string]interface{}{
        "extra": 128,
        "name":  "İsa",
    },
}
funcs := runtime.Functions{
    "gb": func(x float64) float64 { return x * 1024 },
    "mb": func(x float64) float64 { return x },
}

output, _ := runtime.EvalTemplate(input, ctx, funcs)
fmt.Println(output)
```

## license

MIT
