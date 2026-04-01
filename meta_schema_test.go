package navigator

import (
	"os"
	"strings"
	"testing"
)

func TestMetaSchema_RootMissingInfo(t *testing.T) {
	yaml := `openapi: "3.0.3"
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	if !issueCodesContain(idx.Issues, "meta.required") {
		t.Fatalf("expected meta.required, got: %s", formatIssues(idx.Issues))
	}
	var sawInfo bool
	for _, iss := range idx.Issues {
		if iss.Code == "meta.required" && strings.Contains(iss.Message, "info") {
			sawInfo = true
		}
	}
	if !sawInfo {
		t.Fatalf("expected curated message mentioning info: %s", formatIssues(idx.Issues))
	}
}

func TestMetaSchema_MinimalRootPasses(t *testing.T) {
	data, err := os.ReadFile("testdata/validate/ok_minimal.yaml")
	if err != nil {
		t.Fatal(err)
	}
	idx := ParseContent(data, "file:///ok.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	for _, iss := range idx.Issues {
		if iss.Category == CategoryMeta && iss.Severity == SeverityError {
			t.Errorf("unexpected meta error %s: %s", iss.Code, iss.Message)
		}
	}
}

func TestMetaSchema_FragmentSchemaObject(t *testing.T) {
	yaml := `type: object
properties:
  id:
    type: string
`
	idx := ParseContent([]byte(yaml), "file:///frag.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	for _, iss := range idx.Issues {
		if iss.Category == CategoryMeta && iss.Severity == SeverityError {
			t.Errorf("unexpected meta error on fragment: %s: %s", iss.Code, iss.Message)
		}
	}
}

func TestMetaSchema_SkipMetaSchema(t *testing.T) {
	yaml := `openapi: "3.0.3"
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	issues := idx.Revalidate(ValidationOptions{SkipMetaSchema: true})
	for _, iss := range issues {
		if iss.Category == CategoryMeta {
			t.Fatalf("unexpected meta issue with SkipMetaSchema: %s", iss.Code)
		}
	}
}

func TestMetaSchema_CompiledSchemasLoad(t *testing.T) {
	_, _, err := compiledMetaSchemas()
	if err != nil {
		t.Fatalf("compiledMetaSchemas: %v", err)
	}
}

func TestMetaSchema_CompiledArazzoSchemaLoad(t *testing.T) {
	_, err := compiledArazzoMetaSchema()
	if err != nil {
		t.Fatalf("compiledArazzoMetaSchema: %v", err)
	}
}

func TestMetaSchema_CompiledVersionedSchemasLoad(t *testing.T) {
	roots, frags, err := compiledVersionedMetaSchemas()
	if err != nil {
		t.Fatalf("compiledVersionedMetaSchemas: %v", err)
	}
	for _, ver := range []Version{Version20, Version30, Version31, Version32} {
		if roots[ver] == nil {
			t.Fatalf("missing versioned root schema for %s", ver)
		}
	}
	if frags[Version20][FragmentSchema] == nil {
		t.Fatal("missing Swagger 2.0 schema fragment meta-schema")
	}
	for _, ver := range []Version{Version30, Version31, Version32} {
		if frags[ver][FragmentSchema] == nil {
			t.Fatalf("missing schema fragment meta-schema for %s", ver)
		}
		if frags[ver][FragmentOperation] == nil {
			t.Fatalf("missing operation fragment meta-schema for %s", ver)
		}
		if frags[ver][FragmentComponents] == nil {
			t.Fatalf("missing components fragment meta-schema for %s", ver)
		}
	}
}

