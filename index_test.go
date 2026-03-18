package navigator

import (
	"sort"
	"testing"
)

func TestIndex_IsOpenAPI(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	if !idx.IsOpenAPI() {
		t.Error("should be OpenAPI")
	}

	fragIdx := mustParse(t, loadFixture(t, "fragments/schema.yaml"))
	if fragIdx.IsOpenAPI() {
		t.Error("fragment should not be OpenAPI")
	}
}

func TestIndex_HasPath(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	if !idx.HasPath("/pets") {
		t.Error("should have /pets")
	}
	if idx.HasPath("/missing") {
		t.Error("should not have /missing")
	}
}

func TestIndex_AllOperations(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	ops := idx.AllOperations()
	if len(ops) != 3 {
		t.Errorf("expected 3 operations, got %d", len(ops))
	}
}

func TestIndex_SchemaNames(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	names := idx.SchemaNames()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "Error" || names[1] != "Pet" {
		t.Errorf("SchemaNames = %v", names)
	}
}

func TestIndex_ComponentNames(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))

	schemas := idx.ComponentNames("schemas")
	sort.Strings(schemas)
	if len(schemas) != 2 {
		t.Errorf("expected 2 schemas, got %v", schemas)
	}

	params := idx.ComponentNames("parameters")
	if len(params) != 1 {
		t.Errorf("expected 1 parameter, got %v", params)
	}

	responses := idx.ComponentNames("responses")
	if len(responses) != 1 {
		t.Errorf("expected 1 response, got %v", responses)
	}

	ss := idx.ComponentNames("securitySchemes")
	if len(ss) != 1 {
		t.Errorf("expected 1 securityScheme, got %v", ss)
	}
}

func TestIndex_RefsTo(t *testing.T) {
	idx := mustParse(t, loadFixture(t, "single/minimal.yaml"))
	refs := idx.RefsTo("#/components/schemas/Pet")
	if len(refs) == 0 {
		t.Error("expected refs pointing to Pet schema")
	}
}

func TestIndex_NilSafety(t *testing.T) {
	var idx *Index
	if idx.IsOpenAPI() {
		t.Error("nil should not be OpenAPI")
	}
	if idx.HasPath("/x") {
		t.Error("nil should not have path")
	}
	if idx.AllOperations() != nil {
		t.Error("nil should return nil operations")
	}
	if idx.URI() != "" {
		t.Error("nil should return empty URI")
	}
	if idx.Content() != nil {
		t.Error("nil should return nil content")
	}
	if idx.Tree() != nil {
		t.Error("nil should return nil tree")
	}
	if idx.SemanticRoot() != nil {
		t.Error("nil should return nil semantic root")
	}
}
