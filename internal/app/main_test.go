package app

import (
	"encoding/json"
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

func TestOptimizeSelectsBestPrompt(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	templatePath := filepath.Join(dir, "prompt.tmpl")
	matrixPath := filepath.Join(dir, "matrix.json")

	if err := os.WriteFile(templatePath, []byte("Tone={{.tone}} Risk={{.risk}}"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	matrixBody, err := json.Marshal(map[string][]string{
		"tone": {"direct", "skeptical"},
		"risk": {"low", "high"},
	})
	if err != nil {
		t.Fatalf("marshal matrix: %v", err)
	}
	if err := os.WriteFile(matrixPath, matrixBody, 0o644); err != nil {
		t.Fatalf("write matrix: %v", err)
	}

	candidates, err := buildCandidates(templatePath, matrixPath)
	if err != nil {
		t.Fatalf("buildCandidates() error = %v", err)
	}

	report := optimize(candidates, []CandidateScore{
		{Index: "001", Score: 0.72},
		{Index: "002", Score: 0.91},
		{Index: "003", Score: 0.65},
		{Index: "004", Score: 0.84},
	}, 2)

	if report.Best == nil || report.Best.Index != "002" {
		t.Fatalf("best = %+v, want index 002", report.Best)
	}
	if len(report.Top) != 2 {
		t.Fatalf("len(top) = %d, want 2", len(report.Top))
	}
	if len(report.FactorEffects) == 0 {
		t.Fatalf("expected factor effects")
	}
}
