package navigator

// DocumentKind classifies the top-level API description family.
type DocumentKind int

const (
	DocumentKindUnknown DocumentKind = iota
	DocumentKindOpenAPI
	DocumentKindArazzo
)

// String returns a human-readable name for the document kind.
func (k DocumentKind) String() string {
	switch k {
	case DocumentKindOpenAPI:
		return "openapi"
	case DocumentKindArazzo:
		return "arazzo"
	default:
		return "unknown"
	}
}

// DetectDocumentKind classifies the document family from root keys.
func DetectDocumentKind(rootKeys []string) DocumentKind {
	if len(rootKeys) == 0 {
		return DocumentKindUnknown
	}
	if looksLikeArazzoKeys(rootKeys) {
		return DocumentKindArazzo
	}
	return DocumentKindOpenAPI
}

func looksLikeArazzoKeys(rootKeys []string) bool {
	keySet := make(map[string]bool, len(rootKeys))
	for _, key := range rootKeys {
		keySet[key] = true
	}
	if keySet["arazzo"] {
		return true
	}
	return keySet["sourceDescriptions"] && keySet["workflows"]
}
