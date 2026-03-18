package navigator

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// Position represents a zero-based line/character position in a document.
type Position struct {
	Line      uint32
	Character uint32
}

// Range represents a span between two positions in a document.
type Range struct {
	Start Position
	End   Position
}

// ContainsPosition returns true if pos is within r (inclusive start, exclusive end).
func ContainsPosition(r Range, pos Position) bool {
	if pos.Line < r.Start.Line || pos.Line > r.End.Line {
		return false
	}
	if pos.Line == r.Start.Line && pos.Character < r.Start.Character {
		return false
	}
	if pos.Line == r.End.Line && pos.Character >= r.End.Character {
		return false
	}
	return true
}

// IsEmpty returns true if the range has zero extent.
func IsEmpty(r Range) bool {
	return r.Start.Line == r.End.Line && r.Start.Character == r.End.Character
}

// FileStartRange is a zero-length range at the start of a file.
var FileStartRange = Range{
	Start: Position{Line: 0, Character: 0},
	End:   Position{Line: 0, Character: 0},
}

// Loc tracks source location for any model element, linking back to the
// tree-sitter node that produced it. Node is nil for standalone-parsed indexes.
type Loc struct {
	Range Range
	Node  *tree_sitter.Node
}

// LocFromNode creates a Loc from a tree-sitter node.
func LocFromNode(node *tree_sitter.Node) Loc {
	if node == nil {
		return Loc{}
	}
	start := node.StartPosition()
	end := node.EndPosition()
	return Loc{
		Range: Range{
			Start: Position{Line: uint32(start.Row), Character: uint32(start.Column)},
			End:   Position{Line: uint32(end.Row), Character: uint32(end.Column)},
		},
		Node: node,
	}
}

// LocOrFallback returns primary if it has a non-zero range, otherwise fallback.
func LocOrFallback(primary, fallback Loc) Loc {
	if primary.Range.Start.Line == 0 && primary.Range.Start.Character == 0 &&
		primary.Range.End.Line == 0 && primary.Range.End.Character == 0 {
		return fallback
	}
	return primary
}

// DescriptionValue holds a decoded description string alongside source
// geometry needed for translating markdown-relative positions back to
// absolute ranges.
type DescriptionValue struct {
	Text       string
	Loc        Loc
	LineOffset int // 0 for plain/quoted scalars, 1 for block scalars (| or >)
	IndentCols int // leading whitespace stripped from block scalar content lines
}

// Node is a lightweight wrapper for raw values that haven't been typed into a
// specific model struct (e.g. examples, extensions, defaults).
type Node struct {
	Value   string
	RawNode *tree_sitter.Node
	Loc     Loc
}
