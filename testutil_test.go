package navigator

import (
	"os"
	"path/filepath"
	"testing"
)

func loadFixture(t *testing.T, relPath string) []byte {
	t.Helper()
	path := filepath.Join("testdata", relPath)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", relPath, err)
	}
	return data
}

func mustParse(t *testing.T, content []byte) *Index {
	t.Helper()
	idx := Parse(content)
	if idx == nil {
		t.Fatal("Parse returned nil")
	}
	return idx
}

func writeMultiFile(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}
