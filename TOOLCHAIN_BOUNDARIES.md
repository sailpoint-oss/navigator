# OpenAPI Toolchain Boundaries

This document defines the intended contract between the OpenAPI toolchain repositories in this workspace. It is the short reference for "who owns what" across parsing, validation, linting, editor UX, extraction, and runtime testing.

## Repository ownership

| Repo | Owns | Does not own |
|------|------|--------------|
| `tree-sitter-openapi` | YAML/JSON grammar, CST fidelity, injections, bindings | OpenAPI IR, semantic validation, rule execution |
| `navigator` | Parse, typed IR, pointers/ranges, `$ref` indexes, multi-file project/resolver APIs, parse-time issues | Style/security linting, editor UX, org-level reporting |
| `barrelman` | Rules, rulesets, severities, diagnostic packaging, lint execution | Parsing, workspace graph primitives, OpenAPI IR ownership |
| `telescope` | LSP, CLI/SDK UX, workspace lifecycle, code actions, hover/definition/reference flows, user-visible orchestration | Canonical parser ownership, long-term duplicate structural validators |
| `cartographer` | Extraction, rewrite, bundle, compare, audit, report generation, multi-repo pipelines | Core parser or linter ownership |
| `barometer` | Runtime contract execution, response validation, Arazzo workflows | Static parsing/lint ownership |

## Validation layering

- `navigator` is the canonical source of parse-time diagnostics: syntax, structural shape, schema-shape, and meta-schema validation.
- `barrelman` maps navigator issues into diagnostics and layers semantic/style/security/cross-file rules on top.
- `telescope` presents those diagnostics in editor and SDK flows. It should not invent a second meaning for "lint" or "validate".
- `cartographer` and `barometer` consume the shared contracts; they should not maintain parallel OpenAPI parsing stacks.

## Pipeline ownership

- `navigator/graph` owns the shared workspace substrate and the built-in `raw -> parse -> bind` stages.
- `barrelman` owns lint-rule execution on top of indexed documents.
- `telescope` owns the higher-level user experience that combines parse/bind, linting, validation presentation, and analysis workflows.
- Stage names such as `lint`, `validate`, and `analyze` may exist as shared vocabulary, but Navigator's default pipeline does not currently ship built-in implementations for them.

## Public contracts

- Grammar and CST consumers should depend on `tree-sitter-openapi` directly.
- Parse/index/ref consumers should depend on `navigator`.
- Rule authors and lint runners should depend on `barrelman`.
- Editor integrations should depend on `telescope`, which in turn consumes `navigator` and `barrelman`.
- Downstream generators and runtime tools should prefer `navigator.Index` and `barrelman` diagnostics as their shared compatibility surface.

## Compatibility and release order

- Repositories are versioned independently. The current pinned compatibility state lives in each consumer's `go.mod`, especially:
  - `barrelman/go.mod`
  - `telescope/server/go.mod`
  - `cartographer/cartographer/go.mod`
  - `barometer/go.mod`
- If `navigator` changes public parse/index/range/resolver contracts, update consumers in this order:
  1. `navigator`
  2. `barrelman`
  3. `telescope/server`
  4. `cartographer/cartographer`
  5. `barometer`
- If `barrelman` changes rule metadata, diagnostics, or config contracts, update consumers in this order:
  1. `barrelman`
  2. `telescope/server`
  3. `cartographer/cartographer`
- Keep local development aligned with a workspace `go.work` file instead of temporary `replace` directives whenever multiple sibling repos are changing together.
- Use `TOOLCHAIN_FIXTURE_MATRIX.md` as the source of truth for cross-repo fixture IDs, parity anchors, and smoke-test intent before cutting releases.

## Coordinated smoke checks

- Navigator core: `go test ./...`
- Barrelman: `go test ./...`
- Telescope server: `go test ./...`
- Cartographer: `cd cartographer && go test ./cmd ./spec ./unified ./telescope`
- Barometer: `go test ./internal/openapi ./...`
- GitHub Actions: run `.github/workflows/downstream-smoke.yml` in `navigator` before publishing compatibility-sensitive changes.

## Practical guidance

- If you need a CST or syntax highlighting, go to `tree-sitter-openapi`.
- If you need a `Document`, `Index`, `Project`, `Range`, `PointerIndex`, or cross-file `$ref` resolution, go to `navigator`.
- If you need a rule registry, `LintContent`, `LintFiles`, severity filtering, or rule metadata, go to `barrelman`.
- If you need hover, definition, references, rename, code actions, or LSP/CLI orchestration, go to `telescope`.
- If you need extraction, bundling, compare/audit, or report generation, go to `cartographer`.
- If you need live HTTP contract execution or Arazzo workflows, go to `barometer`.
