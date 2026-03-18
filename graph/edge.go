package graph

// EdgeKind classifies the relationship between two graph nodes.
type EdgeKind int

const (
	EdgeRef      EdgeKind = iota // $ref relationship
	EdgeExternal                 // external document reference
	EdgePathRef                  // file-relative $ref edge with fragment
)

// Edge represents a directed relationship between two documents.
type Edge struct {
	SourceURI     string
	SourcePointer string
	TargetURI     string
	TargetPointer string
	Kind          EdgeKind
	RefValue      string // raw $ref value as written in source
}
