package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sailpoint-oss/navigator"
)

func TestWriteReport_IncludesFixtureDiagnosticsAndPasses(t *testing.T) {
	repoRoot := t.TempDir()
	testdataDir := filepath.Join(repoRoot, "testdata")
	if err := os.MkdirAll(testdataDir, 0o755); err != nil {
		t.Fatalf("mkdir testdata: %v", err)
	}

	valid := `openapi: "3.1.0"
info:
  title: Good API
  version: "1.0.0"
paths: {}
`
	invalid := `openapi: "3.1.0"
paths: {}
`
	if err := os.WriteFile(filepath.Join(testdataDir, "good.yaml"), []byte(valid), 0o644); err != nil {
		t.Fatalf("write good fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testdataDir, "missing-info.yaml"), []byte(invalid), 0o644); err != nil {
		t.Fatalf("write invalid fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testdataDir, "notes.txt"), []byte("ignored"), 0o644); err != nil {
		t.Fatalf("write ignored fixture: %v", err)
	}

	var out bytes.Buffer
	if err := writeReport(&out, repoRoot); err != nil {
		t.Fatalf("writeReport: %v", err)
	}
	report := out.String()

	for _, want := range []string{
		"## Navigator validation preview",
		"### `testdata/good.yaml`",
		"### `testdata/missing-info.yaml`",
		"_No diagnostics._",
		"`structural.missing-info`",
		"validation-preview-comment.yml",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
	if strings.Contains(report, "notes.txt") {
		t.Fatalf("report should ignore non-YAML files:\n%s", report)
	}
}

func TestWriteIssueTable_TruncatesLargeIssueSets(t *testing.T) {
	issues := make([]navigator.Issue, 0, maxIssuesPerFile+1)
	for i := 0; i < maxIssuesPerFile+1; i++ {
		issues = append(issues, navigator.Issue{
			Code:     "structural.test",
			Message:  "issue message",
			Pointer:  "/paths/~1example",
			Severity: navigator.SeverityError,
			Category: navigator.CategoryStructural,
		})
	}

	var out bytes.Buffer
	writeIssueTable(&out, issues)
	report := out.String()

	if count := strings.Count(report, "`structural.test`"); count != maxIssuesPerFile {
		t.Fatalf("expected %d rendered issues, got %d", maxIssuesPerFile, count)
	}
	if !strings.Contains(report, "_Showing first 40 issues._") {
		t.Fatalf("expected truncation notice, got:\n%s", report)
	}
}
