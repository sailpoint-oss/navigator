import z from "zod";
import { SchemaObjectSchema } from "../schema";
import { BaseSchemaObjectSchema } from "./base";

/**
 * Composition Schema Object for OpenAPI 3.2.
 *
 * Accepts schemas using JSON Schema 2020-12 composition keywords
 * (allOf, anyOf, oneOf, not, if/then/else) without requiring a `type` field.
 */
export const CompositionSchemaObject = BaseSchemaObjectSchema.extend({
	get allOf() {
		return z
			.array(z.lazy(() => SchemaObjectSchema) as z.ZodType)
			.optional()
			.meta({
				title: "allOf",
				description:
					"Require the instance to validate against all subschemas.",
			});
	},
	get anyOf() {
		return z
			.array(z.lazy(() => SchemaObjectSchema) as z.ZodType)
			.optional()
			.meta({
				title: "anyOf",
				description:
					"Require the instance to validate against at least one subschema.",
			});
	},
	get oneOf() {
		return z
			.array(z.lazy(() => SchemaObjectSchema) as z.ZodType)
			.optional()
			.meta({
				title: "oneOf",
				description:
					"Require the instance to validate against exactly one subschema.",
			});
	},
	get not() {
		return z
			.lazy(() => SchemaObjectSchema)
			.optional()
			.meta({
				title: "not",
				description:
					"Require the instance to NOT validate against the given subschema.",
			});
	},
	get ["if"]() {
		return z
			.lazy(() => SchemaObjectSchema)
			.optional()
			.meta({
				title: "if",
				description:
					"Conditional validation: if the instance matches this schema, apply then; otherwise apply else.",
			});
	},
	// biome-ignore lint/suspicious/noThenProperty: then is a valid JSON Schema keyword
	get then() {
		return z
			.lazy(() => SchemaObjectSchema)
			.optional()
			.meta({
				title: "then",
				description: "Subschema applied when if matches.",
			});
	},
	get ["else"]() {
		return z
			.lazy(() => SchemaObjectSchema)
			.optional()
			.meta({
				title: "else",
				description: "Subschema applied when if does not match.",
			});
	},
}).meta({
	title: "Composition Schema",
	description:
		"A schema using composition keywords (allOf, anyOf, oneOf, not, if/then/else) without a type constraint.",
	examples: [
		{
			oneOf: [
				{ $ref: "#/components/schemas/Cat" },
				{ $ref: "#/components/schemas/Dog" },
			],
		},
	],
});

export type CompositionSchema = z.infer<typeof CompositionSchemaObject>;
