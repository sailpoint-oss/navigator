package navigator

import (
	"reflect"
	"sort"
	"testing"
)

type arazzoSnapshot struct {
	Kind        DocumentKind
	WorkflowIDs []string
	SourceNames []string
	StepIDs     []string
	IssueCodes  []string
}

func TestParse_ArazzoYAMLParity(t *testing.T) {
	content := loadFixture(t, "arazzo/simple.yaml")
	standalone := ParseURI("file:///arazzo/simple.yaml", content)
	treesitter := ParseContent(content, "file:///arazzo/simple.yaml")
	if standalone == nil || treesitter == nil {
		t.Fatalf("both parsers must succeed: standalone=%v treesitter=%v", standalone != nil, treesitter != nil)
	}

	if !standalone.IsArazzo() {
		t.Fatalf("standalone kind = %s, want arazzo", standalone.Kind)
	}
	if !treesitter.IsArazzo() {
		t.Fatalf("tree-sitter kind = %s, want arazzo", treesitter.Kind)
	}
	if standalone.PrimaryValue() != standalone.Arazzo {
		t.Fatal("standalone PrimaryValue should return the typed Arazzo document")
	}
	if treesitter.PrimaryValue() != treesitter.Arazzo {
		t.Fatal("tree-sitter PrimaryValue should return the typed Arazzo document")
	}

	if standalone.Arazzo == nil || treesitter.Arazzo == nil {
		t.Fatal("typed Arazzo document should be populated")
	}
	if standalone.Arazzo.Info == nil || standalone.Arazzo.Info.Title != "Simple workflow" {
		t.Fatalf("standalone title = %+v, want Simple workflow", standalone.Arazzo.Info)
	}
	if treesitter.Arazzo.Info == nil || treesitter.Arazzo.Info.Title != "Simple workflow" {
		t.Fatalf("tree-sitter title = %+v, want Simple workflow", treesitter.Arazzo.Info)
	}

	if !reflect.DeepEqual(snapshotForArazzo(standalone), snapshotForArazzo(treesitter)) {
		t.Fatalf("Arazzo parity mismatch\nstandalone=%+v\ntreesitter=%+v", snapshotForArazzo(standalone), snapshotForArazzo(treesitter))
	}
	assertNoErrorSeverity(t, standalone.Issues)
	assertNoErrorSeverity(t, treesitter.Issues)
}

func TestParse_ArazzoJSON(t *testing.T) {
	content := loadFixture(t, "arazzo/simple.json")
	idx := ParseContent(content, "file:///arazzo/simple.json")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}
	if !idx.IsArazzo() {
		t.Fatalf("kind = %s, want arazzo", idx.Kind)
	}
	if idx.Arazzo == nil || len(idx.Arazzo.Workflows) != 1 {
		t.Fatalf("workflows = %+v, want one workflow", idx.Arazzo)
	}
	if got := idx.Arazzo.Workflows[0].Steps[0].OperationID; got != "listPets" {
		t.Fatalf("operationId = %q, want listPets", got)
	}
	assertNoErrorSeverity(t, idx.Issues)
}

func TestValidate_ArazzoStructuralIssues(t *testing.T) {
	idx := ParseContent(loadFixture(t, "arazzo/invalid_duplicate_ids.yaml"), "file:///arazzo/invalid_duplicate_ids.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	if !idx.IsArazzo() {
		t.Fatalf("kind = %s, want arazzo", idx.Kind)
	}

	wantCodes := []string{
		"structural.duplicate-workflow-id",
		"structural.duplicate-step-id",
		"structural.missing-step-id",
		"structural.step-target-count",
	}
	for _, code := range wantCodes {
		if !issueCodesContain(idx.Issues, code) {
			t.Fatalf("expected %s in issues, got: %s", code, formatIssues(idx.Issues))
		}
	}
}

func TestValidate_ArazzoMetaSchema(t *testing.T) {
	idx := ParseContent(loadFixture(t, "arazzo/invalid_missing_workflows.yaml"), "file:///arazzo/invalid_missing_workflows.yaml")
	if idx == nil {
		t.Fatal("nil index")
	}
	if !idx.IsArazzo() {
		t.Fatalf("kind = %s, want arazzo", idx.Kind)
	}
	if !issueCodesContain(idx.Issues, "meta.required") {
		t.Fatalf("expected meta.required, got: %s", formatIssues(idx.Issues))
	}
	if !issueCodesContain(idx.Issues, "structural.missing-workflows") {
		t.Fatalf("expected structural.missing-workflows, got: %s", formatIssues(idx.Issues))
	}
}

func snapshotForArazzo(idx *Index) arazzoSnapshot {
	s := arazzoSnapshot{}
	if idx == nil {
		return s
	}
	s.Kind = idx.Kind
	if idx.Arazzo != nil {
		for _, source := range idx.Arazzo.SourceDescriptions {
			s.SourceNames = append(s.SourceNames, source.Name)
		}
		for _, workflow := range idx.Arazzo.Workflows {
			s.WorkflowIDs = append(s.WorkflowIDs, workflow.WorkflowID)
			for _, step := range workflow.Steps {
				s.StepIDs = append(s.StepIDs, step.StepID)
			}
		}
	}
	for _, issue := range idx.Issues {
		s.IssueCodes = append(s.IssueCodes, issue.Code)
	}
	sort.Strings(s.SourceNames)
	sort.Strings(s.WorkflowIDs)
	sort.Strings(s.StepIDs)
	sort.Strings(s.IssueCodes)
	return s
}
