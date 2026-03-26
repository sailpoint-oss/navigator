package navigator

import (
	"fmt"
	"strings"
)

// ResolveRef resolves a $ref string to a model element within this index.
// Only local references (starting with #) are supported; cross-file
// resolution requires CrossFileResolver.
func (idx *Index) ResolveRef(ref string) (interface{}, error) {
	if idx == nil {
		return nil, fmt.Errorf("resolve: nil index")
	}
	if ref == "" {
		return nil, fmt.Errorf("resolve: empty ref")
	}
	if !strings.HasPrefix(ref, "#") {
		// Strip the file path for cross-file refs, resolving only the fragment.
		if i := strings.Index(ref, "#"); i >= 0 {
			ref = ref[i:]
		} else {
			return nil, fmt.Errorf("resolve: external ref %q has no fragment", ref)
		}
	}
	ref = strings.TrimPrefix(ref, "#")
	if ref == "" || ref == "/" {
		return idx.PrimaryValue(), nil
	}
	ref = strings.TrimPrefix(ref, "/")
	parts := strings.Split(ref, "/")
	for i, p := range parts {
		parts[i] = unescapeJSONPointer(p)
	}
	if idx.Document != nil && idx.Document.DocType == DocTypeRoot {
		if val, err := idx.resolveRefParts(parts); err == nil {
			return val, nil
		}
	}
	return idx.resolveSemanticRefParts(parts)
}

func (idx *Index) resolveRefParts(parts []string) (interface{}, error) {
	if len(parts) == 0 || idx.Document == nil {
		return nil, fmt.Errorf("resolve: cannot resolve empty path")
	}
	switch parts[0] {
	case "components":
		return idx.resolveComponent(parts[1:])
	case "paths":
		return idx.resolvePath(parts[1:])
	case "info":
		if len(parts) == 1 {
			return idx.Document.Info, nil
		}
	case "servers":
		return idx.resolveByIndex(parts[1:], len(idx.Document.Servers), func(i int) interface{} {
			return &idx.Document.Servers[i]
		})
	case "tags":
		return idx.resolveByIndex(parts[1:], len(idx.Document.Tags), func(i int) interface{} {
			return &idx.Document.Tags[i]
		})
	}
	return nil, fmt.Errorf("resolve: unknown root key %q", parts[0])
}

func (idx *Index) resolveComponent(parts []string) (interface{}, error) {
	if len(parts) < 2 || idx.Document.Components == nil {
		return nil, fmt.Errorf("resolve: components ref needs at least kind/name")
	}
	if len(parts) > 2 {
		return nil, fmt.Errorf("resolve: deep component path requires semantic resolution")
	}
	kind, name := parts[0], parts[1]
	switch kind {
	case "schemas":
		if s, ok := idx.Document.Components.Schemas[name]; ok {
			return s, nil
		}
	case "responses":
		if r, ok := idx.Document.Components.Responses[name]; ok {
			return r, nil
		}
	case "parameters":
		if p, ok := idx.Document.Components.Parameters[name]; ok {
			return p, nil
		}
	case "examples":
		if e, ok := idx.Document.Components.Examples[name]; ok {
			return e, nil
		}
	case "requestBodies":
		if rb, ok := idx.Document.Components.RequestBodies[name]; ok {
			return rb, nil
		}
	case "headers":
		if h, ok := idx.Document.Components.Headers[name]; ok {
			return h, nil
		}
	case "securitySchemes":
		if ss, ok := idx.Document.Components.SecuritySchemes[name]; ok {
			return ss, nil
		}
	case "links":
		if l, ok := idx.Document.Components.Links[name]; ok {
			return l, nil
		}
	case "pathItems":
		if pi, ok := idx.Document.Components.PathItems[name]; ok {
			return pi, nil
		}
	default:
		return nil, fmt.Errorf("resolve: unknown component kind %q", kind)
	}
	return nil, fmt.Errorf("resolve: component %s/%s not found", kind, name)
}

func (idx *Index) resolvePath(parts []string) (interface{}, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("resolve: paths ref needs a path template")
	}
	// Unescape JSON Pointer before reconstructing the path template
	decoded := unescapeJSONPointer(parts[0])
	pathTemplate := decoded
	if !strings.HasPrefix(pathTemplate, "/") {
		pathTemplate = "/" + pathTemplate
	}
	item, ok := idx.Document.Paths[pathTemplate]
	if !ok {
		return nil, fmt.Errorf("resolve: path %q not found", pathTemplate)
	}
	if len(parts) == 1 {
		return item, nil
	}
	method := strings.ToLower(parts[1])
	if len(parts) > 2 {
		return nil, fmt.Errorf("resolve: deep path pointer requires semantic resolution")
	}
	for _, mo := range item.Operations() {
		if mo.Method == method {
			return mo.Operation, nil
		}
	}
	switch method {
	case "parameters":
		return item.Parameters, nil
	}
	return nil, fmt.Errorf("resolve: method or key %q not found on path %q", parts[1], pathTemplate)
}

func (idx *Index) resolveByIndex(parts []string, length int, getter func(int) interface{}) (interface{}, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("resolve: need an array index")
	}
	if len(parts) > 1 {
		return nil, fmt.Errorf("resolve: deep array pointer requires semantic resolution")
	}
	i := 0
	for _, ch := range parts[0] {
		if ch < '0' || ch > '9' {
			return nil, fmt.Errorf("resolve: invalid array index %q", parts[0])
		}
		i = i*10 + int(ch-'0')
	}
	if i >= length {
		return nil, fmt.Errorf("resolve: index %d out of range (length %d)", i, length)
	}
	return getter(i), nil
}

// unescapeJSONPointer decodes RFC 6901 escapes.
func unescapeJSONPointer(s string) string {
	s = strings.ReplaceAll(s, "~1", "/")
	s = strings.ReplaceAll(s, "~0", "~")
	return s
}

// EscapeJSONPointer encodes a string for use in a JSON Pointer segment.
func EscapeJSONPointer(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}

// Resolve is an alias for ResolveRef, matching the method name used by
// telescope's LSP handlers and rules engine.
func (idx *Index) Resolve(ref string) (interface{}, error) {
	return idx.ResolveRef(ref)
}

// ComponentRefPath returns a $ref path for a component, e.g. "#/components/schemas/Pet".
func ComponentRefPath(kind, name string) string {
	return "#/components/" + kind + "/" + EscapeJSONPointer(name)
}
