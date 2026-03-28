# OpenAPI Toolchain Boundaries

This document defines the intended contract between the API-description toolchain repositories in this workspace. It is the short reference for "who owns what" across parsing, validation, linting, editor UX, extraction, orchestration, and runtime testing for OpenAPI and Arazzo.

## Repository ownership

| Repo | Owns | Does not own |
|------|------|--------------|
| `tree-sitter-openapi` | YAML/JSON grammar, CST fidelity, injections, and bindings for OpenAPI and Arazzo documents | Typed IR, schema/meta validation, rule execution, orchestration |
| `navigator` | Parse/load OpenAPI and Arazzo, typed IR, pointers/ranges, `$ref` indexes and multi-file project/resolver APIs where applicable, canonical parse-time issues, schema/meta validation | Style/security linting, editor UX, source-code extraction, runtime execution |
| `barrelman` | Static rule execution, rulesets, severities, filtering, and diagnostic packaging on top of Navigator indexes | Parser ownership, canonical schema validation, runtime execution, editor UX |
| `cartographer` | Service-local source-code-to-OpenAPI extraction, `.cartographer/cartographer.yaml`, extraction CLI/library/API | Spec-side lint/validate UX, org orchestration, runtime execution |
| `telescope` | Spec-side CLI/LSP/editor UX over workspace files, diagnostics aggregation, bundle preview, code actions, hover/definition/reference flows, custom rule runtimes | Canonical parser ownership, source-code extraction, runtime engine ownership |
| `meridian` | Codebase-side orchestration, `analyze` and `pipeline`, org-level config, report aggregation, optional runtime stage composition | Extraction engine implementation, spec-side LSP UX, runtime engine ownership |
| `barometer` | Runtime contract execution, response validation, Arazzo workflow execution, reusable runtime report contracts, CLI/package/action surfaces | Static parsing/lint ownership, editor UX, codebase extraction/orchestration ownership |

## Validation layering

- `navigator` is the canonical source of parse-time diagnostics for OpenAPI and Arazzo: syntax, structural shape, schema-shape where applicable, and meta-schema validation.
- `barrelman` maps Navigator issues into diagnostics and layers semantic/style/security/cross-file rules on top. Raw parse/schema correctness stays in Navigator.
- `telescope` presents those diagnostics in spec-side CLI, SDK, and editor flows. It should not invent a second meaning for "lint" or "validate".
- `meridian` composes extraction plus static analysis and reporting for codebase-facing CI/CD, but it should still consume shared Navigator and Barrelman contracts instead of recreating them.
- Runtime tools should consume the shared static contracts; they should not maintain parallel OpenAPI or Arazzo parsing stacks.

## Pipeline ownership

- `navigator/graph` owns the shared workspace substrate and the built-in `raw -> parse -> bind` stages.
- `barrelman` owns static lint-rule execution on top of indexed documents.
- `telescope` owns the higher-level spec-side user experience that combines parse/bind, linting, validation presentation, and analysis workflows over workspace files.
- `cartographer` owns source-code extraction into a single service-local OpenAPI result.
- `meridian` owns higher-level codebase and multi-repo orchestration after extraction, including report generation and optional runtime-stage composition.
- `barometer` owns runtime execution once static documents or extracted specs are available.
- Stage names such as `lint`, `validate`, `analyze`, and `test` may exist as shared vocabulary, but each repo should keep its ownership boundary clear.

## Public contracts

- Grammar and CST consumers should depend on `tree-sitter-openapi` directly.
- Parse/index/ref consumers should depend on `navigator`.
- Rule authors and lint runners should depend on `barrelman`.
- Source-code extraction consumers should depend on `cartographer/extraction`.
- Editor integrations should depend on `telescope`, which in turn consumes `navigator` and `barrelman`.
- Codebase orchestration should depend on `meridian`, which may consume `cartographer/extraction`, `navigator`, `barrelman`, and `barometer`.
- Downstream generators, runtime tools, and orchestration layers should prefer `navigator.Index`, Navigator issues, Barrelman diagnostics, and Barometer's versioned runtime report schema as their shared compatibility surfaces.

## Compatibility and release order

- Repositories are versioned independently. The current pinned compatibility state lives in each consumer's `go.mod`, especially:
  - `barrelman/go.mod`
  - `meridian/go.mod`
  - `telescope/server/go.mod`
  - `barometer/go.mod`
- If `navigator` changes public parse/index/range/resolver contracts, update consumers in this order:
  1. `navigator`
  2. `barrelman`
  3. `telescope/server`
  4. `barometer`
  5. `meridian`
- If `cartographer/extraction` changes its public extraction/config contract, update consumers in this order:
  1. `cartographer`
  2. `meridian`
- If `barrelman` changes rule metadata, diagnostics, or config contracts, update consumers in this order:
  1. `barrelman`
  2. `telescope/server`
  3. `meridian`
- If `barometer` changes runtime report or action contracts, update consumers in this order:
  1. `barometer`
  2. `telescope`
  3. `meridian`
- Keep local development aligned with a workspace `go.work` file instead of temporary `replace` directives whenever multiple sibling repos are changing together.
- Use `TOOLCHAIN_FIXTURE_MATRIX.md` as the source of truth for cross-repo fixture IDs, parity anchors, and smoke-test intent before cutting releases.

## Coordinated smoke checks

- Navigator core: `go test ./...`
- Cartographer: `cd cartographer && go test ./...`
- Barrelman: `go test ./...`
- Telescope server: `go test ./...`
- Barometer: `go test ./internal/openapi ./...`
- Meridian: `GOWORK=off go test ./...`
- GitHub Actions: run `.github/workflows/downstream-smoke.yml` in `navigator` before publishing compatibility-sensitive changes.

## Practical guidance

- If you need a CST, grammar, or syntax highlighting, go to `tree-sitter-openapi`.
- If you need an OpenAPI `Document`, an `ArazzoDocument`, an `Index`, a `Project`, precise ranges, or canonical parse/schema validation, go to `navigator`.
- If you need static rule execution, severity filtering, or rule metadata, go to `barrelman`.
- If you need service-local source extraction, go to `cartographer`.
- If you need hover, definition, references, rename, code actions, or spec-side CLI/LSP UX, go to `telescope`.
- If you need codebase or org-level orchestration after extraction, go to `meridian`.
- If you need live HTTP contract execution or Arazzo workflow execution, go to `barometer`.
