package navigator

import "testing"

func TestIssues_CarryOpenAPIDocumentKind(t *testing.T) {
	idx := ParseContent([]byte(`openapi: "3.0.3"
info:
  title: Test
  version: "1.0.0"
paths:
  /pets:
    get: {}
`), "file:///openapi.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}
	if len(idx.Issues) == 0 {
		t.Fatal("expected issues for missing responses")
	}
	for _, issue := range idx.Issues {
		if issue.DocumentKind != DocumentKindOpenAPI {
			t.Fatalf("issue %q kind = %s, want openapi", issue.Code, issue.DocumentKind)
		}
	}
}

func TestIssues_CarryArazzoDocumentKind(t *testing.T) {
	idx := ParseContent(loadFixture(t, "arazzo/invalid_missing_workflows.yaml"), "file:///arazzo/invalid_missing_workflows.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}
	if len(idx.Issues) == 0 {
		t.Fatal("expected issues for invalid arazzo fixture")
	}
	for _, issue := range idx.Issues {
		if issue.DocumentKind != DocumentKindArazzo {
			t.Fatalf("issue %q kind = %s, want arazzo", issue.Code, issue.DocumentKind)
		}
	}
}
