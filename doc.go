// Package navigator provides uniform OpenAPI and Arazzo file handling for the
// entire workspace toolchain. It parses, indexes, and resolves $refs across
// single-file and massively multi-file OpenAPI specifications (Swagger 2.0,
// OpenAPI 3.0–3.2), and it parses plus schema-validates Arazzo 1.0.x workflow
// documents through the same loading entrypoints.
//
// Navigator is consumed by editor, contract-testing, and downstream API
// tooling. It is the single source of truth for OpenAPI and Arazzo document
// models, parsing, indexing, and cross-file resolution where applicable.
//
// # Design Principles
//
// Lazy by default: files are loaded and parsed only when needed. [OpenProject]
// loads the root document; transitive $ref targets load on demand. [Project.LoadAll]
// is the explicit break-glass to force-load everything.
//
// Break-glass at every level: every abstraction exposes its internals. An [Index]
// provides [Index.Tree] for the raw tree-sitter CST, [Index.SemanticRoot] for
// the IR, and [Index.Document] / [Index.Arazzo] / [Index.PrimaryValue] for the
// typed model. A [Project] exposes its [Project.Cache], [Project.Graph], and
// [Project.Resolver] for direct manipulation.
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
// Multi-file project (workspace consumer pattern):
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
// # Validation
//
// Every successful parse populates [Index.Issues] with [Issue] values covering
// syntax (CST/yaml duplicate keys, tree-sitter ERROR nodes), structural rules,
// lightweight schema-shaped checks where applicable, and meta-schema
// validation. OpenAPI root and fragment schemas are embedded from generated
// Draft 2020-12 JSON Schemas under `schemas/ts/openapi/`; Arazzo uses the
// published Arazzo 1.0.x root schema. All meta validation runs through
// github.com/santhosh-tekuri/jsonschema/v6 with ECMA-262 pattern matching via
// github.com/dlclark/regexp2. Meta-schema runs on decoded source bytes, not the
// IR; [Issue] messages for [CategoryMeta] are navigator-curated, not raw engine
// strings. Each [Issue] also carries its [Issue.DocumentKind] so downstream
// tools can distinguish OpenAPI from Arazzo diagnostics without re-inspecting
// the index. Meta-schema is on by default; set
// [ValidationOptions.SkipMetaSchema] to disable it. Indexes built without
// source bytes (e.g. [NewIndexFromDocument]) skip the meta pass. Tree-sitter
// and standalone parsers produce equivalent structural and schema-shape
// diagnostics.
//
// [Index.Document] may be nil when the source is not an OpenAPI document, and
// [Index.Arazzo] may be nil when the source is not an Arazzo document. Callers
// should check [Index.IsOpenAPI], [Index.IsArazzo], or use [Index.PrimaryValue]
// before accessing family-specific fields.
//
// # Package Layout
//
// The main navigator package contains all types, parsing, indexing, single-file
// and cross-file resolution, workspace infrastructure, and the project facade.
//
// The [github.com/sailpoint-oss/navigator/graph] sub-package adds LSP-grade
// workspace graph management with pipeline processing and snapshot support.
//
// Directory cmd/validation-preview is a CLI that writes Markdown tables of
// [Issue] lines for YAML under testdata/ (PR validation preview in CI).
package navigator
