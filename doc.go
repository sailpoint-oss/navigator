// Package navigator provides uniform OpenAPI file handling for the entire
// workspace toolchain. It parses, indexes, and resolves $refs across single-file
// and massively multi-file OpenAPI specifications (Swagger 2.0, OpenAPI 3.0–3.2).
//
// Navigator is consumed by telescope (LSP), barometer (contract testing),
// and cartographer (spec extraction/bundling). It is the single source of truth
// for OpenAPI document models, parsing, indexing, and cross-file resolution.
//
// # Design Principles
//
// Lazy by default: files are loaded and parsed only when needed. [OpenProject]
// loads the root document; transitive $ref targets load on demand. [Project.LoadAll]
// is the explicit break-glass to force-load everything.
//
// Break-glass at every level: every abstraction exposes its internals. An [Index]
// provides [Index.Tree] for the raw tree-sitter CST, [Index.SemanticRoot] for the IR, and
// [Index.Document] for the typed model. A [Project] exposes its [Project.Cache],
// [Project.Graph], and [Project.Resolver] for direct manipulation.
//
// Workspace-wide shared state: [IndexCache] and [FileGraph] are singletons shared
// across multiple roots. A [Project] is a thin lens over this shared state.
//
// # Quick Start
//
// Simple single-file parsing (barometer pattern):
//
//	idx := navigator.Parse(content)
//	for _, op := range idx.AllOperations() {
//	    fmt.Println(op.Operation.OperationID)
//	}
//
// Multi-file project (cartographer pattern):
//
//	project, err := navigator.OpenProject("path/to/root.yaml")
//	if err != nil { ... }
//	project.LoadAll() // eagerly loads all transitive deps
//	result, err := project.Resolve(project.RootURI(), "./schemas/pet.yaml")
//
// Full workspace graph (telescope pattern):
//
//	import "github.com/sailpoint-oss/navigator/graph"
//	ws := navigator.NewWorkspace()
//	g := graph.New(ws.Cache, ws.Graph)
//	g.AddSource(navigator.NewMemorySource(uri, content))
//	runner := graph.NewPipelineRunner(graph.DefaultStages()...)
//	runner.RunAll(ctx, g, uri)
//
// # Tree-sitter Integration
//
// When CGO is available, [ParseContent] and [ParseTree] use the tree-sitter
// grammar from tree-sitter-openapi for precise source locations on every model
// element. The standalone [Parse] and [ParseURI] functions use gopkg.in/yaml.v3
// and require no CGO, making them suitable for CLI tools and CI environments.
//
// # Package Layout
//
// The main navigator package contains all types, parsing, indexing, single-file
// and cross-file resolution, workspace infrastructure, and the project facade.
//
// The [github.com/sailpoint-oss/navigator/graph] sub-package adds LSP-grade
// workspace graph management with pipeline processing and snapshot support.
package navigator
