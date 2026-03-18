package navigator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func generateSpec(t *testing.T, fileCount int) string {
	t.Helper()
	dir := t.TempDir()

	// Generate root
	root := `openapi: "3.0.3"
info:
  title: Scale Test
  version: "1.0"
paths:`
	for i := 0; i < fileCount && i < 100; i++ {
		root += fmt.Sprintf(`
  /path%d:
    get:
      operationId: get%d
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./schemas/schema%d.yaml"`, i, i, i)
	}
	os.WriteFile(filepath.Join(dir, "api.yaml"), []byte(root), 0644)

	// Generate schemas directory
	os.MkdirAll(filepath.Join(dir, "schemas"), 0755)
	for i := 0; i < fileCount; i++ {
		var schema string
		if i < fileCount-1 {
			schema = fmt.Sprintf(`type: object
properties:
  id:
    type: integer
  next:
    $ref: "./schema%d.yaml"`, i+1)
		} else {
			schema = `type: object
properties:
  id:
    type: integer
  name:
    type: string`
		}
		os.WriteFile(filepath.Join(dir, "schemas", fmt.Sprintf("schema%d.yaml", i)), []byte(schema), 0644)
	}

	return dir
}

func TestScale_LazyLoading(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping scale test in short mode")
	}

	dir := generateSpec(t, 100)
	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Only root should be loaded initially
	if proj.Cache().Len() != 1 {
		t.Errorf("lazy: cache len = %d, want 1", proj.Cache().Len())
	}

	// Resolve a specific deep ref
	result, err := proj.Resolve(proj.RootURI(), "./schemas/schema0.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if result.TargetIndex == nil {
		t.Error("should resolve schema0")
	}

	loaded := proj.Cache().Len()
	if loaded >= 100 {
		t.Errorf("lazy loading should not load all files, loaded %d", loaded)
	}

	// LoadAll should load everything
	proj.LoadAll()
	if proj.Cache().Len() < 100 {
		t.Errorf("LoadAll should load at least 100 files, got %d", proj.Cache().Len())
	}
}

func TestScale_LargeGraph(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping scale test in short mode")
	}

	g := NewFileGraph()
	for i := 0; i < 1000; i++ {
		from := fmt.Sprintf("file:///%d.yaml", i)
		to := fmt.Sprintf("file:///%d.yaml", i+1)
		g.AddEdge(RefEdge{FromURI: from, ToURI: to})
	}

	// Should handle large transitive queries
	deps := g.TransitiveDependentsOf("file:///999.yaml")
	if len(deps) != 999 {
		t.Errorf("expected 999 transitive dependents, got %d", len(deps))
	}
}
