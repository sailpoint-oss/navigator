package navigator

import "testing"

func TestSemanticNodeAt_TraversesMappingsAndSequences(t *testing.T) {
	root := &SemanticNode{
		Kind: NodeMapping,
		Children: map[string]*SemanticNode{
			"items": {
				Kind: NodeSequence,
				Items: []*SemanticNode{
					{Kind: NodeScalar, Value: "first"},
					{Kind: NodeScalar, Value: "second"},
				},
			},
		},
	}

	node, err := semanticNodeAt(root, []string{"items", "1"})
	if err != nil {
		t.Fatalf("semanticNodeAt: %v", err)
	}
	if node.StringValue() != "second" {
		t.Fatalf("StringValue = %q, want second", node.StringValue())
	}
}

func TestSemanticNodeAt_ReportsScalarTraversalError(t *testing.T) {
	root := &SemanticNode{
		Kind: NodeMapping,
		Children: map[string]*SemanticNode{
			"info": {Kind: NodeScalar, Value: "title"},
		},
	}

	if _, err := semanticNodeAt(root, []string{"info", "title"}); err == nil {
		t.Fatal("expected scalar traversal error")
	}
}

func TestResolveSemanticRefParts_ReturnsTypedInfo(t *testing.T) {
	idx := ParseContent([]byte(`openapi: "3.1.0"
info:
  title: Semantic API
  version: "1.0.0"
paths: {}
`), "file:///semantic.yaml")
	if idx == nil {
		t.Fatal("ParseContent returned nil")
	}

	val, err := idx.resolveSemanticRefParts([]string{"info"})
	if err != nil {
		t.Fatalf("resolveSemanticRefParts: %v", err)
	}
	info, ok := val.(*Info)
	if !ok {
		t.Fatalf("expected *Info, got %T", val)
	}
	if info.Title != "Semantic API" || info.Version != "1.0.0" {
		t.Fatalf("unexpected info: %+v", info)
	}
}
