package navigator

// Severity classifies how serious an Issue is for tooling and UX.
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
)

// IssueCategory groups issues by validation layer.
type IssueCategory int

const (
	CategorySyntax IssueCategory = iota
	CategoryStructural
	CategorySchema
	CategoryMeta
)

// Issue is a single diagnostic produced during parsing / validation.
// Codes are stable, lowercase, dot-separated identifiers (e.g. "structural.missing-responses").
type Issue struct {
	Code         string
	Message      string
	Pointer      string // JSON Pointer (RFC 6901), "" if unknown
	Range        Range
	Severity     Severity
	Category     IssueCategory
	DocumentKind DocumentKind
	SpecVersion  Version // document version when known
}
