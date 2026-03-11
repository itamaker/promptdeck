package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"text/template"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}

	switch args[0] {
	case "render":
		return runRender(args[1:])
	case "matrix":
		return runMatrix(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", args[0])
		usage()
		return 2
	}
}

func usage() {
	fmt.Println("promptdeck renders reusable prompt templates.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  promptdeck render -template examples/review.tmpl -vars examples/vars.json")
	fmt.Println("  promptdeck matrix -template examples/review.tmpl -matrix examples/matrix.json")
}

func runRender(args []string) int {
	fs := flag.NewFlagSet("render", flag.ContinueOnError)
	templatePath := fs.String("template", "", "path to a Go text template")
	varsPath := fs.String("vars", "", "path to a JSON object or array of objects")
	outPath := fs.String("out", "", "optional output file")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *templatePath == "" || *varsPath == "" {
		fmt.Fprintln(os.Stderr, "both -template and -vars are required")
		return 2
	}

	body, err := renderWithVars(*templatePath, *varsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(body), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}

	fmt.Print(body)
	return 0
}

func runMatrix(args []string) int {
	fs := flag.NewFlagSet("matrix", flag.ContinueOnError)
	templatePath := fs.String("template", "", "path to a Go text template")
	matrixPath := fs.String("matrix", "", "path to a JSON object of string arrays")
	outDir := fs.String("out-dir", "", "optional output directory for rendered prompts")
	ext := fs.String("ext", ".txt", "output extension when -out-dir is provided")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *templatePath == "" || *matrixPath == "" {
		fmt.Fprintln(os.Stderr, "both -template and -matrix are required")
		return 2
	}

	rendered, err := renderMatrix(*templatePath, *matrixPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if *outDir == "" {
		for i, prompt := range rendered {
			if i > 0 {
				fmt.Println("\n---")
			}
			fmt.Print(prompt)
		}
		return 0
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	for i, prompt := range rendered {
		name := fmt.Sprintf("%03d%s", i+1, *ext)
		if err := os.WriteFile(filepath.Join(*outDir, name), []byte(prompt), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	fmt.Printf("wrote %d prompts to %s\n", len(rendered), *outDir)
	return 0
}

func renderWithVars(templatePath string, varsPath string) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	body, err := os.ReadFile(varsPath)
	if err != nil {
		return "", fmt.Errorf("read vars: %w", err)
	}

	var single map[string]any
	if err := json.Unmarshal(body, &single); err == nil && single != nil {
		return executeTemplate(tmpl, single)
	}

	var many []map[string]any
	if err := json.Unmarshal(body, &many); err != nil {
		return "", fmt.Errorf("decode vars: %w", err)
	}

	var combined string
	for i, item := range many {
		rendered, err := executeTemplate(tmpl, item)
		if err != nil {
			return "", err
		}
		if i > 0 {
			combined += "\n---\n"
		}
		combined += rendered
	}
	return combined, nil
}

func renderMatrix(templatePath string, matrixPath string) ([]string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	body, err := os.ReadFile(matrixPath)
	if err != nil {
		return nil, fmt.Errorf("read matrix: %w", err)
	}

	var matrix map[string][]string
	if err := json.Unmarshal(body, &matrix); err != nil {
		return nil, fmt.Errorf("decode matrix: %w", err)
	}

	combos := cartesianMatrix(matrix)
	rendered := make([]string, 0, len(combos))
	for _, combo := range combos {
		output, err := executeTemplate(tmpl, combo)
		if err != nil {
			return nil, err
		}
		rendered = append(rendered, output)
	}
	return rendered, nil
}

func executeTemplate(tmpl *template.Template, vars map[string]any) (string, error) {
	var output stringWriter
	if err := tmpl.Execute(&output, vars); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return output.String(), nil
}

func cartesianMatrix(matrix map[string][]string) []map[string]any {
	if len(matrix) == 0 {
		return []map[string]any{{}}
	}

	keys := make([]string, 0, len(matrix))
	for key := range matrix {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	var results []map[string]any
	var walk func(int, map[string]any)
	walk = func(index int, current map[string]any) {
		if index == len(keys) {
			item := make(map[string]any, len(current))
			for key, value := range current {
				item[key] = value
			}
			results = append(results, item)
			return
		}

		key := keys[index]
		values := matrix[key]
		if len(values) == 0 {
			current[key] = ""
			walk(index+1, current)
			delete(current, key)
			return
		}

		for _, value := range values {
			current[key] = value
			walk(index+1, current)
		}
		delete(current, key)
	}

	walk(0, map[string]any{
		"index": "0",
	})

	for i := range results {
		results[i]["index"] = strconv.Itoa(i + 1)
	}
	return results
}

type stringWriter struct {
	body []byte
}

func (w *stringWriter) Write(p []byte) (int, error) {
	w.body = append(w.body, p...)
	return len(p), nil
}

func (w *stringWriter) String() string {
	return string(w.body)
}
