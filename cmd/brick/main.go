package main

import (
	"fmt"
	"github.com/isaeken/brickengine-go/runtime"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: brick <file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	ctx := runtime.Context{}
	funcs := runtime.DefaultFunctions()

	output, err := runtime.RunScript(string(content), ctx, funcs)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Println(output)
}
