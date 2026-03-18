package navigator

import "testing"

func TestParse_MinimalYAML(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := mustParse(t, content)

	if !idx.IsOpenAPI() {
		t.Fatal("should be detected as OpenAPI root")
	}
	if idx.Document.Version != "3.0.3" {
		t.Errorf("version = %q, want %q", idx.Document.Version, "3.0.3")
	}
	if idx.Document.Info == nil || idx.Document.Info.Title != "Minimal API" {
		t.Error("info.title not parsed")
	}

	if len(idx.Operations) != 3 {
		t.Errorf("expected 3 operations, got %d", len(idx.Operations))
	}
	if _, ok := idx.Operations["listPets"]; !ok {
		t.Error("missing listPets operation")
	}
	if _, ok := idx.Operations["createPet"]; !ok {
		t.Error("missing createPet operation")
	}
	if _, ok := idx.Operations["getPet"]; !ok {
		t.Error("missing getPet operation")
	}

	if len(idx.Schemas) != 2 {
		t.Errorf("expected 2 schemas, got %d", len(idx.Schemas))
	}
	if _, ok := idx.Schemas["Pet"]; !ok {
		t.Error("missing Pet schema")
	}
	if _, ok := idx.Schemas["Error"]; !ok {
		t.Error("missing Error schema")
	}

	if len(idx.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(idx.Parameters))
	}

	if len(idx.Responses) != 1 {
		t.Errorf("expected 1 response, got %d", len(idx.Responses))
	}

	if len(idx.SecuritySchemes) != 1 {
		t.Errorf("expected 1 security scheme, got %d", len(idx.SecuritySchemes))
	}

	if len(idx.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(idx.Tags))
	}

	if len(idx.AllRefs) == 0 {
		t.Error("expected refs to be collected")
	}

	if len(idx.Document.Security) != 1 {
		t.Errorf("expected 1 security requirement, got %d", len(idx.Document.Security))
	}
}

func TestParse_MinimalJSON(t *testing.T) {
	content := loadFixture(t, "single/minimal.json")
	idx := mustParse(t, content)

	if !idx.IsOpenAPI() {
		t.Fatal("should detect JSON OpenAPI root")
	}
	if idx.Document.Version != "3.0.3" {
		t.Errorf("version = %q, want %q", idx.Document.Version, "3.0.3")
	}
	if _, ok := idx.Operations["listItems"]; !ok {
		t.Error("missing listItems operation from JSON spec")
	}
	if _, ok := idx.Schemas["Item"]; !ok {
		t.Error("missing Item schema from JSON spec")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	idx := Parse([]byte(""))
	if idx != nil {
		t.Error("empty input should return nil")
	}
}

func TestParse_MalformedYAML(t *testing.T) {
	idx := Parse([]byte("{{invalid yaml"))
	if idx != nil {
		t.Error("malformed YAML should return nil")
	}
}

func TestParse_SchemaFragment(t *testing.T) {
	content := loadFixture(t, "fragments/schema.yaml")
	idx := mustParse(t, content)
	if idx.IsOpenAPI() {
		t.Error("schema fragment should not be detected as OpenAPI root")
	}
}

func TestParse_RefsCollected(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := mustParse(t, content)

	found := false
	for _, ref := range idx.AllRefs {
		if ref.Target == "#/components/schemas/Pet" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected $ref to #/components/schemas/Pet")
	}
}

func TestParse_OperationsByPath(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := mustParse(t, content)

	ops, ok := idx.OperationsByPath["/pets"]
	if !ok {
		t.Fatal("missing /pets in OperationsByPath")
	}
	if len(ops) != 2 {
		t.Errorf("expected 2 operations on /pets, got %d", len(ops))
	}
}

func TestParse_SchemaProperties(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := mustParse(t, content)

	pet := idx.Schemas["Pet"]
	if pet == nil {
		t.Fatal("Pet schema not found")
	}
	if len(pet.Properties) != 3 {
		t.Errorf("expected 3 properties on Pet, got %d", len(pet.Properties))
	}
	if len(pet.Required) != 1 || pet.Required[0] != "name" {
		t.Errorf("expected required=[name], got %v", pet.Required)
	}
}

func TestParseURI_SetsURI(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := ParseURI("file:///test.yaml", content)
	if idx == nil {
		t.Fatal("ParseURI returned nil")
	}
	if idx.URI() != "file:///test.yaml" {
		t.Errorf("URI = %q, want %q", idx.URI(), "file:///test.yaml")
	}
	for _, ref := range idx.AllRefs {
		if ref.URI != "file:///test.yaml" {
			t.Errorf("ref URI = %q, want %q", ref.URI, "file:///test.yaml")
		}
	}
}