func TestMetaSchema_TypeMismatchMessage(t *testing.T) {
	yaml := `openapi: 3
info:
  title: T
  version: "1"
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	// With a version union (3.0 / 3.1 / 3.2 / 2.0), a non-string `openapi`
	// value may surface as meta.type on /openapi or as failed union (meta.any-of / meta.ref).
	var saw bool
	for _, iss := range idx.Issues {
		if iss.Category != CategoryMeta {
			continue
		}
		if iss.Code == "meta.type" && strings.Contains(iss.Message, "string") {
			saw = true
			break
		}
		if iss.Code == "meta.any-of" || iss.Code == "meta.ref" {
			saw = true
			break
		}
	}
	if !saw {
		t.Fatalf("expected meta signal for wrong openapi type (type or union at /openapi), got: %s", formatIssues(idx.Issues))
	}
}

func TestMetaSchema_FixturesFromDisk(t *testing.T) {
	tests := []struct {
		file      string
		wantMeta  string
		wantClean bool
	}{
		{"testdata/meta/invalid_root_no_info.yaml", "meta.required", false},
		{"testdata/meta/ok_fragment_schema.yaml", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			idx := ParseContent(data, "file:///fixture.yaml")
			if idx == nil {
				t.Fatal("nil index")
			}
			if tc.wantClean {
				for _, iss := range idx.Issues {
					if iss.Category == CategoryMeta && iss.Severity == SeverityError {
						t.Errorf("unexpected meta %s: %s", iss.Code, iss.Message)
					}
				}
				return
			}
			if !issueCodesContain(idx.Issues, tc.wantMeta) {
				t.Fatalf("want code %q, got %s", tc.wantMeta, formatIssues(idx.Issues))
			}
		})
	}
}

func TestMetaSchema_ArazzoFixturesFromDisk(t *testing.T) {
	tests := []struct {
		file      string
		wantMeta  string
		wantClean bool
	}{
		{"testdata/arazzo/simple.yaml", "", true},
		{"testdata/arazzo/invalid_missing_workflows.yaml", "meta.required", false},
	}
	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			idx := ParseContent(data, "file:///fixture.yaml")
			if idx == nil {
				t.Fatal("nil index")
			}
			if tc.wantClean {
				for _, iss := range idx.Issues {
					if iss.Category == CategoryMeta && iss.Severity == SeverityError {
						t.Errorf("unexpected meta %s: %s", iss.Code, iss.Message)
					}
				}
				return
			}
			if !issueCodesContain(idx.Issues, tc.wantMeta) {
				t.Fatalf("want code %q, got %s", tc.wantMeta, formatIssues(idx.Issues))
			}
		})
	}
}

func TestMetaSchema_AllMetaCodesInRegistry(t *testing.T) {
	// Documents chosen to exercise meta validation paths.
	docs := []string{
		`openapi: "3.0.3"
paths: {}`,
		`openapi: 1
info:
  title: T
  version: "1"
paths: {}`,
		`arazzo: "1.0.1"
info:
  title: T
  version: "1"
sourceDescriptions:
  - name: api
    url: ./petstore.yaml`,
	}
	for _, raw := range docs {
		idx := ParseContent([]byte(raw), "file:///x.yaml")
		if idx == nil {
			continue
		}
		for _, iss := range idx.Issues {
			if iss.Category != CategoryMeta {
				continue
			}
			if _, ok := knownIssueCodes[iss.Code]; !ok {
				t.Errorf("meta code %q not in knownIssueCodes", iss.Code)
			}
		}
	}
}

func TestMetaSchema_ContextAwareMessages(t *testing.T) {
	yaml := `openapi: "3.1.0"
info:
  title: Test
  version: "1.0"
components:
  schemas:
    Pet:
      yo: yo
      type: object
paths:
  /test:
    get:
      summary: t
      foo: bar
      responses:
        "200":
          description: ok
`
	idx := ParseContent([]byte(yaml), "file:///test.yaml")
	t.Logf("Total issues: %d", len(idx.Issues))
	for _, iss := range idx.Issues {
		t.Logf("  [%s] ptr=%s msg=%s", iss.Code, iss.Pointer, iss.Message)
	}

	// Verify messages use contextual OpenAPI terminology
	var foundSchemaCtx, foundOperationCtx bool
	for _, iss := range idx.Issues {
		if strings.Contains(iss.Message, "Schema") || strings.Contains(iss.Message, "Reference") {
			foundSchemaCtx = true
		}
		if strings.Contains(iss.Message, "Operation") {
			foundOperationCtx = true
		}
	}
	if !foundSchemaCtx {
		t.Error("expected at least one message mentioning 'Schema' or 'Reference'")
	}
	if !foundOperationCtx {
		t.Error("expected at least one message mentioning 'Operation'")
	}
}
