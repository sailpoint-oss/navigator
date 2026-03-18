package navigator

import (
	"testing"
)

func TestParseContent_YAML(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := ParseContent(content, "file:///api.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}
	if !idx.IsOpenAPI() {
		t.Error("should be OpenAPI")
	}
	if idx.Tree() == nil {
		t.Error("tree-sitter tree should be non-nil")
	}
	if idx.SemanticRoot() == nil {
		t.Error("semantic root should be non-nil")
	}
	if idx.Content() == nil {
		t.Error("content should be non-nil")
	}
}

func TestParseContent_JSON(t *testing.T) {
	content := loadFixture(t, "single/minimal.json")
	idx := ParseContent(content, "file:///api.json")
	if idx == nil {
		t.Fatal("ParseContent returned nil for JSON")
	}
	if !idx.IsOpenAPI() {
		t.Error("JSON spec should be OpenAPI")
	}
}

func TestParseContent_OperationParity(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	standalone := Parse(content)
	treesitter := ParseContent(content, "file:///api.yaml")

	if standalone == nil || treesitter == nil {
		t.Fatal("both parsers should succeed")
	}

	if len(standalone.Operations) != len(treesitter.Operations) {
		t.Errorf("operation count mismatch: standalone=%d treesitter=%d",
			len(standalone.Operations), len(treesitter.Operations))
	}

	for opID := range standalone.Operations {
		if _, ok := treesitter.Operations[opID]; !ok {
			t.Errorf("treesitter missing operation %q", opID)
		}
	}
}

func TestParseContent_SchemaParity(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	standalone := Parse(content)
	treesitter := ParseContent(content, "file:///api.yaml")

	if len(standalone.Schemas) != len(treesitter.Schemas) {
		t.Errorf("schema count mismatch: standalone=%d treesitter=%d",
			len(standalone.Schemas), len(treesitter.Schemas))
	}

	for name := range standalone.Schemas {
		if _, ok := treesitter.Schemas[name]; !ok {
			t.Errorf("treesitter missing schema %q", name)
		}
	}
}

func TestParseContent_RefParity(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	standalone := Parse(content)
	treesitter := ParseContent(content, "file:///api.yaml")

	if len(standalone.AllRefs) != len(treesitter.AllRefs) {
		t.Errorf("ref count mismatch: standalone=%d treesitter=%d",
			len(standalone.AllRefs), len(treesitter.AllRefs))
	}
}

func TestParseContent_Version(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	idx := ParseContent(content, "file:///api.yaml")
	if idx.Version != Version30 {
		t.Errorf("version = %q, want %q", idx.Version, Version30)
	}
}

func TestParseContent_Fragment(t *testing.T) {
	content := loadFixture(t, "fragments/schema.yaml")
	idx := ParseContent(content, "file:///schema.yaml")
	if idx == nil {
		t.Fatal("should parse fragment")
	}
	if idx.IsOpenAPI() {
		t.Error("fragment should not be OpenAPI root")
	}
}

func TestParseContent_Empty(t *testing.T) {
	idx := ParseContent([]byte(""), "file:///empty.yaml")
	if idx != nil {
		t.Error("empty content should return nil")
	}
}
