# Navigator

[![Go Reference](https://pkg.go.dev/badge/github.com/sailpoint-oss/navigator.svg)](https://pkg.go.dev/github.com/sailpoint-oss/navigator)
[![CI](https://github.com/sailpoint-oss/navigator/actions/workflows/ci.yml/badge.svg)](https://github.com/sailpoint-oss/navigator/actions/workflows/ci.yml)

Uniform OpenAPI file handling for the workspace toolchain. Navigator owns parsing, typed IR construction, pointers/ranges, multi-file `$ref` resolution, and parse-time validation for Swagger 2.0 and OpenAPI 3.0--3.2 documents.

## Toolchain Role

Navigator is the canonical parse/index/ref layer in this workspace:

- `tree-sitter-openapi` owns grammar and CST quality.
- `navigator` owns parse/index/ref/validation contracts.
- `barrelman` owns semantic/style/security rules and diagnostic packaging.
- `telescope` owns editor and SDK experiences built on top of Navigator + Barrelman.

See [`TOOLCHAIN_BOUNDARIES.md`](TOOLCHAIN_BOUNDARIES.md) for the cross-repo ownership model.

## Installation

```bash
go get github.com/sailpoint-oss/navigator
```

## Quick Start

### Parse a Single File

```go
import navigator "github.com/sailpoint-oss/navigator"

idx := navigator.Parse(content)
if idx == nil {
    log.Fatal("failed to parse")
}

for _, op := range idx.AllOperations() {
    fmt.Printf("%s %s  operationId=%s\n", op.Method, op.Path, op.Operation.OperationID)
}
```

### Open a Multi-File Project

```go
project, err := navigator.OpenProject("path/to/root.yaml")
if err != nil {
    log.Fatal(err)
}
project.LoadAll() // eagerly load all transitive $ref targets

result, err := project.Resolve(project.RootURI(), "./schemas/Pet.yaml")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Resolved to %s (pointer: %s)\n", result.TargetURI, result.Pointer)
```

### Full Workspace Graph (Shared Substrate)

```go
import (
    navigator "github.com/sailpoint-oss/navigator"
    navgraph  "github.com/sailpoint-oss/navigator/graph"
)

ws := navigator.NewWorkspace()
g := navgraph.New(ws.Cache, ws.Graph)
g.AddSource(navigator.NewMemorySource(uri, content))

runner := navgraph.NewPipelineRunner(navgraph.DefaultStages()...) // raw -> parse -> bind
runner.RunAll(ctx, g, uri)
```

`navigator/graph` intentionally stops at the shared parse/bind substrate. Lint-rule execution and higher-level validate/analyze workflows are downstream concerns.

## Design Principles

**Lazy by default** -- files are loaded and parsed only when needed. `OpenProject` loads the root document; transitive `$ref` targets load on demand. `Project.LoadAll()` is the explicit break-glass to force-load everything.

**Break-glass at every level** -- every abstraction exposes its internals. An `Index` provides `Tree()` for the raw tree-sitter CST, `SemanticRoot()` for the IR, and `Document` for the typed model. A `Project` exposes its `Cache()`, `Graph()`, and `Resolver()` for direct manipulation.

**Workspace-wide shared state** -- `IndexCache` and `FileGraph` are singletons shared across multiple roots. A `Project` is a thin lens over this shared state.

## Package Layout

| Package | Import Path                             | Description                                                                                                            |
| ------- | --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| Core    | `github.com/sailpoint-oss/navigator`       | All types, parsing, indexing, single-file and cross-file resolution, workspace infrastructure, and the project facade. |
| Graph   | `github.com/sailpoint-oss/navigator/graph` | Shared workspace graph substrate with built-in `raw`, `parse`, and `bind` stages used by downstream orchestrators.   |

### Core Package Highlights

| Area           | Key Types / Functions                                                                |
| -------------- | ------------------------------------------------------------------------------------ |
| Parsing        | `Parse`, `ParseURI` (standalone, no CGO) / `ParseContent`, `ParseTree` (tree-sitter) |
| Document Model | `Document`, `PathItem`, `Operation`, `Schema`, `Components`, `Parameter`, `Response` |
| Indexing       | `Index`, `PointerIndex`, `IndexCache`                                                |
| Resolution     | `ResolveRef` (local), `CrossFileResolver` (cross-file)                               |
| Workspace      | `Project`, `Workspace`, `FileGraph`, `Discovery`                                     |
| Sources        | `DocumentSource`, `FilesystemSource`, `MemorySource`                                 |
| Locations      | `Loc`, `Range`, `Position`, `SemanticNode`                                           |

## Tree-sitter vs Standalone Parsing

Navigator provides two parsing paths:

|                          | Standalone (`Parse` / `ParseURI`) | Tree-sitter (`ParseContent` / `ParseTree`) |
| ------------------------ | --------------------------------- | ------------------------------------------ |
| **CGO required**         | No                                | Yes                                        |
| **Source locations**     | Approximate (yaml.v3 line/column) | Precise (tree-sitter byte offsets)         |
| **`idx.Tree()`**         | `nil`                             | `*tree_sitter.Tree`                        |
| **`idx.SemanticRoot()`** | `nil`                             | `*SemanticNode`                            |
| **Best for**             | CLI tools, CI, contract testing   | LSP, linting with precise diagnostics      |

Both parsers produce logically equivalent `*Index` output. The same operations, schemas, refs, and document structure are present regardless of parser.

## Dependencies

| Dependency                                                               | Purpose                                   |
| ------------------------------------------------------------------------ | ----------------------------------------- |
| [tree-sitter-openapi](https://github.com/sailpoint-oss/tree-sitter-openapi) | Tree-sitter grammar for OpenAPI YAML/JSON |
| [go-tree-sitter](https://github.com/tree-sitter/go-tree-sitter)          | Tree-sitter Go runtime                    |
| [fsnotify](https://github.com/fsnotify/fsnotify)                         | Filesystem watching                       |
| [yaml.v3](https://gopkg.in/yaml.v3)                                      | Standalone YAML parsing (no CGO)          |

## Releasing

Navigator uses Go module versioning, but the repository's current release workflow is automated:

- `.github/workflows/release.yml` runs on pushes to `main`.
- It bumps the next patch tag automatically and creates the GitHub Release.
- For compatibility-sensitive changes, run the `Downstream Smoke` workflow first and verify the anchors listed in `TOOLCHAIN_FIXTURE_MATRIX.md`.

Release coordination guidance lives in:

- `TOOLCHAIN_BOUNDARIES.md` for ownership, bump order, and smoke checks
- `TOOLCHAIN_FIXTURE_MATRIX.md` for the cross-repo parity matrix

The [Go module proxy](https://proxy.golang.org/) automatically indexes published tags after the release workflow pushes them.

## License

See [LICENSE](LICENSE) for details.
