package navigator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscovery_Scan(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
paths: {}`,
		"schemas/Pet.yaml": `type: object
properties:
  name:
    type: string`,
		"readme.md": "# Readme",
	})

	d := NewDiscovery(nil)
	if err := d.Scan(dir); err != nil {
		t.Fatal(err)
	}

	roots := d.Roots()
	if len(roots) != 1 {
		t.Errorf("expected 1 root, got %d", len(roots))
	}

	files := d.AllFiles()
	yamlCount := 0
	for _, f := range files {
		if f.Path != "" {
			yamlCount++
		}
	}
	if yamlCount < 2 {
		t.Errorf("expected at least 2 yaml files, got %d", yamlCount)
	}
}

func TestDiscovery_RootDetection(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"root.yaml":     `openapi: "3.1.0"` + "\ninfo:\n  title: T\n  version: '1'\npaths: {}",
		"fragment.yaml": `type: object`,
	})

	d := NewDiscovery(nil)
	d.Scan(dir)

	rootURI := PathToURI(filepath.Join(dir, "root.yaml"))
	fragURI := PathToURI(filepath.Join(dir, "fragment.yaml"))

	if !d.IsRoot(rootURI) {
		t.Error("root.yaml should be detected as root")
	}
	if d.IsRoot(fragURI) {
		t.Error("fragment.yaml should not be root")
	}
}

func TestDiscovery_ExcludePatterns(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml":        `openapi: "3.0.0"` + "\ninfo:\n  title: T\n  version: '1'\npaths: {}",
		"vendor/lib.yaml": `openapi: "3.0.0"` + "\ninfo:\n  title: V\n  version: '1'\npaths: {}",
	})

	d := NewDiscovery(nil)
	d.Scan(dir)

	// vendor is auto-excluded
	roots := d.Roots()
	if len(roots) != 1 {
		t.Errorf("expected 1 root (vendor excluded), got %d", len(roots))
	}
}

func TestDiscovery_UpdateFile(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `type: object`,
	})

	d := NewDiscovery(nil)
	d.Scan(dir)

	if len(d.Roots()) != 0 {
		t.Error("initially should have no roots")
	}

	path := filepath.Join(dir, "api.yaml")
	os.WriteFile(path, []byte(`openapi: "3.0.0"`+"\ninfo:\n  title: T\n  version: '1'\npaths: {}"), 0644)
	d.UpdateFile(path)

	if len(d.Roots()) != 1 {
		t.Error("after update should have 1 root")
	}
}

func TestDiscovery_RemoveFile(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.0"` + "\ninfo:\n  title: T\n  version: '1'\npaths: {}",
	})

	d := NewDiscovery(nil)
	d.Scan(dir)

	if len(d.Roots()) != 1 {
		t.Fatal("expected 1 root")
	}

	d.RemoveFile(filepath.Join(dir, "api.yaml"))
	if len(d.Roots()) != 0 {
		t.Error("after remove should have 0 roots")
	}
}

func TestDiscovery_FileByURI(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.0"` + "\ninfo:\n  title: T\n  version: '1'\npaths: {}",
	})

	d := NewDiscovery(nil)
	d.Scan(dir)

	uri := PathToURI(filepath.Join(dir, "api.yaml"))
	f := d.FileByURI(uri)
	if f == nil {
		t.Fatal("should find file by URI")
	}
	if f.Role != RoleRoot {
		t.Errorf("role = %v, want root", f.Role)
	}
}
