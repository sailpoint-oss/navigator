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

func TestParseContent_ContactLicenseParity(t *testing.T) {
	content := loadFixture(t, "golden/kitchen-sink/specs/api.yaml")
	standalone := Parse(content)
	treesitter := ParseContent(content, "file:///api.yaml")

	if standalone == nil || treesitter == nil {
		t.Fatal("both parsers should succeed")
	}

	if standalone.Document.Info.Contact == nil {
		t.Fatal("standalone should parse contact")
	}
	if treesitter.Document.Info.Contact == nil {
		t.Fatal("treesitter should parse contact")
	}
	if standalone.Document.Info.Contact.Name != treesitter.Document.Info.Contact.Name {
		t.Errorf("contact name mismatch: standalone=%q treesitter=%q",
			standalone.Document.Info.Contact.Name, treesitter.Document.Info.Contact.Name)
	}
	if standalone.Document.Info.Contact.Email != treesitter.Document.Info.Contact.Email {
		t.Errorf("contact email mismatch: standalone=%q treesitter=%q",
			standalone.Document.Info.Contact.Email, treesitter.Document.Info.Contact.Email)
	}

	if standalone.Document.Info.License == nil {
		t.Fatal("standalone should parse license")
	}
	if treesitter.Document.Info.License == nil {
		t.Fatal("treesitter should parse license")
	}
	if standalone.Document.Info.License.Name != treesitter.Document.Info.License.Name {
		t.Errorf("license name mismatch: standalone=%q treesitter=%q",
			standalone.Document.Info.License.Name, treesitter.Document.Info.License.Name)
	}
}

func TestParseContent_ContactFields(t *testing.T) {
	content := []byte(`openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
  contact:
    name: Acme Support
    url: https://acme.com/support
    email: support@acme.com
  license:
    name: Apache-2.0
    url: https://www.apache.org/licenses/LICENSE-2.0
paths: {}
`)
	idx := ParseContent(content, "file:///api.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}
	info := idx.Document.Info
	if info.Contact == nil {
		t.Fatal("Contact should not be nil")
	}
	if info.Contact.Name != "Acme Support" {
		t.Errorf("contact name = %q, want %q", info.Contact.Name, "Acme Support")
	}
	if info.Contact.URL != "https://acme.com/support" {
		t.Errorf("contact url = %q, want %q", info.Contact.URL, "https://acme.com/support")
	}
	if info.Contact.Email != "support@acme.com" {
		t.Errorf("contact email = %q, want %q", info.Contact.Email, "support@acme.com")
	}
	if info.License == nil {
		t.Fatal("License should not be nil")
	}
	if info.License.Name != "Apache-2.0" {
		t.Errorf("license name = %q, want %q", info.License.Name, "Apache-2.0")
	}
	if info.License.URL != "https://www.apache.org/licenses/LICENSE-2.0" {
		t.Errorf("license url = %q, want %q", info.License.URL, "https://www.apache.org/licenses/LICENSE-2.0")
	}
}

func TestParseContent_SchemaExample(t *testing.T) {
	content := []byte(`openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
paths: {}
components:
  schemas:
    Status:
      type: string
      example: active
      enum:
        - active
        - inactive
`)
	idx := ParseContent(content, "file:///api.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}
	schema := idx.Document.Components.Schemas["Status"]
	if schema == nil {
		t.Fatal("Status schema not found")
	}
	if schema.Example == nil {
		t.Fatal("Schema.Example should not be nil")
	}
	if schema.Example.Value != "active" {
		t.Errorf("schema example = %q, want %q", schema.Example.Value, "active")
	}
}

func TestParseContent_MediaTypeExample(t *testing.T) {
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
                type: string
              example: hello
    post:
      operationId: createPet
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: string
            example: world
      responses:
        "201":
          description: Created
`)
	idx := ParseContent(content, "file:///api.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}

	resp := idx.Document.Paths["/pets"].Get.Responses["200"]
	if resp == nil {
		t.Fatal("200 response not found")
	}
	mt := resp.Content["application/json"]
	if mt == nil {
		t.Fatal("application/json media type not found in response")
	}
	if mt.Example == nil {
		t.Fatal("response MediaType.Example should not be nil")
	}
	if mt.Example.Value != "hello" {
		t.Errorf("response example = %q, want %q", mt.Example.Value, "hello")
	}

	rb := idx.Document.Paths["/pets"].Post.RequestBody
	if rb == nil {
		t.Fatal("requestBody not found")
	}
	rbMt := rb.Content["application/json"]
	if rbMt == nil {
		t.Fatal("application/json media type not found in requestBody")
	}
	if rbMt.Example == nil {
		t.Fatal("requestBody MediaType.Example should not be nil")
	}
	if rbMt.Example.Value != "world" {
		t.Errorf("requestBody example = %q, want %q", rbMt.Example.Value, "world")
	}
}

func TestParseContent_MediaTypeExamples(t *testing.T) {
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
                type: string
              examples:
                good:
                  summary: A good example
                  value: hello
                bad:
                  summary: A bad example
                  value: oops
`)
	idx := ParseContent(content, "file:///api.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}

	mt := idx.Document.Paths["/pets"].Get.Responses["200"].Content["application/json"]
	if mt == nil {
		t.Fatal("media type not found")
	}
	if len(mt.Examples) != 2 {
		t.Fatalf("expected 2 examples, got %d", len(mt.Examples))
	}
	good := mt.Examples["good"]
	if good == nil {
		t.Fatal("'good' example not found")
	}
	if good.Summary != "A good example" {
		t.Errorf("good.Summary = %q, want %q", good.Summary, "A good example")
	}
	if good.Value == nil || good.Value.Value != "hello" {
		t.Errorf("good.Value = %v, want 'hello'", good.Value)
	}
	bad := mt.Examples["bad"]
	if bad == nil {
		t.Fatal("'bad' example not found")
	}
	if bad.Summary != "A bad example" {
		t.Errorf("bad.Summary = %q, want %q", bad.Summary, "A bad example")
	}
}
