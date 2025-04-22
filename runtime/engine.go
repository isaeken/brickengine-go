package runtime

import (
	"fmt"
	"github.com/isaeken/brickengine-go/lexer"
	"github.com/isaeken/brickengine-go/parser"
	"regexp"
	"strings"
)

var exprRegex = regexp.MustCompile(`\{\{\s*(.*?)\s*\}\}`)

func EvalTemplate(input string, ctx Context, funcs Functions) (string, error) {
	matches := exprRegex.FindAllStringSubmatchIndex(input, -1)
	if len(matches) == 0 {
		return input, nil
	}

	var result strings.Builder
	last := 0

	for _, match := range matches {
		fullStart, fullEnd := match[0], match[1]
		exprStart, exprEnd := match[2], match[3]
		exprText := input[exprStart:exprEnd]

		result.WriteString(input[last:fullStart])

		l := lexer.New(exprText)
		p := parser.New(l)
		expr, err := p.ParseExpression()
		if err != nil {
			return "", fmt.Errorf("parse error in '{{ %s }}': %w", exprText, err)
		}

		val, err := Evaluate(expr, ctx, funcs)
		if err != nil {
			return "", fmt.Errorf("evaluation error in '{{ %s }}': %w", exprText, err)
		}

		result.WriteString(fmt.Sprint(val))
		last = fullEnd
	}

	result.WriteString(input[last:])
	return result.String(), nil
}
