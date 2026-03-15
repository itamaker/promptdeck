package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

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
	manifestPath := fs.String("manifest", "", "optional output manifest path")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *templatePath == "" || *matrixPath == "" {
		fmt.Fprintln(os.Stderr, "both -template and -matrix are required")
		return 2
	}

	candidates, err := buildCandidates(*templatePath, *matrixPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if *manifestPath != "" {
		body, err := json.MarshalIndent(candidates, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if err := os.WriteFile(*manifestPath, append(body, '\n'), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	if *outDir == "" {
		for i, prompt := range promptsFromCandidates(candidates) {
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
	for _, candidate := range candidates {
		name := fmt.Sprintf("%s%s", candidate.Index, *ext)
		if err := os.WriteFile(filepath.Join(*outDir, name), []byte(candidate.Prompt), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	fmt.Printf("wrote %d prompts to %s\n", len(candidates), *outDir)
	return 0
}

func runOptimize(args []string) int {
	fs := flag.NewFlagSet("optimize", flag.ContinueOnError)
	templatePath := fs.String("template", "", "path to a Go text template")
	matrixPath := fs.String("matrix", "", "path to a JSON object of string arrays")
	scoresPath := fs.String("scores", "", "path to a JSON array of prompt scores")
	topN := fs.Int("top", 3, "number of ranked prompts to return")
	outPath := fs.String("out", "", "optional path to write the best prompt")
	jsonOutput := fs.Bool("json", false, "emit machine-readable JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *templatePath == "" || *matrixPath == "" || *scoresPath == "" {
		fmt.Fprintln(os.Stderr, "template, matrix, and scores are required")
		return 2
	}

	candidates, err := buildCandidates(*templatePath, *matrixPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	scores, err := loadScores(*scoresPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	report := optimize(candidates, scores, *topN)
	if report.Best != nil && *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(report.Best.Prompt), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	if *jsonOutput {
		body, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Println(string(body))
		return 0
	}

	printOptimizationReport(report)
	return 0
}
