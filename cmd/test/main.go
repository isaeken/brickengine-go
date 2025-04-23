package main

import (
	"fmt"
	"github.com/isaeken/brickengine-go/runtime"
	"os"
	"path/filepath"
	"regexp"
	rn "runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func main() {
	scriptDirs := []string{"examples/basic", "examples/benchmarks", "examples/fails"}
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

			debug.FreeOSMemory()
			var memStart, memEnd rn.MemStats
			rn.ReadMemStats(&memStart)

			start := time.Now()
			result, err := runtime.RunScript(string(content), ctx, funcs)
			duration := time.Since(start)

			rn.ReadMemStats(&memEnd)
			memUsage := memEnd.Alloc - memStart.Alloc

			isFailTest := strings.Contains(file, "/fails/")
			check := checkGolden(file, result)

			if isFailTest {
				if err != nil || check {
					fmt.Printf("%s✅ Expected Fail%s [%s, %.2f KB]\n", green, reset, formatDuration(duration), float64(memUsage)/1024)
					passed++
				} else {
					fmt.Printf("%s❌ Unexpected Pass%s [%s, %.2f KB]\n", red, reset, formatDuration(duration), float64(memUsage)/1024)
				}
			} else {
				if err != nil || !check {
					fmt.Printf("%s❌ Failed: %v%s [%s, %.2f KB]\n", red, err, reset, formatDuration(duration), float64(memUsage)/1024)
				} else {
					fmt.Printf("%s✅ Passed%s [%s, %.2f KB]\n", green, reset, formatDuration(duration), float64(memUsage)/1024)
					passed++
				}
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

			debug.FreeOSMemory()
			var memStart, memEnd rn.MemStats
			rn.ReadMemStats(&memStart)

			start := time.Now()
			result, err := runtime.RunTemplate(string(content), ctx, funcs)
			duration := time.Since(start)

			rn.ReadMemStats(&memEnd)
			memUsage := memEnd.Alloc - memStart.Alloc
			check := checkGolden(file, result)

			if err != nil || !check {
				fmt.Printf("%s❌ Failed: %v%s [%s, %.2f KB]\n", red, err, reset, formatDuration(duration), float64(memUsage)/1024)
			} else {
				fmt.Printf("%s✅ Passed%s [%s, %.2f KB]\n", green, reset, formatDuration(duration), float64(memUsage)/1024)
				passed++
			}
		}
	}

	fmt.Printf("\n📊 Test Results: %d / %d passed\n", passed, total)
	if passed != total {
		os.Exit(1)
	}
}

func formatDuration(d time.Duration) string {
	ms := d.Milliseconds()
	if ms < 1 {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%dms", ms)
}

func checkGolden(file string, result interface{}) bool {
	goldenPath := strings.TrimSuffix(file, filepath.Ext(file)) + ".golden"
	resultStr := fmt.Sprint(result)

	existing, _ := os.ReadFile(goldenPath)
	golden := string(existing)

	if golden == "" {
		fmt.Println("    ⚠️  No golden file found: " + goldenPath)
		printIndentedOutput(resultStr)
		return true
	}

	actual := normalizeGoldenOutput(strings.Trim(resultStr, "\n"))
	expected := normalizeGoldenOutput(strings.Trim(golden, "\n"))

	if actual != expected {
		fmt.Println("    ❌ Mismatch with golden output")
		printDiff(golden, resultStr)
		return false
	}

	return true
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
		if i < len(ex) && i < len(ac) && ex[i] != ac[i] {
			fmt.Printf("    %s- %s%s\n", red, ex[i], reset)
			fmt.Printf("    %s+ %s%s\n", green, ac[i], reset)
		} else if i < len(ex) {
			fmt.Printf("    %s- %s%s\n", red, ex[i], reset)
		} else if i < len(ac) {
			fmt.Printf("    %s+ %s%s\n", green, ac[i], reset)
		}
	}
	fmt.Println()
}

func normalizeGoldenOutput(output string) string {
	reUUID := regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)
	reTime := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)

	output = reUUID.ReplaceAllString(output, "UUID")
	output = reTime.ReplaceAllString(output, "NOW")

	return output
}
