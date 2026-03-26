/**
 * Emit embedded meta JSON Schemas for Go (Draft 2020-12).
 *
 * Zod sources are vendored from Telescope (`schemas/ts/openapi/`, Zod 4).
 */
import { mkdirSync, writeFileSync } from "fs";
import { join } from "path";
import { z } from "zod";
import {
  META_FRAGMENT_ID,
  META_JSON_SCHEMA_URI,
  META_ROOT_ID,
  metaFragmentVersionID,
  metaRootVersionID,
} from "./ids";
import {
  Header2Schema,
  OpenAPI2Schema,
  Operation2Schema,
  Parameter2Schema,
  PathItem2Schema,
  Response2Schema,
  SchemaObject2Schema,
} from "./openapi/openapi-2.0-module";
import {
  Components30Schema,
  Header30Schema,
  OpenAPI30Schema,
  Operation30Schema,
  Parameter30Schema,
  PathItem30Schema,
  RequestBody30Schema,
  Response30Schema,
  SchemaObject30Schema,
  SecurityScheme30Schema,
  Server30Schema,
} from "./openapi/openapi-3.0-module";
import {
  Components31Schema,
  Header31Schema,
  OpenAPI31Schema,
  Operation31Schema,
  Parameter31Schema,
  PathItem31Schema,
  RequestBody31Schema,
  Response31Schema,
  SchemaObject31Schema,
  SecurityScheme31Schema,
  Server31Schema,
} from "./openapi/openapi-3.1-module";
import {
  Components32Schema,
  Header32Schema,
  OpenAPI32Schema,
  Operation32Schema,
  Parameter32Schema,
  PathItem32Schema,
  RequestBody32Schema,
  Response32Schema,
  SchemaObject32Schema,
  SecurityScheme32Schema,
  Server32Schema,
} from "./openapi/openapi-3.2-module";

const here = import.meta.dir;
const metaDir = join(here, "..", "meta");

