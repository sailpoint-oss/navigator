package navigator

// OperationRef links an operation to its path and method.
type OperationRef struct {
	Path      string
	Method    string
	Operation *Operation
}

// RefUsage tracks where a $ref is used within a document.
type RefUsage struct {
	URI    string // document URI where this ref appears
	Loc    Loc
	Target string // the $ref target value (e.g. "#/components/schemas/Foo")
	From   string // JSON Pointer to the location of the $ref in the source
}

// RefEdge represents a directed $ref dependency between two files.
type RefEdge struct {
	FromURI     string
	FromPointer string
	ToURI       string
	ToPointer   string
	RefValue    string // raw $ref value as written in source
}
