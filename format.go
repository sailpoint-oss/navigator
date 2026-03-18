package navigator

import (
	"bytes"
	"strings"
)

// FileFormat represents the serialization format of an OpenAPI document.
type FileFormat int

const (
	FormatUnknown FileFormat = iota
	FormatYAML
	FormatJSON
)

// String returns a human-readable name for the format.
func (f FileFormat) String() string {
	switch f {
	case FormatYAML:
		return "yaml"
	case FormatJSON:
		return "json"
	default:
		return "unknown"
	}
}

// FormatFromURI determines the file format from a URI or file path extension.
func FormatFromURI(uri string) FileFormat {
	lower := strings.ToLower(uri)
	if strings.HasSuffix(lower, ".json") {
		return FormatJSON
	}
	if strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml") {
		return FormatYAML
	}
	return FormatUnknown
}

// FormatFromContent sniffs the content to determine the format.
func FormatFromContent(content []byte) FileFormat {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return FormatUnknown
	}
	if trimmed[0] == '{' || trimmed[0] == '[' {
		return FormatJSON
	}
	return FormatYAML
}

// DetectFormat returns the format from URI if possible, falling back to content sniffing.
func DetectFormat(uri string, content []byte) FileFormat {
	if f := FormatFromURI(uri); f != FormatUnknown {
		return f
	}
	return FormatFromContent(content)
}
