package navigator

import (
	"reflect"
	"sort"
	"testing"
)

type parseParitySnapshot struct {
	Kind         DocumentKind
	Version      Version
	DocType      DocType
	OperationIDs []string
	SchemaNames  []string
	RefTargets   []string
	IssueCodes   []string
}

func TestParseParity_AcrossFixtures(t *testing.T) {
	tests := []struct {
		name string
		path string
		uri  string
	}{
		{name: "minimal yaml", path: "single/minimal.yaml", uri: "file:///minimal.yaml"},
		{name: "minimal json", path: "single/minimal.json", uri: "file:///minimal.json"},
		{name: "arazzo yaml", path: "arazzo/simple.yaml", uri: "file:///arazzo.yaml"},
		{name: "kitchen sink", path: "golden/kitchen-sink/specs/api.yaml", uri: "file:///kitchen-sink.yaml"},
		{name: "schema fragment", path: "fragments/schema.yaml", uri: "file:///schema.yaml"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			content := loadFixture(t, tc.path)
			standalone := ParseURI(tc.uri, content)
			treesitter := ParseContent(content, tc.uri)
			if standalone == nil || treesitter == nil {
				t.Fatalf("both parsers must succeed: standalone=%v treesitter=%v", standalone != nil, treesitter != nil)
			}

			gotStandalone := snapshotForParity(standalone)
			gotTreeSitter := snapshotForParity(treesitter)
			if !reflect.DeepEqual(gotStandalone, gotTreeSitter) {
				t.Fatalf("parity mismatch\nstandalone=%+v\ntreesitter=%+v", gotStandalone, gotTreeSitter)
			}
		})
	}
}

func TestParseContent_RangeFidelityForOperationsAndRefs(t *testing.T) {
	content := loadFixture(t, "golden/kitchen-sink/specs/api.yaml")
	idx := ParseContent(content, "file:///kitchen-sink.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}
	if idx.Tree() == nil {
		t.Fatal("tree-sitter parse should retain the CST")
	}
	if idx.SemanticRoot() == nil {
		t.Fatal("tree-sitter parse should retain the semantic root")
	}
	if idx.Document == nil || idx.Document.Info == nil {
		t.Fatal("expected parsed document info")
	}
	if IsEmpty(idx.Document.Loc.Range) {
		t.Fatal("document range should be populated")
	}
	if IsEmpty(idx.Document.Info.Loc.Range) {
		t.Fatal("info range should be populated")
	}
	if len(idx.AllRefs) == 0 {
		t.Fatal("expected at least one $ref in kitchen-sink fixture")
	}
	for _, ref := range idx.AllRefs {
		if IsEmpty(ref.Loc.Range) {
			t.Fatalf("$ref %q should have a non-empty range", ref.Target)
		}
	}
	for _, opRef := range idx.AllOperations() {
		if opRef == nil || opRef.Operation == nil {
			t.Fatal("operation ref should be non-nil")
		}
		if IsEmpty(opRef.Operation.Loc.Range) {
			t.Fatalf("operation %s %s should have a non-empty range", opRef.Method, opRef.Path)
		}
	}
}

func snapshotForParity(idx *Index) parseParitySnapshot {
	s := parseParitySnapshot{
		Kind:    idx.Kind,
		Version: idx.Version,
	}
	if idx != nil && idx.Document != nil {
		s.DocType = idx.Document.DocType
	}
	for id := range idx.Operations {
		s.OperationIDs = append(s.OperationIDs, id)
	}
	for name := range idx.Schemas {
		s.SchemaNames = append(s.SchemaNames, name)
	}
	for _, ref := range idx.AllRefs {
		s.RefTargets = append(s.RefTargets, ref.Target)
	}
	for _, iss := range idx.Issues {
		s.IssueCodes = append(s.IssueCodes, iss.Code)
	}
	sort.Strings(s.OperationIDs)
	sort.Strings(s.SchemaNames)
	sort.Strings(s.RefTargets)
	sort.Strings(s.IssueCodes)
	return s
}