/** Same post-process as Telescope `server/schemas/scripts/export.ts`. */
function tightenAdditionalProperties(schema: Record<string, unknown>): Record<string, unknown> {
  const defs = (schema.$defs ?? schema.definitions ?? {}) as Record<string, unknown>;

  const emptyDefs = new Set<string>();
  for (const [name, def] of Object.entries(defs)) {
    if (
      def !== null &&
      typeof def === "object" &&
      !Array.isArray(def) &&
      Object.keys(def as Record<string, unknown>).length === 0
    ) {
      emptyDefs.add(name);
    }
  }

  function isEmptySchemaRef(val: unknown): boolean {
    if (val === null || typeof val !== "object" || Array.isArray(val)) return false;
    const obj = val as Record<string, unknown>;
    const ref = obj.$ref;
    if (typeof ref !== "string") return false;
    const match = ref.match(/^#\/\$defs\/(.+)$/);
    if (!match) return false;
    return emptyDefs.has(match[1]);
  }

  function walk(obj: unknown): unknown {
    if (Array.isArray(obj)) return obj.map(walk);
    if (obj !== null && typeof obj === "object") {
      const out: Record<string, unknown> = {};
      for (const [k, v] of Object.entries(obj as Record<string, unknown>)) {
        if (k === "additionalProperties" && (isEmptySchemaRef(v) || isEmptyObj(v))) {
          out[k] = false;
        } else {
          out[k] = walk(v);
        }
      }
      if (looksLikeObjectSchema(out)) {
        const patternProperties =
          out.patternProperties !== null &&
          typeof out.patternProperties === "object" &&
          !Array.isArray(out.patternProperties)
            ? { ...(out.patternProperties as Record<string, unknown>) }
            : {};
        if (patternProperties["^x-"] === undefined) {
          patternProperties["^x-"] = true;
        }
        out.patternProperties = patternProperties;
      }
      return out;
    }
    return obj;
  }

  function isEmptyObj(v: unknown): boolean {
    return (
      v !== null &&
      typeof v === "object" &&
      !Array.isArray(v) &&
      Object.keys(v as Record<string, unknown>).length === 0
    );
  }

  function looksLikeObjectSchema(v: Record<string, unknown>): boolean {
    return (
      v.type === "object" ||
      v.properties !== undefined ||
      v.required !== undefined ||
      v.additionalProperties !== undefined
    );
  }

  // Do not delete empty $defs: other parts of the schema may still $ref them
  // after additionalProperties is tightened (Zod reuse graph).
  return walk(schema) as Record<string, unknown>;
}

function wrapMeta(
  zodJson: Record<string, unknown>,
  extra: {
    $id: string;
    title: string;
    description?: string;
  },
): Record<string, unknown> {
  const { $schema: _drop, ...rest } = zodJson;
  const out: Record<string, unknown> = {
    $schema: META_JSON_SCHEMA_URI,
    $id: extra.$id,
    title: extra.title,
    ...rest,
  };
  if (extra.description !== undefined) {
    out.description = extra.description;
  }
  return out;
}

function writeSchema(filename: string, body: Record<string, unknown>) {
  const path = join(metaDir, filename);
  writeFileSync(path, `${JSON.stringify(body, null, 2)}\n`, "utf-8");
}

function writeNamedSchema(
  filename: string,
  schema: z.ZodType,
  meta: {
    $id: string;
    title: string;
    description?: string;
  },
) {
  writeSchema(filename, wrapMeta(toDraft2020(schema), meta));
}

function toDraft2020(zodSchema: z.ZodType): Record<string, unknown> {
  const jsonSchema = z.toJSONSchema(zodSchema, {
    target: "draft-2020-12",
    reused: "ref",
    unrepresentable: "any",
  });
  return tightenAdditionalProperties(jsonSchema as Record<string, unknown>);
}

mkdirSync(metaDir, { recursive: true });

export const navigatorRootMetaSchema = z.union([
  OpenAPI30Schema,
  OpenAPI31Schema,
  OpenAPI32Schema,
  OpenAPI2Schema,
]);

export const navigatorFragmentMetaSchema = z.union([
  PathItem30Schema,
  Operation30Schema,
  Parameter30Schema,
  RequestBody30Schema,
  Response30Schema,
  Header30Schema,
  SecurityScheme30Schema,
  Components30Schema,
  SchemaObject30Schema,
  Server30Schema,
  PathItem31Schema,
  Operation31Schema,
  Parameter31Schema,
  RequestBody31Schema,
  Response31Schema,
  Header31Schema,
  SecurityScheme31Schema,
  Components31Schema,
  SchemaObject31Schema,
  Server31Schema,
  PathItem32Schema,
  Operation32Schema,
  Parameter32Schema,
  RequestBody32Schema,
  Response32Schema,
  Header32Schema,
  SecurityScheme32Schema,
  Components32Schema,
  SchemaObject32Schema,
  Server32Schema,
  SchemaObject2Schema,
]);

writeNamedSchema("openapi-root.schema.json", navigatorRootMetaSchema, {
  $id: META_ROOT_ID,
  title: "Navigator OpenAPI root document (meta-validation)",
});

writeNamedSchema("openapi-fragment.schema.json", navigatorFragmentMetaSchema, {
  $id: META_FRAGMENT_ID,
  title: "Navigator OpenAPI fragment document (meta-validation)",
  description:
    "Non-root files: path items, operations, schema objects, components snippets, etc. (Telescope Zod sources).",
});

const versionedRootSchemas = [
  {
    version: "2.0",
    filename: "openapi-2.0-root.schema.json",
    schema: OpenAPI2Schema,
    title: "Navigator OpenAPI 2.0 root document (meta-validation)",
  },
  {
    version: "3.0",
    filename: "openapi-3.0-root.schema.json",
    schema: OpenAPI30Schema,
    title: "Navigator OpenAPI 3.0 root document (meta-validation)",
  },
  {
    version: "3.1",
    filename: "openapi-3.1-root.schema.json",
    schema: OpenAPI31Schema,
    title: "Navigator OpenAPI 3.1 root document (meta-validation)",
  },
  {
    version: "3.2",
    filename: "openapi-3.2-root.schema.json",
    schema: OpenAPI32Schema,
    title: "Navigator OpenAPI 3.2 root document (meta-validation)",
  },
] as const;

for (const item of versionedRootSchemas) {
  writeNamedSchema(item.filename, item.schema, {
    $id: metaRootVersionID(item.version),
    title: item.title,
  });
}

const versionedFragmentSchemas = [
  { version: "2.0", kind: "schema", schema: SchemaObject2Schema },
  { version: "2.0", kind: "path-item", schema: PathItem2Schema },
  { version: "2.0", kind: "operation", schema: Operation2Schema },
  { version: "2.0", kind: "parameter", schema: Parameter2Schema },
  { version: "2.0", kind: "response", schema: Response2Schema },
  { version: "2.0", kind: "header", schema: Header2Schema },
  { version: "3.0", kind: "schema", schema: SchemaObject30Schema },
  { version: "3.0", kind: "path-item", schema: PathItem30Schema },
  { version: "3.0", kind: "operation", schema: Operation30Schema },
  { version: "3.0", kind: "parameter", schema: Parameter30Schema },
  { version: "3.0", kind: "request-body", schema: RequestBody30Schema },
  { version: "3.0", kind: "response", schema: Response30Schema },
  { version: "3.0", kind: "header", schema: Header30Schema },
  { version: "3.0", kind: "security-scheme", schema: SecurityScheme30Schema },
  { version: "3.0", kind: "components", schema: Components30Schema },
  { version: "3.0", kind: "server", schema: Server30Schema },
  { version: "3.1", kind: "schema", schema: SchemaObject31Schema },
  { version: "3.1", kind: "path-item", schema: PathItem31Schema },
  { version: "3.1", kind: "operation", schema: Operation31Schema },
  { version: "3.1", kind: "parameter", schema: Parameter31Schema },
  { version: "3.1", kind: "request-body", schema: RequestBody31Schema },
  { version: "3.1", kind: "response", schema: Response31Schema },
  { version: "3.1", kind: "header", schema: Header31Schema },
  { version: "3.1", kind: "security-scheme", schema: SecurityScheme31Schema },
  { version: "3.1", kind: "components", schema: Components31Schema },
  { version: "3.1", kind: "server", schema: Server31Schema },
  { version: "3.2", kind: "schema", schema: SchemaObject32Schema },
  { version: "3.2", kind: "path-item", schema: PathItem32Schema },
  { version: "3.2", kind: "operation", schema: Operation32Schema },
  { version: "3.2", kind: "parameter", schema: Parameter32Schema },
  { version: "3.2", kind: "request-body", schema: RequestBody32Schema },
  { version: "3.2", kind: "response", schema: Response32Schema },
  { version: "3.2", kind: "header", schema: Header32Schema },
  { version: "3.2", kind: "security-scheme", schema: SecurityScheme32Schema },
  { version: "3.2", kind: "components", schema: Components32Schema },
  { version: "3.2", kind: "server", schema: Server32Schema },
] as const;

for (const item of versionedFragmentSchemas) {
  writeNamedSchema(`openapi-${item.version}-${item.kind}.schema.json`, item.schema, {
    $id: metaFragmentVersionID(item.version, item.kind),
    title: `Navigator OpenAPI ${item.version} ${item.kind} fragment (meta-validation)`,
  });
}

console.log("Wrote schemas/meta/openapi-root.schema.json");
console.log("Wrote schemas/meta/openapi-fragment.schema.json");
