package navigator

import "testing"

func TestPathItem_Operations(t *testing.T) {
	pi := &PathItem{
		Get:    &Operation{OperationID: "getOp"},
		Post:   &Operation{OperationID: "postOp"},
		Delete: &Operation{OperationID: "deleteOp"},
	}
	ops := pi.Operations()
	if len(ops) != 3 {
		t.Fatalf("expected 3 operations, got %d", len(ops))
	}
	methods := map[string]string{}
	for _, op := range ops {
		methods[op.Method] = op.Operation.OperationID
	}
	if methods["get"] != "getOp" {
		t.Error("missing get")
	}
	if methods["post"] != "postOp" {
		t.Error("missing post")
	}
	if methods["delete"] != "deleteOp" {
		t.Error("missing delete")
	}
}

func TestPathItem_Operations_Nil(t *testing.T) {
	var pi *PathItem
	if ops := pi.Operations(); ops != nil {
		t.Error("nil PathItem should return nil operations")
	}
}

func TestOperation_TagNames(t *testing.T) {
	op := &Operation{
		Tags: []TagUsage{{Name: "pets"}, {Name: "admin"}},
	}
	names := op.TagNames()
	if len(names) != 2 || names[0] != "pets" || names[1] != "admin" {
		t.Errorf("unexpected tag names: %v", names)
	}
}

func TestOperation_HasTag(t *testing.T) {
	op := &Operation{
		Tags: []TagUsage{{Name: "pets"}, {Name: "admin"}},
	}
	if _, ok := op.HasTag("pets"); !ok {
		t.Error("should have pets tag")
	}
	if _, ok := op.HasTag("missing"); ok {
		t.Error("should not have missing tag")
	}
}

func TestSecurityRequirement_SchemeNames(t *testing.T) {
	sr := SecurityRequirement{
		Entries: []SecurityRequirementEntry{
			{Name: "oauth2", Scopes: []string{"read", "write"}},
			{Name: "apiKey"},
		},
	}
	names := sr.SchemeNames()
	if len(names) != 2 {
		t.Fatalf("expected 2 scheme names, got %d", len(names))
	}
}

func TestSecurityRequirement_HasScheme(t *testing.T) {
	sr := SecurityRequirement{
		Entries: []SecurityRequirementEntry{
			{Name: "oauth2", Scopes: []string{"read"}},
		},
	}
	entry, ok := sr.HasScheme("oauth2")
	if !ok {
		t.Fatal("should have oauth2 scheme")
	}
	if len(entry.Scopes) != 1 || entry.Scopes[0] != "read" {
		t.Errorf("unexpected scopes: %v", entry.Scopes)
	}
	if _, ok := sr.HasScheme("missing"); ok {
		t.Error("should not have missing scheme")
	}
}

func TestLocOrFallback(t *testing.T) {
	zero := Loc{}
	nonZero := Loc{Range: Range{Start: Position{Line: 5, Character: 3}}}

	if got := LocOrFallback(nonZero, zero); got.Range.Start.Line != 5 {
		t.Error("should return primary when non-zero")
	}
	if got := LocOrFallback(zero, nonZero); got.Range.Start.Line != 5 {
		t.Error("should return fallback when primary is zero")
	}
}
