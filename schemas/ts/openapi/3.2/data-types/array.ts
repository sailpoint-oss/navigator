import z from "zod";
import { SchemaObjectSchema } from "../schema";
import { BaseSchemaObjectSchema } from "./base";

export const ArraySchemaObject = BaseSchemaObjectSchema.extend({
	type: z.literal("array"),
	get items() {
		return z.lazy(() => SchemaObjectSchema).optional().meta({
			title: "items",
			description: "Schema for array items.",
		});
	},
	get prefixItems() {
		return z
			.array(z.lazy(() => SchemaObjectSchema) as z.ZodType)
			.optional()
			.meta({
				title: "prefixItems",
				description:
					"Tuple validation: schemas for items by index (array position).",
			});
	},
	get contains() {
		return z.lazy(() => SchemaObjectSchema).optional().meta({
			title: "contains",
			description:
				"Require the array to contain at least one item matching this subschema.",
		});
	},
	minItems: z.number().int().min(0).optional().meta({ title: "minItems" }),
	maxItems: z.number().int().min(0).optional().meta({ title: "maxItems" }),
	minContains: z
		.number()
		.int()
		.min(0)
		.optional()
		.meta({ title: "minContains" }),
	maxContains: z
		.number()
		.int()
		.min(0)
		.optional()
		.meta({ title: "maxContains" }),
	uniqueItems: z.boolean().optional().meta({ title: "uniqueItems" }),
}).meta({
	title: "Array Object",
	description: "An Array Object defines the shape of an array value.",
	examples: [{ type: "array", items: { type: "string" } }],
});

export type ArraySchema = z.infer<typeof ArraySchemaObject>;
