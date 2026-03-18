package navigator

import "testing"

func TestSemanticNode_GetAndStringValue(t *testing.T) {
	root := &SemanticNode{
		Kind: NodeMapping,
		Children: map[string]*SemanticNode{
			"name": {Kind: NodeScalar, Value: "hello", Key: "name"},
			"count": {Kind: NodeScalar, Value: "42", Key: "count"},
		},
	}

	n := root.Get("name")
	if n == nil || n.StringValue() != "hello" {
		t.Error("Get/StringValue failed for name")
	}

	if root.Get("missing") != nil {
		t.Error("Get should return nil for missing key")
	}
}

func TestSemanticNode_Index(t *testing.T) {
	root := &SemanticNode{
		Kind: NodeSequence,
		Items: []*SemanticNode{
			{Kind: NodeScalar, Value: "first"},
			{Kind: NodeScalar, Value: "second"},
		},
	}

	if root.Index(0).StringValue() != "first" {
		t.Error("Index(0) failed")
	}
	if root.Index(1).StringValue() != "second" {
		t.Error("Index(1) failed")
	}
	if root.Index(2) != nil {
		t.Error("Index(2) should be nil")
	}
	if root.Index(-1) != nil {
		t.Error("Index(-1) should be nil")
	}
}

func TestSemanticNode_Walk(t *testing.T) {
	root := &SemanticNode{
		Kind: NodeMapping,
		Children: map[string]*SemanticNode{
			"a": {Kind: NodeScalar, Value: "1", Key: "a"},
			"b": {
				Kind: NodeMapping,
				Key:  "b",
				Children: map[string]*SemanticNode{
					"c": {Kind: NodeScalar, Value: "2", Key: "c"},
				},
			},
		},
	}

	visited := make(map[string]bool)
	root.Walk(func(path string, node *SemanticNode) {
		visited[path] = true
	})

	if !visited[""] {
		t.Error("root not visited")
	}
	if !visited["/a"] {
		t.Error("/a not visited")
	}
	if !visited["/b"] {
		t.Error("/b not visited")
	}
	if !visited["/b/c"] {
		t.Error("/b/c not visited")
	}
}

func TestSemanticNode_Len(t *testing.T) {
	mapping := &SemanticNode{
		Kind:     NodeMapping,
		Children: map[string]*SemanticNode{"a": {}, "b": {}},
	}
	if mapping.Len() != 2 {
		t.Errorf("mapping Len = %d", mapping.Len())
	}

	seq := &SemanticNode{
		Kind:  NodeSequence,
		Items: []*SemanticNode{{}, {}, {}},
	}
	if seq.Len() != 3 {
		t.Errorf("sequence Len = %d", seq.Len())
	}

	scalar := &SemanticNode{Kind: NodeScalar}
	if scalar.Len() != 0 {
		t.Errorf("scalar Len = %d", scalar.Len())
	}
}

func TestSemanticNode_NilSafety(t *testing.T) {
	var n *SemanticNode
	if n.Get("x") != nil {
		t.Error("nil Get should return nil")
	}
	if n.Index(0) != nil {
		t.Error("nil Index should return nil")
	}
	if n.StringValue() != "" {
		t.Error("nil StringValue should be empty")
	}
	if n.Len() != 0 {
		t.Error("nil Len should be 0")
	}
}
