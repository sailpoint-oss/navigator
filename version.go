package navigator

import "strings"

// Version represents a known OpenAPI specification version.
type Version string

const (
	VersionUnknown Version = ""
	Version20      Version = "2.0"
	Version30      Version = "3.0"
	Version31      Version = "3.1"
	Version32      Version = "3.2"
)

// VersionFromString parses a version string (e.g. "3.1.0") into a Version constant.
func VersionFromString(s string) Version {
	if len(s) == 0 {
		return VersionUnknown
	}
	switch {
	case s == "2.0" || (len(s) >= 2 && s[:2] == "2."):
		return Version20
	case len(s) >= 3 && s[:3] == "3.0":
		return Version30
	case len(s) >= 3 && s[:3] == "3.1":
		return Version31
	case len(s) >= 3 && s[:3] == "3.2":
		return Version32
	default:
		return VersionUnknown
	}
}

// Format is used in rule metadata to indicate which spec versions a rule applies to.
type Format string

const (
	FormatOAS2  Format = "oas2"
	FormatOAS3  Format = "oas3"
	FormatOAS31 Format = "oas3.1"
	FormatOAS32 Format = "oas3.2"
)

// FormatsForVersion returns the Format tags that apply to a given version.
func FormatsForVersion(v Version) []Format {
	switch v {
	case Version20:
		return []Format{FormatOAS2}
	case Version30:
		return []Format{FormatOAS3}
	case Version31:
		return []Format{FormatOAS3, FormatOAS31}
	case Version32:
		return []Format{FormatOAS3, FormatOAS32}
	default:
		return nil
	}
}

// AtLeast returns true if this version is >= the given major.minor.patch.
func (v Version) AtLeast(major, minor, patch int) bool {
	pv := parseVersionParts(string(v))
	return compareParts(pv, [3]int{major, minor, patch}) >= 0
}

// IsZero returns true if the version has not been set.
func (v Version) IsZero() bool {
	return v == VersionUnknown
}

// DocType classifies a document.
type DocType int

const (
	DocTypeUnknown    DocType = iota
	DocTypeRoot               // A root OpenAPI document (has openapi or swagger key)
	DocTypeFragment           // A partial document referenced via $ref or with recognized fragment shape
	DocTypeNonOpenAPI         // Not an OpenAPI document (no root key, no fragment shape, not a graph member)
)

// String returns a human-readable name for the DocType.
func (d DocType) String() string {
	switch d {
	case DocTypeRoot:
		return "root"
	case DocTypeFragment:
		return "fragment"
	case DocTypeNonOpenAPI:
		return "non-openapi"
	default:
		return "unknown"
	}
}

// DetectDocType classifies a document using root-level keys and optional graph
// membership. A file with an openapi/swagger key is always Root. A known $ref
// target (isGraphMember) or a file whose keys match a recognized OpenAPI
// fragment shape is Fragment. Everything else is NonOpenAPI.
func DetectDocType(rootKeys []string, isGraphMember bool) DocType {
	for _, k := range rootKeys {
		switch k {
		case "openapi", "swagger":
			return DocTypeRoot
		}
	}
	if isGraphMember {
		return DocTypeFragment
	}
	if DetectFragmentType(rootKeys) != FragmentUnknown {
		return DocTypeFragment
	}
	return DocTypeNonOpenAPI
}

func parseVersionParts(s string) [3]int {
	var parts [3]int
	segs := strings.SplitN(s, ".", 3)
	for i, seg := range segs {
		if i >= 3 {
			break
		}
		for _, ch := range seg {
			if ch >= '0' && ch <= '9' {
				parts[i] = parts[i]*10 + int(ch-'0')
			} else {
				break
			}
		}
	}
	return parts
}

func compareParts(a, b [3]int) int {
	for i := 0; i < 3; i++ {
		if a[i] != b[i] {
			if a[i] > b[i] {
				return 1
			}
			return -1
		}
	}
	return 0
}
