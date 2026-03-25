import z from "zod";
import { SchemaObjectSchema } from "../schema";
import { BaseSchemaObjectSchema } from "./base";

/**
 * Composition Schema Object for OpenAPI 3.0.
 *
 * Accepts schemas using composition keywords (allOf, anyOf, oneOf, not)
 * without requiring a `type` field. OpenAPI 3.0 supports these four
 * composition keywords but not if/then/else.
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
}).meta({
	title: "Composition Schema",
	description:
		"A schema using composition keywords (allOf, anyOf, oneOf, not) without a type constraint.",
	examples: [
		{
			allOf: [
				{ $ref: "#/components/schemas/Base" },
				{ type: "object", properties: { id: { type: "integer" } } },
			],
		},
	],
});

export type CompositionSchema = z.infer<typeof CompositionSchemaObject>;
