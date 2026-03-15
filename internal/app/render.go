package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

func renderWithVars(templatePath string, varsPath string) (string, error) {
	tmpl, err := parseTemplate(templatePath)
	if err != nil {
		return "", err
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
	candidates, err := buildCandidates(templatePath, matrixPath)
	if err != nil {
		return nil, err
	}
	return promptsFromCandidates(candidates), nil
}

func buildCandidates(templatePath string, matrixPath string) ([]PromptCandidate, error) {
	tmpl, err := parseTemplate(templatePath)
	if err != nil {
		return nil, err
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
	candidates := make([]PromptCandidate, 0, len(combos))
	for _, combo := range combos {
		output, err := executeTemplate(tmpl, combo)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, PromptCandidate{
			Index:  fmt.Sprint(combo["index"]),
			Vars:   cloneVars(combo),
			Prompt: output,
		})
	}
	return candidates, nil
}

func parseTemplate(templatePath string) (*template.Template, error) {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"join":  strings.Join,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
	}).ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	return tmpl, nil
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
		return []map[string]any{{"index": "001"}}
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

	walk(0, map[string]any{})

	for i := range results {
		results[i]["index"] = fmt.Sprintf("%03d", i+1)
	}
	return results
}

func promptsFromCandidates(candidates []PromptCandidate) []string {
	prompts := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		prompts = append(prompts, candidate.Prompt)
	}
	return prompts
}

func cloneVars(vars map[string]any) map[string]any {
	out := make(map[string]any, len(vars))
	for key, value := range vars {
		out[key] = value
	}
	return out
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
