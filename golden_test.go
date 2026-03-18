package navigator

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden files")

type goldenExpected struct {
	Version        string   `json:"version"`
	DocType        string   `json:"doctype"`
	OperationCount int      `json:"operation_count"`
	SchemaCount    int      `json:"schema_count"`
	RefCount       int      `json:"ref_count"`
	OperationIDs   []string `json:"operation_ids"`
	SchemaNames    []string `json:"schema_names"`
}

func TestGolden(t *testing.T) {
	goldenDir := filepath.Join("testdata", "golden")
	entries, err := os.ReadDir(goldenDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			specPath := filepath.Join(goldenDir, name, "specs", "api.yaml")
			expectedPath := filepath.Join(goldenDir, name, "expected", "index.json")

			content, err := os.ReadFile(specPath)
			if err != nil {
				t.Fatalf("missing spec: %v", err)
			}

			idx := Parse(content)
			if idx == nil {
				t.Fatal("parse returned nil")
			}

			actual := buildGoldenExpected(idx)

			if *updateGolden {
				data, _ := json.MarshalIndent(actual, "", "  ")
				data = append(data, '\n')
				os.WriteFile(expectedPath, data, 0644)
				t.Log("updated golden file")
				return
			}

			expectedData, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("missing expected file: %v", err)
			}

			var expected goldenExpected
			if err := json.Unmarshal(expectedData, &expected); err != nil {
				t.Fatalf("invalid expected JSON: %v", err)
			}

			if actual.Version != expected.Version {
				t.Errorf("version: got %q, want %q", actual.Version, expected.Version)
			}
			if actual.DocType != expected.DocType {
				t.Errorf("doctype: got %q, want %q", actual.DocType, expected.DocType)
			}
			if actual.OperationCount != expected.OperationCount {
				t.Errorf("operation_count: got %d, want %d", actual.OperationCount, expected.OperationCount)
			}
			if actual.SchemaCount != expected.SchemaCount {
				t.Errorf("schema_count: got %d, want %d", actual.SchemaCount, expected.SchemaCount)
			}
			if actual.RefCount != expected.RefCount {
				t.Errorf("ref_count: got %d, want %d", actual.RefCount, expected.RefCount)
			}

			sort.Strings(actual.OperationIDs)
			sort.Strings(expected.OperationIDs)
			if len(actual.OperationIDs) != len(expected.OperationIDs) {
				t.Errorf("operation_ids: got %v, want %v", actual.OperationIDs, expected.OperationIDs)
			}

			sort.Strings(actual.SchemaNames)
			sort.Strings(expected.SchemaNames)
			if len(actual.SchemaNames) != len(expected.SchemaNames) {
				t.Errorf("schema_names: got %v, want %v", actual.SchemaNames, expected.SchemaNames)
			}
		})
	}
}

func buildGoldenExpected(idx *Index) goldenExpected {
	ge := goldenExpected{
		Version:        string(idx.Version),
		OperationCount: len(idx.Operations),
		SchemaCount:    len(idx.Schemas),
		RefCount:       len(idx.AllRefs),
	}
	if idx.Version.IsZero() {
		ge.Version = ""
	}
	if idx.IsOpenAPI() {
		ge.DocType = "root"
	} else {
		ge.DocType = "fragment"
	}
	for id := range idx.Operations {
		ge.OperationIDs = append(ge.OperationIDs, id)
	}
	for name := range idx.Schemas {
		ge.SchemaNames = append(ge.SchemaNames, name)
	}
	if ge.OperationIDs == nil {
		ge.OperationIDs = []string{}
	}
	if ge.SchemaNames == nil {
		ge.SchemaNames = []string{}
	}
	sort.Strings(ge.OperationIDs)
	sort.Strings(ge.SchemaNames)
	return ge
}
