package graph

import (
	"time"

	navigator "github.com/sailpoint-oss/navigator"
)

// StageName identifies a processing stage in the pipeline.
type StageName string

const (
	StageRaw      StageName = "raw"
	StageParse    StageName = "parse"
	StageLint     StageName = "lint"
	StageBind     StageName = "bind"
	StageValidate StageName = "validate"
	StageAnalyze  StageName = "analyze"
)

// StageResult holds the cached output of a pipeline stage.
type StageResult struct {
	Stage       StageName
	Data        interface{}
	Version     int64
	Diagnostics []Diagnostic
	Duration    time.Duration
}

// Diagnostic represents a problem found during processing. This type mirrors
// telescope's ctypes.Diagnostic closely enough for the graph pipeline, while
// remaining free of gossip/protocol dependencies.
type Diagnostic struct {
	URI             string
	Range           navigator.Range
	Severity        int
	Code            string
	CodeDescription string
	Source          string
	Message         string
	Tags            []int
	Related         []RelatedInformation
	Data            interface{}
}

// RelatedInformation provides additional context for a diagnostic.
type RelatedInformation struct {
	URI     string
	Range   navigator.Range
	Message string
}

// Fix describes a suggested code fix for a diagnostic.
type Fix struct {
	Title string
	Edits []TextEdit
}

// TextEdit describes a text change.
type TextEdit struct {
	Range   navigator.Range
	NewText string
}

// ParseOutput bundles the semantic parse results for a document. Stored as
// the Data field of a StageParse StageResult by telescope's pipeline.
type ParseOutput struct {
	SemanticNode *navigator.SemanticNode
	PointerIndex *navigator.PointerIndex
	VirtualDocs  []interface{}
}

// GraphNode holds all state for a single document in the workspace graph.
type GraphNode struct {
	Source       navigator.DocumentSource
	Version      int64
	Raw          []byte
	Index        *navigator.Index
	StageResults map[StageName]*StageResult
	DirtyStages  map[StageName]bool
	Diagnostics  []Diagnostic
}

func newGraphNode(src navigator.DocumentSource) *GraphNode {
	return &GraphNode{
		Source:       src,
		StageResults: make(map[StageName]*StageResult),
		DirtyStages:  make(map[StageName]bool),
	}
}
