package navigator

// PointerIndex maps JSON Pointer strings to source ranges. Used to quickly
// locate any element within a document by its pointer path.
type PointerIndex struct {
	entries map[string]Range
}

// NewPointerIndex creates an empty pointer index.
func NewPointerIndex() *PointerIndex {
	return &PointerIndex{entries: make(map[string]Range)}
}

// Set records a mapping from a JSON pointer to a source range.
func (p *PointerIndex) Set(pointer string, r Range) {
	if p.entries == nil {
		p.entries = make(map[string]Range)
	}
	p.entries[pointer] = r
}

// Get returns the range for a JSON pointer, or false if not found.
func (p *PointerIndex) Get(pointer string) (Range, bool) {
	if p == nil || p.entries == nil {
		return Range{}, false
	}
	r, ok := p.entries[pointer]
	return r, ok
}

// All returns a copy of all entries.
func (p *PointerIndex) All() map[string]Range {
	if p == nil || p.entries == nil {
		return nil
	}
	out := make(map[string]Range, len(p.entries))
	for k, v := range p.entries {
		out[k] = v
	}
	return out
}

// Len returns the number of entries.
func (p *PointerIndex) Len() int {
	if p == nil || p.entries == nil {
		return 0
	}
	return len(p.entries)
}

// BuildPointerIndex walks a SemanticNode tree and builds a PointerIndex.
func BuildPointerIndex(root *SemanticNode) *PointerIndex {
	idx := NewPointerIndex()
	if root == nil {
		return idx
	}
	root.Walk(func(pointer string, node *SemanticNode) {
		if node != nil {
			idx.Set(pointer, node.Range)
		}
	})
	return idx
}
