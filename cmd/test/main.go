package main

import (
	"fmt"
	"github.com/isaeken/brickengine-go/runtime"
	"os"
	"path/filepath"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func main() {
	scriptDirs := []string{"examples/basic", "examples/network"}
	templateDirs := []string{"examples/templates"}
	total := 0
	passed := 0

	for _, dir := range scriptDirs {
		files, _ := filepath.Glob(filepath.Join(dir, "*.bee"))
		for _, file := range files {
			total++
			fmt.Printf("🔍 %-40s ", file)

			content, _ := os.ReadFile(file)
			ctx := runtime.Context{}
			funcs := runtime.DefaultFunctions()

			result, err := runtime.RunScript(string(content), ctx, funcs)
			status := getStatus(err)
			fmt.Println(status)

			if err == nil {
				passed++
				checkGolden(file, result)
			}
		}
	}

	for _, dir := range templateDirs {
		files, _ := filepath.Glob(filepath.Join(dir, "*.yaml"))
		for _, file := range files {
			total++
			fmt.Printf("🧾 %-40s ", file)

			content, _ := os.ReadFile(file)
			ctx := runtime.Context{}
			funcs := runtime.DefaultFunctions()

			result, err := runtime.RunTemplate(string(content), ctx, funcs)
			status := getStatus(err)
			fmt.Println(status)

			if err == nil {
				passed++
				checkGolden(file, result)
			}
		}
	}

	fmt.Printf("\n📊 Test Results: %d / %d passed\n", passed, total)
	if passed != total {
		os.Exit(1)
	}
}

func getStatus(err error) string {
	if err != nil {
		return fmt.Sprintf("%s❌ Failed: %v%s", red, err, reset)
	}
	return fmt.Sprintf("%s✅ Passed%s", green, reset)
}

func checkGolden(file string, result interface{}) {
	goldenPath := strings.TrimSuffix(file, filepath.Ext(file)) + ".golden"
	resultStr := fmt.Sprint(result)

	existing, _ := os.ReadFile(goldenPath)
	golden := string(existing)

	if golden == "" {
		fmt.Println("    ⚠️  No golden file found")
		printIndentedOutput(resultStr)
		return
	}

	actual := strings.Trim(resultStr, "\n")
	expected := strings.Trim(golden, "\n")

	if actual != expected {
		fmt.Println("    ❌ Mismatch with golden output")
		printDiff(golden, resultStr)
		return
	}
}

func printIndentedOutput(output string) {
	fmt.Println("    Output:")
	for _, line := range strings.Split(output, "\n") {
		fmt.Printf("    %s\n", line)
	}
	fmt.Println()
}

func printDiff(expected, actual string) {
	ex := strings.Split(expected, "\n")
	ac := strings.Split(actual, "\n")
	fmt.Println("    Diff:")
	for i := 0; i < len(ex) || i < len(ac); i++ {
		if i < len(ex) && i < len(ac) {
			if ex[i] != ac[i] {
				fmt.Printf("    %s- %s%s\n", red, ex[i], reset)
				fmt.Printf("    %s+ %s%s\n", green, ac[i], reset)
			}
		} else if i < len(ex) {
			fmt.Printf("    %s- %s%s\n", red, ex[i], reset)
		} else {
			fmt.Printf("    %s+ %s%s\n", green, ac[i], reset)
		}
	}
	fmt.Println()
}
