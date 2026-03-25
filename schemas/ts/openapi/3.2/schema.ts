import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ArraySchemaObject } from "./data-types/array";
import { BooleanSchemaObject } from "./data-types/boolean";
import { CompositionSchemaObject } from "./data-types/composition";
import { IntegerSchemaObject } from "./data-types/integer";
import { NullSchemaObject } from "./data-types/null";
import { NumberSchemaObject } from "./data-types/number";
import { ObjectSchemaObject } from "./data-types/object";
import { StringObjectSchema } from "./data-types/string";

/**
 * Schema Object union with proper ordering for better error messages.
 *
 * Note: $ref is NOT allowed as a sibling in schema objects. If $ref is present,
 * the object must be a Reference Object (only $ref, summary, description allowed).
 *
 * CompositionSchemaObject is last so typed variants are matched first when a
 * `type` field is present. Typeless composition schemas (allOf, anyOf, etc.)
 * fall through to the composition variant.
 */
export const SchemaObjectSchema = z
	.union([
		ReferenceObjectSchema,
		StringObjectSchema,
		NumberSchemaObject,
		IntegerSchemaObject,
		BooleanSchemaObject,
		NullSchemaObject,
		ArraySchemaObject,
		ObjectSchemaObject,
		CompositionSchemaObject,
	])
	.meta({
		title: "Schema Object",
		description: "A Schema Object defines the shape of a JSON value.",
		examples: [
			{ type: "string", format: "email" },
			{ type: "object", properties: { id: { type: "integer" } } },
		],
	});

export type SchemaObject = z.infer<typeof SchemaObjectSchema>;
