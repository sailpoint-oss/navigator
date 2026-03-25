package navigator

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// Index provides fast lookups into a parsed OpenAPI document.
//
// Document may be nil when the source is not a YAML/JSON mapping (e.g. a bare
// sequence or scalar). In that case Issues will contain a structural diagnostic
// and the lookup maps will be empty. Callers should check IsOpenAPI before
// assuming Document fields are populated.
type Index struct {
	Document         *Document
	// Issues holds syntax, structural, and schema-shape diagnostics produced
	// during parsing. Always populated when an Index is returned from any parse
	// function; use Revalidate to re-run with custom ValidationOptions.
	Issues           []Issue
	Operations       map[string]*OperationRef   // operationId -> ref
	OperationsByPath map[string][]OperationRef  // path template -> operations
	Schemas          map[string]*Schema         // component name -> schema
	Parameters       map[string]*Parameter      // component name -> parameter
	Responses        map[string]*Response       // component name -> response
	SecuritySchemes  map[string]*SecurityScheme // scheme name -> scheme
	Refs             map[string][]RefUsage      // $ref target -> usages
	AllRefs          []RefUsage                 // all $ref usages in the doc
	Tags             map[string]*Tag            // tag name -> tag
	Version          Version
	Format           FileFormat

	uri          string
	content      []byte
	tree         *tree_sitter.Tree
	semanticRoot *SemanticNode
}

// URI returns the document URI associated with this index.
func (idx *Index) URI() string {
	if idx == nil {
		return ""
	}
	return idx.uri
}

// Content returns the raw bytes of the source document.
func (idx *Index) Content() []byte {
	if idx == nil {
		return nil
	}
	return idx.content
}

// Tree returns the tree-sitter CST (nil for standalone-parsed indexes).
func (idx *Index) Tree() *tree_sitter.Tree {
	if idx == nil {
		return nil
	}
	return idx.tree
}

// SemanticRoot returns the SemanticNode IR root (nil for standalone-parsed indexes).
func (idx *Index) SemanticRoot() *SemanticNode {
	if idx == nil {
		return nil
	}
	return idx.semanticRoot
}

// IsOpenAPI returns true if the index represents a valid root OpenAPI document.
func (idx *Index) IsOpenAPI() bool {
	return idx != nil && idx.Document != nil && idx.Document.DocType == DocTypeRoot
}

// HasPath returns true if the given path template exists.
func (idx *Index) HasPath(path string) bool {
	if idx == nil || idx.Document == nil {
		return false
	}
	_, ok := idx.Document.Paths[path]
	return ok
}

// AllOperations returns all indexed operations.
func (idx *Index) AllOperations() []*OperationRef {
	if idx == nil {
		return nil
	}
	ops := make([]*OperationRef, 0, len(idx.Operations))
	for _, op := range idx.Operations {
		ops = append(ops, op)
	}
	return ops
}

// SchemaNames returns all component schema names.
func (idx *Index) SchemaNames() []string {
	if idx == nil {
		return nil
	}
	names := make([]string, 0, len(idx.Schemas))
	for name := range idx.Schemas {
		names = append(names, name)
	}
	return names
}

// ComponentNames returns the names of all components of a given kind.
func (idx *Index) ComponentNames(kind string) []string {
	if idx == nil || idx.Document == nil || idx.Document.Components == nil {
		return nil
	}
	switch kind {
	case "schemas":
		return mapKeys(idx.Document.Components.Schemas)
	case "responses":
		return mapKeys(idx.Document.Components.Responses)
	case "parameters":
		return mapKeys(idx.Document.Components.Parameters)
	case "examples":
		return mapKeys(idx.Document.Components.Examples)
	case "requestBodies":
		return mapKeys(idx.Document.Components.RequestBodies)
	case "headers":
		return mapKeys(idx.Document.Components.Headers)
	case "securitySchemes":
		return mapKeys(idx.Document.Components.SecuritySchemes)
	case "links":
		return mapKeys(idx.Document.Components.Links)
	default:
		return nil
	}
}

// RefsTo returns all usages that reference the given target.
func (idx *Index) RefsTo(target string) []RefUsage {
	if idx == nil {
		return nil
	}
	return idx.Refs[target]
}

func mapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// NewIndexFromDocument builds a fully indexed Index from an existing Document.
// Useful when the Document was obtained from another source (e.g. a telescope
// SDK Index) and you need navigator's ResolveRef and lookup methods.
// Issues is always populated via default validation.
func NewIndexFromDocument(doc *Document) *Index {
	idx := newEmptyIndex()
	idx.Document = doc
	if doc != nil {
		idx.Version = doc.ParsedVersion
	}
	idx.indexPaths()
	idx.indexComponents()
	idx.indexTags()
	idx.applyDefaultValidation()
	return idx
}

func newEmptyIndex() *Index {
	return &Index{
		Operations:       make(map[string]*OperationRef),
		OperationsByPath: make(map[string][]OperationRef),
		Schemas:          make(map[string]*Schema),
		Parameters:       make(map[string]*Parameter),
		Responses:        make(map[string]*Response),
		SecuritySchemes:  make(map[string]*SecurityScheme),
		Refs:             make(map[string][]RefUsage),
		Tags:             make(map[string]*Tag),
	}
}

func (idx *Index) indexPaths() {
	if idx.Document == nil {
		return
	}
	for path, item := range idx.Document.Paths {
		if item == nil {
			continue
		}
		var ops []OperationRef
		for _, mo := range item.Operations() {
			ref := OperationRef{
				Path:      path,
				Method:    mo.Method,
				Operation: mo.Operation,
			}
			ops = append(ops, ref)
			if mo.Operation.OperationID != "" {
				idx.Operations[mo.Operation.OperationID] = &ref
			}
		}
		idx.OperationsByPath[path] = ops
	}
}

func (idx *Index) indexComponents() {
	if idx.Document == nil || idx.Document.Components == nil {
		return
	}
	c := idx.Document.Components
	for name, s := range c.Schemas {
		idx.Schemas[name] = s
	}
	for name, p := range c.Parameters {
		idx.Parameters[name] = p
	}
	for name, r := range c.Responses {
		idx.Responses[name] = r
	}
	for name, ss := range c.SecuritySchemes {
		idx.SecuritySchemes[name] = ss
	}
}

func (idx *Index) indexTags() {
	if idx.Document == nil {
		return
	}
	for i := range idx.Document.Tags {
		t := &idx.Document.Tags[i]
		idx.Tags[t.Name] = t
	}
}
