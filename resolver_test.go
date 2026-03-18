package navigator

import "testing"

func TestCrossFileResolver_LocalRef(t *testing.T) {
	cache := NewIndexCache()
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	uri := "file:///api.yaml"
	cache.Set(uri, idx)

	r := NewCrossFileResolver(cache, NewFileGraph())
	result, err := r.Resolve(uri, "#/components/schemas/Pet")
	if err != nil {
		t.Fatal(err)
	}
	if result.TargetURI != uri {
		t.Errorf("TargetURI = %q", result.TargetURI)
	}
	if _, ok := result.Value.(*Schema); !ok {
		t.Errorf("expected *Schema, got %T", result.Value)
	}
}

func TestCrossFileResolver_ExternalRef(t *testing.T) {
	dir := writeMultiFile(t, map[string]string{
		"api.yaml": `openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: "./schemas/Pet.yaml"`,
		"schemas/Pet.yaml": `type: object
properties:
  name:
    type: string`,
	})

	proj, err := OpenProject(dir + "/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	result, err := proj.Resolve(proj.RootURI(), "./schemas/Pet.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if result.TargetIndex == nil {
		t.Error("expected target index")
	}
}

func TestCrossFileResolver_CycleDetection(t *testing.T) {
	cache := NewIndexCache()
	graph := NewFileGraph()

	aContent := []byte(`type: object
properties:
  partner:
    $ref: "./b.yaml"`)
	bContent := []byte(`type: object
properties:
  partner:
    $ref: "./a.yaml"`)

	aURI := "file:///workspace/a.yaml"
	bURI := "file:///workspace/b.yaml"

	cache.Set(aURI, ParseURI(aURI, aContent))
	cache.Set(bURI, ParseURI(bURI, bContent))
	graph.AddEdge(RefEdge{FromURI: aURI, ToURI: bURI})
	graph.AddEdge(RefEdge{FromURI: bURI, ToURI: aURI})

	r := NewCrossFileResolver(cache, graph)

	// First hop a->b should work
	result, err := r.Resolve(aURI, "./b.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if result.TargetURI != bURI {
		t.Errorf("first hop target = %q", result.TargetURI)
	}
}

func TestCrossFileResolver_MissingFile(t *testing.T) {
	cache := NewIndexCache()
	r := NewCrossFileResolver(cache, NewFileGraph())
	_, err := r.Resolve("file:///api.yaml", "./missing.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCrossFileResolver_EmptyRef(t *testing.T) {
	r := NewCrossFileResolver(NewIndexCache(), NewFileGraph())
	_, err := r.Resolve("file:///api.yaml", "")
	if err == nil {
		t.Error("expected error for empty ref")
	}
}

func TestCrossFileResolver_CanResolve(t *testing.T) {
	cache := NewIndexCache()
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	uri := "file:///api.yaml"
	cache.Set(uri, idx)

	r := NewCrossFileResolver(cache, NewFileGraph())
	if !r.CanResolve(uri, "#/components/schemas/Pet") {
		t.Error("should be able to resolve Pet")
	}
	if r.CanResolve(uri, "#/components/schemas/NonExistent") {
		t.Error("should not resolve NonExistent")
	}
}
