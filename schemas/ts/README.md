# Meta-schema sources (Zod → JSON Schema)

Embedded OpenAPI meta-schemas under `schemas/meta/*.json` are **generated** by [Bun](https://bun.sh/) from **Zod 4** sources.

The large tree under [`openapi/`](openapi/) is vendored from [Telescope](https://github.com/sailpoint-oss/telescope) (`server/schemas/src`). [`build.ts`](build.ts) composes Telescope’s root modules into the navigator **root** meta-schema (OpenAPI 2.0 + 3.0 + 3.1 + 3.2 union) and Telescope **fragment** exports into the **fragment** meta-schema. To refresh the vendored copy after upstream changes:

```bash
rsync -a --delete ../telescope/server/schemas/src/ ./schemas/ts/openapi/
```

- **Build** (repository root, next to `go.mod`):

  ```bash
  bun install
  bun run schemas:build
  ```

- **Regenerate from Go:** `go generate ./...` (Bun on `PATH`).

Go loads these schemas with **`jsonschema/v6`** using **`regexp2` in ECMAScript mode** (see `meta_ecma_regexp.go`) so JSON Schema `pattern` keywords match Draft 2020-12 / ECMA-262, including patterns Zod emits for `.email()` and `.url()`.

Commit **both** `schemas/ts/**` and `schemas/meta/*.json` so `go test` works without Bun.
