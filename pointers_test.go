package navigator

import "testing"

func TestPointerIndex_SetGet(t *testing.T) {
	pi := NewPointerIndex()
	r := Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 1, Character: 10}}
	pi.Set("/components/schemas/Pet", r)

	got, ok := pi.Get("/components/schemas/Pet")
	if !ok {
		t.Fatal("should find pointer")
	}
	if got.Start.Line != 1 {
		t.Errorf("line = %d", got.Start.Line)
	}
}

func TestPointerIndex_Missing(t *testing.T) {
	pi := NewPointerIndex()
	_, ok := pi.Get("/missing")
	if ok {
		t.Error("should not find missing pointer")
	}
}

func TestBuildPointerIndex(t *testing.T) {
	root := &SemanticNode{
		Kind: NodeMapping,
		Range: Range{Start: Position{0, 0}, End: Position{10, 0}},
		Children: map[string]*SemanticNode{
			"info": {
				Kind:  NodeMapping,
				Key:   "info",
				Range: Range{Start: Position{1, 0}, End: Position{3, 0}},
				Children: map[string]*SemanticNode{
					"title": {Kind: NodeScalar, Value: "Test", Key: "title",
						Range: Range{Start: Position{2, 2}, End: Position{2, 8}}},
				},
			},
		},
	}

	pi := BuildPointerIndex(root)
	all := pi.All()
	if len(all) < 3 {
		t.Errorf("expected at least 3 entries, got %d", len(all))
	}

	if _, ok := pi.Get(""); !ok {
		t.Error("should have root pointer")
	}
	if _, ok := pi.Get("/info"); !ok {
		t.Error("should have /info pointer")
	}
	if _, ok := pi.Get("/info/title"); !ok {
		t.Error("should have /info/title pointer")
	}
}
