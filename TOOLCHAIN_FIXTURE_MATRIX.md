# Toolchain Fixture Matrix

This document defines the shared fixture, golden, and parity matrix for the
OpenAPI toolchain.

The matrix is intentionally **logical**, not a single copied directory. These
repositories are versioned independently, so each repo keeps standalone fixtures
and tests, but they should line up to the same fixture IDs, coverage tiers, and
assertion intent.

## Matrix Rules

1. Every new cross-repo fixture starts with a stable logical ID here.
2. Repos can use local equivalents instead of byte-identical copies when they
   need different formats:
   - tree-sitter corpus files keep CST expectations.
   - Navigator keeps parse/index/validation fixtures.
   - Barrelman keeps lint-focused file fixtures.
   - Telescope keeps SDK/LSP fixtures and mirrored sidecar copies.
   - Cartographer keeps extraction and spec-loader fixtures.
   - Barometer keeps runtime validation fixtures.
3. If a fixture is relevant to a repo's capability, that repo should have a
   local test for it or an explicitly documented gap.

## Capability Tags

- `cst`: grammar and parse-tree coverage
- `parse`: typed OpenAPI indexing and version detection
- `validate`: parse-time structural, schema, or meta-schema validation
- `workspace`: multi-file `$ref` topology and resolver behavior
- `lint`: Barrelman and Telescope diagnostic parity
- `editor`: ranges, pointers, code actions, rename, references
- `extract`: extraction and downstream loader smoke coverage
- `runtime`: live request/response contract execution

## Fixture IDs

### `root-swagger20-minimal`

- Tags: `cst`, `parse`
- tree-sitter-openapi: `tree-sitter-openapi/test/corpus/yaml_versions.txt`
- navigator: `testdata/versions/swagger-2.0.yaml` via `fixture_matrix_test.go`
- telescope: `server/testutil/specs/api-v2.yaml`
- barometer: `testdata/openapi/swagger-2.0.yaml` via `internal/openapi/loader_versions_test.go` (explicit rejection coverage)

### `root-oas30-minimal`

- Tags: `cst`, `parse`, `validate`, `lint`, `editor`, `runtime`
- tree-sitter-openapi: `tree-sitter-openapi/test/corpus/yaml_basic.txt`, `tree-sitter-openapi/test/corpus/yaml_versions.txt`
- navigator: `testdata/versions/openapi-3.0.yaml`, `testdata/validate/ok_minimal.yaml`
- barrelman: `testdata/toolchain/oas30-minimal.yaml`
- telescope: `server/testutil/specs/openapi-3.0.yaml`, `server/lsp/toolchain_parity_test.go`
- barometer: `testdata/openapi/petstore-minimal.yaml`

### `root-oas31-minimal`

- Tags: `cst`, `parse`, `validate`, `editor`, `runtime`
- tree-sitter-openapi: `tree-sitter-openapi/test/corpus/yaml_versions.txt`
- navigator: `testdata/versions/openapi-3.1.yaml`
- telescope: `server/testutil/specs/openapi-3.1.yaml`
- barometer: `testdata/openapi/petstore-3.1.yaml`

### `root-oas32-minimal`

- Tags: `cst`, `parse`, `validate`, `editor`, `runtime`
- tree-sitter-openapi: `tree-sitter-openapi/test/corpus/yaml_versions.txt`
- navigator: `testdata/versions/openapi-3.2.yaml`
- telescope: `server/testutil/specs/openapi-3.2.yaml`
- barometer: `testdata/openapi/petstore-3.2.yaml`

### `root-invalid-missing-info`

- Tags: `validate`, `lint`
- navigator: `testdata/meta/invalid_root_no_info.yaml`
- telescope: `server/testutil/specs/invalid-openapi-structural.yaml`, `server/lsp/toolchain_parity_test.go`

### `fragment-schema-valid`

- Tags: `cst`, `parse`, `validate`
- tree-sitter-openapi: `tree-sitter-openapi/test/corpus/yaml_fragment_schema.txt`
- navigator: `testdata/meta/ok_fragment_schema.yaml`, `testdata/fragments/schema.yaml`

### `fragment-schema-invalid`

- Tags: `parse`, `validate`
- navigator: `testdata/meta/invalid_fragment_schema.yaml` via `fixture_matrix_test.go`

### `workspace-multifile`

- Tags: `parse`, `workspace`, `lint`, `extract`
- navigator: `testdata/multifile/v1/`
- barrelman: `testdata/toolchain/multifile/`
- telescope: `server/testutil/specs/test-multi-file.yaml`, `test-files/openapi/test-multi-file-refs/`
- cartographer: `cartographer/spec/testdata/multifile/`

### `workspace-cycle`

- Tags: `workspace`, `extract`
- navigator: `testdata/cycles/`
- telescope: `test-files/openapi/test-ref-cycle-main.yaml`, `test-files/openapi/test-ref-cycle-b.yaml`
- cartographer: `cartographer/spec/testdata/circular/`

### `workspace-diamond`

- Tags: `workspace`
- navigator: `testdata/diamond/`

### `workspace-missing-target`

- Tags: `workspace`, `lint`, `editor`
- navigator: `resolver_test.go` (`TestCrossFileResolver_MissingFile`)
- barrelman: `testdata/toolchain/multifile-broken/`
- telescope: `server/lsp/handler_test.go` unresolved `$ref` quick-fix coverage

### `editor-kitchen-sink`

- Tags: `cst`, `parse`, `editor`
- tree-sitter-openapi: `tree-sitter-openapi/test/corpus/yaml_ref_and_extensions.txt`, `yaml_description_markdown.txt`
- navigator: `testdata/golden/kitchen-sink/specs/api.yaml`, `parse_parity_test.go`
- telescope: `server/lsp/handler_test.go`, `server/lsp/toolchain_parity_test.go`

### `extract-loader-smoke`

- Tags: `extract`
- cartographer: `cartographer/spec/testdata/singlefile/`, `multifile/`, `crossformat/`, `circular/`

### `runtime-loader-smoke`

- Tags: `runtime`
- barometer: `testdata/openapi/petstore-minimal.yaml`, `testdata/openapi/petstore-with-header.yaml`

## Current Test Anchors

- Navigator:
  `parse_parity_test.go`, `fixture_matrix_test.go`, `project_test.go`,
  `resolver_test.go`, `meta_schema_test.go`
- Barrelman:
  `fixture_matrix_test.go`, `lintfiles_workspace_test.go`
- Telescope:
  `server/testutil/specs/mirror_sync_test.go`, `server/testutil/golden/pipeline_test.go`,
  `server/lsp/toolchain_parity_test.go`, `server/lsp/handler_test.go`
- Cartographer:
  `cartographer/spec/toolchain_fixture_test.go`
- Barometer:
  `internal/openapi/toolchain_fixture_test.go`, `internal/openapi/e2e_test.go`

## Adding A New Fixture

1. Add the logical ID and capability tags here.
2. Add or update the tree-sitter corpus case if CST shape is part of the
   contract.
3. Add or update the Navigator fixture for parse/index/validation ownership.
4. Add or update Barrelman and Telescope parity tests if the fixture should
   yield stable diagnostics.
5. Add or update Cartographer or Barometer smoke coverage when the fixture is
   relevant to extraction or runtime validation.
