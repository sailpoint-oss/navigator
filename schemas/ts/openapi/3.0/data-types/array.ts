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
	minItems: z.number().int().min(0).optional().meta({ title: "minItems" }),
	maxItems: z.number().int().min(0).optional().meta({ title: "maxItems" }),
	uniqueItems: z.boolean().optional().meta({ title: "uniqueItems" }),
}).meta({
	title: "Array Object",
	description: "An Array Object defines the shape of an array value.",
	examples: [{ type: "array", items: { type: "string" } }],
});

export type ArraySchema = z.infer<typeof ArraySchemaObject>;
