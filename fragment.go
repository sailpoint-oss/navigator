package navigator

// FragmentType classifies what kind of OpenAPI fragment a non-root file contains.
type FragmentType int

const (
	FragmentUnknown        FragmentType = iota
	FragmentSchema                      // Standalone schema object
	FragmentPathItem                    // Standalone path item
	FragmentOperation                   // Standalone operation
	FragmentParameter                   // Standalone parameter
	FragmentRequestBody                 // Standalone request body
	FragmentResponse                    // Standalone response
	FragmentHeader                      // Standalone header
	FragmentSecurityScheme              // Standalone security scheme
	FragmentComponents                  // Standalone components object
	FragmentServer                      // Standalone server object
)

// String returns a human-readable name for the fragment type.
func (f FragmentType) String() string {
	switch f {
	case FragmentSchema:
		return "schema"
	case FragmentPathItem:
		return "path-item"
	case FragmentOperation:
		return "operation"
	case FragmentParameter:
		return "parameter"
	case FragmentRequestBody:
		return "request-body"
	case FragmentResponse:
		return "response"
	case FragmentHeader:
		return "header"
	case FragmentSecurityScheme:
		return "security-scheme"
	case FragmentComponents:
		return "components"
	case FragmentServer:
		return "server"
	default:
		return "unknown"
	}
}

// DetectFragmentType examines root-level keys of a non-root document to
// determine what kind of OpenAPI fragment it represents.
func DetectFragmentType(rootKeys []string) FragmentType {
	keySet := make(map[string]bool, len(rootKeys))
	for _, k := range rootKeys {
		keySet[k] = true
	}

	// If it looks like a root doc, it's not a fragment.
	if keySet["openapi"] || keySet["swagger"] {
		return FragmentUnknown
	}

	if keySet["components"] && !keySet["type"] && !keySet["properties"] {
		return FragmentComponents
	}

	// Schema indicators
	if keySet["type"] || keySet["properties"] || keySet["allOf"] || keySet["anyOf"] ||
		keySet["oneOf"] || keySet["items"] || keySet["$ref"] || keySet["enum"] {
		return FragmentSchema
	}

	// PathItem indicators
	if keySet["get"] || keySet["put"] || keySet["post"] || keySet["delete"] ||
		keySet["options"] || keySet["head"] || keySet["patch"] || keySet["trace"] {
		return FragmentPathItem
	}

	// Operation indicators
	if keySet["operationId"] || keySet["responses"] {
		return FragmentOperation
	}

	// Parameter indicators
	if keySet["in"] && keySet["name"] {
		return FragmentParameter
	}

	// RequestBody indicators
	if keySet["content"] && keySet["required"] && !keySet["description"] {
		return FragmentRequestBody
	}
	if keySet["content"] && !keySet["url"] && !keySet["in"] {
		return FragmentResponse
	}

	// SecurityScheme indicators
	if keySet["type"] && (keySet["scheme"] || keySet["flows"] || keySet["openIdConnectUrl"]) {
		return FragmentSecurityScheme
	}

	// Header indicators
	if keySet["schema"] && !keySet["in"] && !keySet["name"] {
		return FragmentHeader
	}

	// Server indicators
	if keySet["url"] && !keySet["type"] {
		return FragmentServer
	}

	return FragmentUnknown
}
