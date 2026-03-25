import z from "zod";
import { ItemsObjectSchema } from "./schema";

export const HeaderObjectSchema = z
	.looseObject({
		description: z.string().optional().meta({ title: "description" }),
		type: z
			.enum(["string", "number", "integer", "boolean", "array"])
			.meta({ title: "type" }),
		format: z.string().optional().meta({ title: "format" }),
		items: ItemsObjectSchema.optional().meta({ title: "items" }),
		collectionFormat: z
			.enum(["csv", "ssv", "tsv", "pipes"])
			.optional()
			.meta({ title: "collectionFormat" }),
		default: z.unknown().optional().meta({ title: "default" }),
		maximum: z.number().optional().meta({ title: "maximum" }),
		exclusiveMaximum: z.boolean().optional().meta({ title: "exclusiveMaximum" }),
		minimum: z.number().optional().meta({ title: "minimum" }),
		exclusiveMinimum: z.boolean().optional().meta({ title: "exclusiveMinimum" }),
		maxLength: z.number().int().min(0).optional().meta({ title: "maxLength" }),
		minLength: z.number().int().min(0).optional().meta({ title: "minLength" }),
		pattern: z.string().optional().meta({ title: "pattern" }),
		maxItems: z.number().int().min(0).optional().meta({ title: "maxItems" }),
		minItems: z.number().int().min(0).optional().meta({ title: "minItems" }),
		uniqueItems: z.boolean().optional().meta({ title: "uniqueItems" }),
		enum: z.array(z.unknown()).optional().meta({ title: "enum" }),
		multipleOf: z.number().optional().meta({ title: "multipleOf" }),
	})
	.meta({ title: "Header" });

export type HeaderObject = z.infer<typeof HeaderObjectSchema>;


