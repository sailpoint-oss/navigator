package navigator

// Document is the top-level OpenAPI document model.
type Document struct {
	Version       string // raw version string, e.g. "3.1.0"
	ParsedVersion Version
	DocType       DocType
	Info          *Info
	Servers       []Server
	Paths         map[string]*PathItem
	Components    *Components
	Security      []SecurityRequirement
	Tags          []Tag
	ExternalDocs  *ExternalDocs
	Schemes       []string // OAS 2.0 only
	Extensions    map[string]*Node
	Loc           Loc
}

// Info provides metadata about the API.
type Info struct {
	Title          string
	Description    DescriptionValue
	TermsOfService string
	Contact        *Contact
	License        *License
	Version        string
	Extensions     map[string]*Node
	Loc            Loc
	TitleLoc       Loc
	VersionLoc     Loc
}

// Contact information for the API.
type Contact struct {
	Name  string
	URL   string
	Email string
	Loc   Loc
}

// License information for the API.
type License struct {
	Name       string
	Identifier string
	URL        string
	Loc        Loc
}

// Server describes a server.
type Server struct {
	URL         string
	Description DescriptionValue
	Variables   map[string]*ServerVariable
	Extensions  map[string]*Node
	Loc         Loc
	URLLoc      Loc
}

// ServerVariable describes a server variable for URL template substitution.
type ServerVariable struct {
	Enum        []string
	Default     string
	Description DescriptionValue
	Loc         Loc
}

// PathItem describes the operations available on a single path.
type PathItem struct {
	Summary     string
	Description DescriptionValue
	Get         *Operation
	Put         *Operation
	Post        *Operation
	Delete      *Operation
	Options     *Operation
	Head        *Operation
	Patch       *Operation
	Trace       *Operation
	Parameters  []*Parameter
	Servers     []Server
	Ref         string
	Extensions  map[string]*Node
	Loc         Loc
	PathLoc     Loc
}

// MethodOperation pairs an HTTP method with its Operation.
type MethodOperation struct {
	Method    string
	Operation *Operation
}

// Operations returns all non-nil operations on this path item with their HTTP method.
func (p *PathItem) Operations() []MethodOperation {
	if p == nil {
		return nil
	}
	var ops []MethodOperation
	if p.Get != nil {
		ops = append(ops, MethodOperation{"get", p.Get})
	}
	if p.Put != nil {
		ops = append(ops, MethodOperation{"put", p.Put})
	}
	if p.Post != nil {
		ops = append(ops, MethodOperation{"post", p.Post})
	}
	if p.Delete != nil {
		ops = append(ops, MethodOperation{"delete", p.Delete})
	}
	if p.Options != nil {
		ops = append(ops, MethodOperation{"options", p.Options})
	}
	if p.Head != nil {
		ops = append(ops, MethodOperation{"head", p.Head})
	}
	if p.Patch != nil {
		ops = append(ops, MethodOperation{"patch", p.Patch})
	}
	if p.Trace != nil {
		ops = append(ops, MethodOperation{"trace", p.Trace})
	}
	return ops
}

// TagUsage tracks an individual tag reference within an operation's tags array.
type TagUsage struct {
	Name string
	Loc  Loc
}

// Operation describes a single API operation on a path.
type Operation struct {
	OperationID    string
	Summary        string
	Description    DescriptionValue
	Tags           []TagUsage
	Parameters     []*Parameter
	RequestBody    *RequestBody
	Responses      map[string]*Response
	Security       []SecurityRequirement
	Deprecated     bool
	Servers        []Server
	ExternalDocs   *ExternalDocs
	Extensions     map[string]*Node
	Loc            Loc
	MethodLoc      Loc
	OperationIDLoc Loc
	TagsLoc        Loc
	ResponsesLoc   Loc
	ParametersLoc  Loc
}

// TagNames returns the tag names as a plain string slice.
func (op *Operation) TagNames() []string {
	if op == nil {
		return nil
	}
	names := make([]string, len(op.Tags))
	for i, t := range op.Tags {
		names[i] = t.Name
	}
	return names
}

// HasTag returns the TagUsage for the given name, if present.
func (op *Operation) HasTag(name string) (TagUsage, bool) {
	if op == nil {
		return TagUsage{}, false
	}
	for _, t := range op.Tags {
		if t.Name == name {
			return t, true
		}
	}
	return TagUsage{}, false
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name            string
	In              string // "query", "header", "path", "cookie"
	Description     DescriptionValue
	Required        bool
	Deprecated      bool
	AllowEmptyValue bool
	Schema          *Schema
	Example         *Node
	Examples        map[string]*Example
	Ref             string
	Extensions      map[string]*Node
	Loc             Loc
	NameLoc         Loc
}

// RequestBody describes a request body.
type RequestBody struct {
	Description DescriptionValue
	Content     map[string]*MediaType
	Required    bool
	Ref         string
	Loc         Loc
	NameLoc     Loc
}

// MediaType describes a media type with schema and examples.
type MediaType struct {
	Schema   *Schema
	Example  *Node
	Examples map[string]*Example
	Loc      Loc
}

// Response describes a single response from an API operation.
type Response struct {
	Description DescriptionValue
	Headers     map[string]*Header
	Content     map[string]*MediaType
	Links       map[string]*Link
	Ref         string
	Extensions  map[string]*Node
	Loc         Loc
	CodeLoc     Loc
	NameLoc     Loc
	HeadersLoc  Loc
}

// Header describes a single header.
type Header struct {
	Description DescriptionValue
	Required    bool
	Deprecated  bool
	Schema      *Schema
	Ref         string
	Loc         Loc
	NameLoc     Loc
}

// Link describes a possible design-time link for a response.
type Link struct {
	OperationRef   string
	OperationID    string
	Description    DescriptionValue
	Ref            string
	Loc            Loc
	NameLoc        Loc
	OperationIDLoc Loc
}

// Tag adds metadata to a single tag used by operations.
type Tag struct {
	Name         string
	Description  DescriptionValue
	ExternalDocs *ExternalDocs
	Extensions   map[string]*Node
	Loc          Loc
	NameLoc      Loc
}

// ExternalDocs points to external documentation.
type ExternalDocs struct {
	Description DescriptionValue
	URL         string
	Loc         Loc
	URLLoc      Loc
}

// SecurityRequirementEntry represents a single scheme entry within a security
// requirement object.
type SecurityRequirementEntry struct {
	Name      string
	Scopes    []string
	ScopeLocs []Loc
	NameLoc   Loc
}

// SecurityRequirement is a security requirement object (one element in the
// security array).
type SecurityRequirement struct {
	Entries []SecurityRequirementEntry
	Loc     Loc
}

// SchemeNames returns the scheme names referenced in this requirement.
func (sr SecurityRequirement) SchemeNames() []string {
	names := make([]string, len(sr.Entries))
	for i, e := range sr.Entries {
		names[i] = e.Name
	}
	return names
}

// HasScheme returns the entry for the given scheme name, if present.
func (sr SecurityRequirement) HasScheme(name string) (SecurityRequirementEntry, bool) {
	for _, e := range sr.Entries {
		if e.Name == name {
			return e, true
		}
	}
	return SecurityRequirementEntry{}, false
}
