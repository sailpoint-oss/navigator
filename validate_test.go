package navigator

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestValidation_MinimalYAMLNoIssues(t *testing.T) {
	idx := ParseContent(loadFixture(t, "single/minimal.yaml"), "file:///api.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	assertNoErrorSeverity(t, idx.Issues)
}

func TestValidation_MinimalJSONNoIssues(t *testing.T) {
	idx := ParseContent(loadFixture(t, "single/minimal.json"), "file:///api.json")
	if idx == nil {
		t.Fatal("nil index")
	}
	assertNoErrorSeverity(t, idx.Issues)
}

func TestValidation_StandaloneMatchesTreeSitterIssueCount(t *testing.T) {
	content := loadFixture(t, "single/minimal.yaml")
	a := ParseURI("file:///a.yaml", content)
	b := ParseContent(content, "file:///b.yaml")
	if a == nil || b == nil {
		t.Fatal("parse failed")
	}
	if len(a.Issues) != len(b.Issues) {
		t.Fatalf("issue count standalone=%d treesitter=%d\nstandalone: %s\ntreesitter: %s",
			len(a.Issues), len(b.Issues), formatIssues(a.Issues), formatIssues(b.Issues))
	}
}

func TestValidation_FragmentNoStructuralRootRules(t *testing.T) {
	idx := ParseContent(loadFixture(t, "fragments/schema.yaml"), "file:///schema.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	for _, iss := range idx.Issues {
		if iss.Category == CategoryStructural && strings.HasPrefix(iss.Code, "structural.missing") {
			t.Errorf("fragment should not trigger %s: %s", iss.Code, iss.Message)
		}
	}
}

func TestValidation_MissingResponses(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths:
  /a:
    get:
      operationId: x
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	if !issueCodesContain(idx.Issues, "structural.missing-responses") {
		t.Fatalf("expected missing-responses, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_MissingInfoTitle(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  version: "1"
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "structural.missing-info-title") {
		t.Fatalf("expected missing-info-title, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_MissingInfoVersion(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "structural.missing-info-version") {
		t.Fatalf("expected missing-info-version, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_PathsWrongType(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: []
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "structural.ir-wrong-type") {
		t.Fatalf("expected ir-wrong-type for paths, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_OpenAPIWrongType(t *testing.T) {
	yaml := `openapi:
  x: y
info:
  title: T
  version: "1"
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "structural.ir-wrong-type") {
		t.Fatalf("expected ir-wrong-type for openapi, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_UnsupportedVersionWarning(t *testing.T) {
	yaml := `openapi: "9.0.0"
info:
  title: T
  version: "1"
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "structural.unsupported-version") {
		t.Fatalf("expected unsupported-version, got: %s", formatIssues(idx.Issues))
	}
	var sawWarn bool
	for _, iss := range idx.Issues {
		if iss.Code == "structural.unsupported-version" && iss.Severity == SeverityWarning {
			sawWarn = true
		}
	}
	if !sawWarn {
		t.Fatal("unsupported-version should be warning severity")
	}
}

func TestValidation_YAMLDuplicateKey(t *testing.T) {
	yaml := `openapi: "3.0.3"
openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
`
	a := ParseURI("", []byte(yaml))
	if a == nil {
		t.Fatal("nil index")
	}
	if !issueCodesContain(a.Issues, "syntax.duplicate-key") {
		t.Fatalf("expected duplicate-key for standalone, got: %s", formatIssues(a.Issues))
	}

	b := ParseContent([]byte(yaml), "file:///x.yaml")
	if b == nil {
		t.Fatal("nil index")
	}
	if !issueCodesContain(b.Issues, "syntax.duplicate-key") {
		t.Fatalf("expected duplicate-key for treesitter, got: %s", formatIssues(b.Issues))
	}
}

func TestValidation_ParameterMissingName(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths:
  /a:
    get:
      parameters:
        - in: query
          schema:
            type: string
      responses:
        "200":
          description: ok
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "schema.parameter-missing-name") {
		t.Fatalf("expected parameter-missing-name, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_SchemaUnknownType(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
components:
  schemas:
    X:
      type: notatype
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "schema.unknown-type") {
		t.Fatalf("expected schema.unknown-type, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_SecuritySchemeHTTPMissingScheme(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
components:
  securitySchemes:
    A:
      type: http
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "schema.security-scheme-http-fields") {
		t.Fatalf("expected security-scheme-http-fields, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_MaxIssues(t *testing.T) {
	// duplicate keys + structural issues — cap at 1
	yaml := `openapi: "3.0.3"
openapi: "3.0.3"
info:
  version: "1"
paths: {}
`
	idx := ParseURI("file:///x.yaml", []byte(yaml))
	if idx == nil {
		t.Fatal("nil")
	}
	opts := ValidationOptions{MaxIssues: 1}
	issues := idx.Revalidate(opts)
	if len(issues) != 1 {
		t.Fatalf("want 1 issue, got %d: %s", len(issues), formatIssues(issues))
	}
}

func TestValidation_ContextCancels(t *testing.T) {
	yaml := `openapi: "3.0.3"
openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
`
	idx := ParseURI("file:///x.yaml", []byte(yaml))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	opts := ValidationOptions{Context: ctx}
	issues := idx.Revalidate(opts)
	if len(issues) != 0 {
		t.Fatalf("canceled context should yield no issues, got: %s", formatIssues(issues))
	}
}

func TestValidation_NewIndexFromDocument(t *testing.T) {
	doc := &Document{
		DocType:       DocTypeRoot,
		Version:       "3.0.3",
		ParsedVersion: Version30,
		Info: &Info{
			Title:   "T",
			Version: "1",
			Loc:     Loc{Range: FileStartRange},
		},
		Paths: map[string]*PathItem{},
		Loc:   Loc{Range: FileStartRange},
	}
	idx := NewIndexFromDocument(doc)
	if idx == nil {
		t.Fatal("nil")
	}
	// paths empty object — no operations; valid
	assertNoErrorSeverity(t, idx.Issues)
}

func TestValidation_NewIndexFromDocument_MissingTitle(t *testing.T) {
	doc := &Document{
		DocType:       DocTypeRoot,
		Version:       "3.0.3",
		ParsedVersion: Version30,
		Info: &Info{
			Version: "1",
			Loc:     Loc{Range: FileStartRange},
		},
		Paths: map[string]*PathItem{},
		Loc:   Loc{Range: FileStartRange},
	}
	idx := NewIndexFromDocument(doc)
	if !issueCodesContain(idx.Issues, "structural.missing-info-title") {
		t.Fatalf("expected missing title: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_SkipFlags(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  version: "1"
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	// baseline: has structural issue (missing title)
	if !issueCodesContain(idx.Issues, "structural.missing-info-title") {
		t.Fatalf("precondition: %s", formatIssues(idx.Issues))
	}
	issues := idx.Revalidate(ValidationOptions{SkipStructural: true})
	for _, iss := range issues {
		if iss.Category == CategoryStructural {
			t.Errorf("unexpected structural issue with SkipStructural: %s", iss.Code)
		}
	}
}

func TestValidation_ServersWrongType(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
servers: {}
paths: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "structural.ir-wrong-type") {
		t.Fatalf("expected ir-wrong-type for servers, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_TagsWrongType(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
tags: {}
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "structural.ir-wrong-type") {
		t.Fatalf("expected ir-wrong-type for tags, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_SecuritySchemeAPIKeyMissingName(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
components:
  securitySchemes:
    A:
      type: apiKey
      in: header
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if !issueCodesContain(idx.Issues, "schema.security-scheme-apikey-fields") {
		t.Fatalf("expected apiKey fields issue, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_OAuth2MissingFlowsStandalone(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
components:
  securitySchemes:
    A:
      type: oauth2
`
	idx := ParseURI("file:///x.yaml", []byte(yaml))
	if idx == nil {
		t.Fatal("nil index")
	}
	if !issueCodesContain(idx.Issues, "schema.security-scheme-oauth2-fields") {
		t.Fatalf("expected oauth2 flows issue, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_JSONCSTParseError(t *testing.T) {
	// Invalid JSON: missing comma between properties
	bad := []byte(`{"openapi":"3.0.3" "info":{"title":"T","version":"1"},"paths":{}}`)
	idx := ParseContent(bad, "file:///x.json")
	if idx == nil {
		t.Fatal("expected index even on partial parse")
	}
	if !issueCodesContain(idx.Issues, "syntax.parse-error") {
		t.Fatalf("expected syntax.parse-error, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_RootNotMapping(t *testing.T) {
	yaml := `- not
- a
- document
`
	idx := ParseContent([]byte(yaml), "file:///x.yaml")
	if idx == nil {
		t.Fatal("nil")
	}
	if !issueCodesContain(idx.Issues, "structural.root-not-mapping") {
		t.Fatalf("expected root-not-mapping: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_RootNotMappingStandalone(t *testing.T) {
	yaml := `- not
- a
- document
`
	idx := ParseURI("file:///x.yaml", []byte(yaml))
	if idx == nil {
		t.Fatal("nil")
	}
	if !issueCodesContain(idx.Issues, "structural.root-not-mapping") {
		t.Fatalf("expected root-not-mapping for standalone: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_OpenAPIWrongTypeStandalone(t *testing.T) {
	yaml := `openapi:
  x: y
info:
  title: T
  version: "1"
paths: {}
`
	idx := ParseURI("file:///x.yaml", []byte(yaml))
	if idx == nil {
		t.Fatal("nil")
	}
	if !issueCodesContain(idx.Issues, "structural.ir-wrong-type") {
		t.Fatalf("expected ir-wrong-type for standalone openapi, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_PathsWrongTypeStandalone(t *testing.T) {
	yaml := `openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: []
`
	idx := ParseURI("file:///x.yaml", []byte(yaml))
	if idx == nil {
		t.Fatal("nil")
	}
	if !issueCodesContain(idx.Issues, "structural.ir-wrong-type") {
		t.Fatalf("expected ir-wrong-type for standalone paths, got: %s", formatIssues(idx.Issues))
	}
}

func TestValidation_LoadValidateFixtures(t *testing.T) {
	tests := []struct {
		file  string
		codes []string
	}{
		{"testdata/validate/ok_minimal.yaml", nil},
		{"testdata/validate/bad_duplicate_key.yaml", []string{"syntax.duplicate-key"}},
		{"testdata/validate/bad_missing_responses.yaml", []string{"structural.missing-responses"}},
		{"testdata/validate/bad_paths_array.yaml", []string{"structural.ir-wrong-type"}},
	}
	parsers := []struct {
		name  string
		parse func([]byte) *Index
	}{
		{"treesitter", func(b []byte) *Index { return ParseContent(b, "file:///fixture.yaml") }},
		{"standalone", func(b []byte) *Index { return ParseURI("file:///fixture.yaml", b) }},
	}
	for _, tc := range tests {
		for _, p := range parsers {
			t.Run(p.name+"/"+tc.file, func(t *testing.T) {
				data, err := os.ReadFile(tc.file)
				if err != nil {
					t.Fatal(err)
				}
				idx := p.parse(data)
				if idx == nil {
					t.Fatal("nil index")
				}
				got := issueCodeSet(idx.Issues)
				for _, want := range tc.codes {
					if !got[want] {
						t.Fatalf("missing code %q in %v\nissues: %s", want, got, formatIssues(idx.Issues))
					}
				}
				if len(tc.codes) == 0 {
					assertNoErrorSeverity(t, idx.Issues)
				}
			})
		}
	}
}

func TestValidation_AllCodesInRegistry(t *testing.T) {
	// Exercises multiple bad documents to collect a wide range of codes,
	// then asserts every emitted code is in knownIssueCodes.
	docs := []string{
		`- not a mapping`,
		`openapi: "3.0.3"
openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}`,
		`openapi:
  x: y
info:
  title: T
  version: "1"
paths: {}`,
		`openapi: "9.0.0"
info:
  version: "1"
paths: []`,
		`openapi: "3.0.3"
info:
  title: T
  version: "1"
paths:
  /a:
    get:
      parameters:
        - in: query
      responses:
        "200":
          description: ok`,
		`openapi: "3.0.3"
info:
  title: T
  version: "1"
paths: {}
components:
  schemas:
    X:
      type: notatype
  securitySchemes:
    A:
      type: apiKey
      in: header
    B:
      type: http
    C:
      type: oauth2
    D:
      type: openIdConnect
    E:
      type: alien`,
	}
	for _, raw := range docs {
		for _, idx := range []*Index{
			ParseContent([]byte(raw), "file:///x.yaml"),
			ParseURI("file:///x.yaml", []byte(raw)),
		} {
			if idx == nil {
				continue
			}
			issues := idx.Revalidate(ValidationOptions{SkipMetaSchema: true})
			for _, iss := range issues {
				if _, ok := knownIssueCodes[iss.Code]; !ok {
					t.Errorf("issue code %q not in knownIssueCodes registry", iss.Code)
				}
			}
		}
	}
}

// --- helpers ---

func assertNoErrorSeverity(t *testing.T, issues []Issue) {
	t.Helper()
	for _, iss := range issues {
		if iss.Severity == SeverityError {
			t.Errorf("unexpected error issue %s: %s", iss.Code, iss.Message)
		}
	}
}

func issueCodesContain(issues []Issue, code string) bool {
	for _, iss := range issues {
		if iss.Code == code {
			return true
		}
	}
	return false
}

func issueCodeSet(issues []Issue) map[string]bool {
	m := make(map[string]bool)
	for _, iss := range issues {
		m[iss.Code] = true
	}
	return m
}

func formatIssues(issues []Issue) string {
	if len(issues) == 0 {
		return "(none)"
	}
	var b strings.Builder
	for _, iss := range issues {
		b.WriteString(iss.Code)
		b.WriteString(": ")
		b.WriteString(iss.Message)
		b.WriteString("; ")
	}
	return strings.TrimSuffix(b.String(), "; ")
}
