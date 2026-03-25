import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ArraySchemaObject } from "./data-types/array";
import { BaseSchemaObjectSchema } from "./data-types/base";
import { BooleanSchemaObject } from "./data-types/boolean";
import { CompositionSchemaObject } from "./data-types/composition";
import { IntegerSchemaObject } from "./data-types/integer";
import { NullSchemaObject } from "./data-types/null";
import { NumberSchemaObject } from "./data-types/number";
import { ObjectSchemaObject } from "./data-types/object";
import { StringObjectSchema } from "./data-types/string";

const TypeArraySchemaObject = BaseSchemaObjectSchema.extend({
	type: z
		.array(
			z.enum([
				"string",
				"number",
				"integer",
				"boolean",
				"null",
				"array",
				"object",
			]),
		)
		.min(1)
		.meta({ title: "type" }),
}).meta({
	title: "Type Array Schema",
	description:
		"OpenAPI 3.1 (JSON Schema) allows type to be an array of primitive types.",
	examples: [{ type: ["string", "null"] }],
});

/**
 * Schema Object union with proper ordering for better error messages.
 *
 *
 * Note: $ref is NOT allowed as a sibling in schema objects. If $ref is present,
 * the object must be a Reference Object (only $ref, summary, description allowed).
 */
export const SchemaObjectSchema = z
	.union([
		ReferenceObjectSchema,
		TypeArraySchemaObject,
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
