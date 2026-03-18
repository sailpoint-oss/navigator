# Navigator

[![Go Reference](https://pkg.go.dev/badge/github.com/LukasParke/navigator.svg)](https://pkg.go.dev/github.com/LukasParke/navigator)
[![CI](https://github.com/LukasParke/navigator/actions/workflows/ci.yml/badge.svg)](https://github.com/LukasParke/navigator/actions/workflows/ci.yml)

Uniform OpenAPI file handling for the entire workspace toolchain. Navigator parses, indexes, and resolves `$ref`s across single-file and massively multi-file OpenAPI specifications (Swagger 2.0, OpenAPI 3.0--3.2).

## Installation

```bash
go get github.com/LukasParke/navigator
```

## Quick Start

### Parse a Single File

```go
import navigator "github.com/LukasParke/navigator"

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

### Full Workspace Graph (LSP Pattern)

```go
import (
    navigator "github.com/LukasParke/navigator"
    navgraph  "github.com/LukasParke/navigator/graph"
)

ws := navigator.NewWorkspace()
g := navgraph.New(ws.Cache, ws.Graph)
g.AddSource(navigator.NewMemorySource(uri, content))

runner := navgraph.NewPipelineRunner(navgraph.DefaultStages()...)
runner.RunAll(ctx, g, uri)
```

## Design Principles

**Lazy by default** -- files are loaded and parsed only when needed. `OpenProject` loads the root document; transitive `$ref` targets load on demand. `Project.LoadAll()` is the explicit break-glass to force-load everything.

**Break-glass at every level** -- every abstraction exposes its internals. An `Index` provides `Tree()` for the raw tree-sitter CST, `SemanticRoot()` for the IR, and `Document` for the typed model. A `Project` exposes its `Cache()`, `Graph()`, and `Resolver()` for direct manipulation.

**Workspace-wide shared state** -- `IndexCache` and `FileGraph` are singletons shared across multiple roots. A `Project` is a thin lens over this shared state.

## Package Layout

| Package | Import Path                             | Description                                                                                                            |
| ------- | --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| Core    | `github.com/LukasParke/navigator`       | All types, parsing, indexing, single-file and cross-file resolution, workspace infrastructure, and the project facade. |
| Graph   | `github.com/LukasParke/navigator/graph` | LSP-grade workspace graph management with pipeline processing and snapshot support.                                    |

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
| [tree-sitter-openapi](https://github.com/LukasParke/tree-sitter-openapi) | Tree-sitter grammar for OpenAPI YAML/JSON |
| [go-tree-sitter](https://github.com/tree-sitter/go-tree-sitter)          | Tree-sitter Go runtime                    |
| [fsnotify](https://github.com/fsnotify/fsnotify)                         | Filesystem watching                       |
| [yaml.v3](https://gopkg.in/yaml.v3)                                      | Standalone YAML parsing (no CGO)          |

## Releasing

Navigator uses Go module versioning. To publish a new version:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The [Go module proxy](https://proxy.golang.org/) automatically indexes tagged versions. A GitHub Release with auto-generated notes is created by CI on every `v*` tag push.

## License

See [LICENSE](LICENSE) for details.
