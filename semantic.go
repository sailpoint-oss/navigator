package navigator

import (
	"fmt"
	"strconv"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// NodeKind classifies a SemanticNode.
type NodeKind int

const (
	NodeMapping  NodeKind = iota
	NodeSequence
	NodeScalar
	NodeNull
)

// SemanticNode is a format-independent intermediate representation of a
// YAML or JSON document tree. It strips away format-specific wrapping
// (block_mapping_pair, flow_pair, pair, etc.) to expose a uniform
// mapping/sequence/scalar structure.
type SemanticNode struct {
	Kind     NodeKind
	Value    any
	Range    Range
	Key      string
	Children map[string]*SemanticNode
	Items    []*SemanticNode
	Tag      string
	Anchor   string
	Alias    string
	CST      *sitter.Node
}

// Get retrieves a child by key (for mappings).
func (n *SemanticNode) Get(key string) *SemanticNode {
	if n == nil || n.Children == nil {
		return nil
	}
	return n.Children[key]
}

// Index retrieves a child by index (for sequences).
func (n *SemanticNode) Index(i int) *SemanticNode {
	if n == nil || i < 0 || i >= len(n.Items) {
		return nil
	}
	return n.Items[i]
}

// StringValue returns the scalar value as a string.
func (n *SemanticNode) StringValue() string {
	if n == nil {
		return ""
	}
	if n.Kind == NodeScalar {
		if s, ok := n.Value.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", n.Value)
	}
	return ""
}

// Len returns the number of children (mapping) or items (sequence).
func (n *SemanticNode) Len() int {
	if n == nil {
		return 0
	}
	switch n.Kind {
	case NodeMapping:
		return len(n.Children)
	case NodeSequence:
		return len(n.Items)
	default:
		return 0
	}
}

// Walk visits all nodes depth-first, calling fn with the JSON pointer path
// and each node.
func (n *SemanticNode) Walk(fn func(path string, node *SemanticNode)) {
	if n == nil {
		return
	}
	n.walk("", fn)
}

func (n *SemanticNode) walk(path string, fn func(string, *SemanticNode)) {
	fn(path, n)
	switch n.Kind {
	case NodeMapping:
		for key, child := range n.Children {
			child.walk(path+"/"+EscapeJSONPointer(key), fn)
		}
	case NodeSequence:
		for i, item := range n.Items {
			item.walk(path+"/"+strconv.Itoa(i), fn)
		}
	}
}
