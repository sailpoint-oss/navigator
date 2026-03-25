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
} from "./ids";
import { OpenAPI2Schema, SchemaObject2Schema } from "./openapi/openapi-2.0-module";
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

const rootJson = toDraft2020(navigatorRootMetaSchema);
writeSchema(
  "openapi-root.schema.json",
  wrapMeta(rootJson, {
    $id: META_ROOT_ID,
    title: "Navigator OpenAPI root document (meta-validation)",
  }),
);

const fragJson = toDraft2020(navigatorFragmentMetaSchema);
writeSchema(
  "openapi-fragment.schema.json",
  wrapMeta(fragJson, {
    $id: META_FRAGMENT_ID,
    title: "Navigator OpenAPI fragment document (meta-validation)",
    description:
      "Non-root files: path items, operations, schema objects, components snippets, etc. (Telescope Zod sources).",
  }),
);

console.log("Wrote schemas/meta/openapi-root.schema.json");
console.log("Wrote schemas/meta/openapi-fragment.schema.json");
