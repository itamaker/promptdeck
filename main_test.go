package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCartesianMatrix(t *testing.T) {
	t.Parallel()

	matrix := map[string][]string{
		"audience": {"engineers", "researchers"},
		"mode":     {"strict", "creative"},
	}

	combos := cartesianMatrix(matrix)
	if len(combos) != 4 {
		t.Fatalf("len(combos) = %d, want 4", len(combos))
	}
}

func TestRenderWithVars(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	templatePath := filepath.Join(dir, "prompt.tmpl")
	varsPath := filepath.Join(dir, "vars.json")

	if err := os.WriteFile(templatePath, []byte("Hello {{.name}}"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	if err := os.WriteFile(varsPath, []byte(`{"name":"OpenClaw"}`), 0o644); err != nil {
		t.Fatalf("write vars: %v", err)
	}

	rendered, err := renderWithVars(templatePath, varsPath)
	if err != nil {
		t.Fatalf("renderWithVars() error = %v", err)
	}
	if !strings.Contains(rendered, "OpenClaw") {
		t.Fatalf("rendered output = %q, want OpenClaw", rendered)
	}
}
