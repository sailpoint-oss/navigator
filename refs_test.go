package navigator

import "testing"

func TestExtractExternalRefs(t *testing.T) {
	content := []byte(`openapi: "3.0.3"
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
                $ref: "./schemas/Pet.yaml"
        "404":
          $ref: "#/components/responses/NotFound"
components:
  responses:
    NotFound:
      description: Not found`)

	uri := "file:///workspace/api.yaml"
	idx := ParseURI(uri, content)
	if idx == nil {
		t.Fatal("parse failed")
	}

	edges := ExtractExternalRefs(uri, idx)
	if len(edges) != 1 {
		t.Fatalf("expected 1 external ref edge, got %d", len(edges))
	}
	if edges[0].FromURI != uri {
		t.Errorf("from = %q", edges[0].FromURI)
	}
}

func TestCollectExternalRefTargets(t *testing.T) {
	content := []byte(`openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
paths:
  /a:
    $ref: "./paths/a.yaml"
  /b:
    $ref: "./paths/b.yaml"
  /c:
    $ref: "./paths/a.yaml"`)

	uri := "file:///workspace/api.yaml"
	idx := ParseURI(uri, content)
	targets := CollectExternalRefTargets(uri, idx)
	if len(targets) != 2 {
		t.Errorf("expected 2 unique targets, got %d: %v", len(targets), targets)
	}
}

func TestUpdateGraphFromIndex(t *testing.T) {
	content := []byte(`openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
paths:
  /pets:
    $ref: "./pets.yaml"`)

	uri := "file:///api.yaml"
	idx := ParseURI(uri, content)
	g := NewFileGraph()

	UpdateGraphFromIndex(g, uri, idx)

	deps := g.DependenciesOf(uri)
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(deps))
	}
}
