import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ItemsObjectSchema, SchemaObjectSchema } from "./schema";

const ParameterLocationSchema = z.enum([
	"query",
	"header",
	"path",
	"formData",
	"body",
]);

const SimpleParameterBaseSchema = z.looseObject({
	name: z
		.string()
		.describe("REQUIRED. The name of the parameter.")
		.meta({ title: "name" }),
	in: ParameterLocationSchema.describe(
		"REQUIRED. The location of the parameter.",
	).meta({
		title: "in",
	}),
	description: z.string().optional().meta({ title: "description" }),
	required: z.boolean().optional().meta({ title: "required" }),
	type: z
		.enum(["string", "number", "integer", "boolean", "array", "file"])
		.describe("REQUIRED. The type of the parameter.")
		.meta({ title: "type" }),
	format: z.string().optional().meta({ title: "format" }),
	allowEmptyValue: z.boolean().optional().meta({ title: "allowEmptyValue" }),
	items: ItemsObjectSchema.optional().meta({ title: "items" }),
	collectionFormat: z
		.enum(["csv", "ssv", "tsv", "pipes", "multi"])
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
});

const BodyParameterObjectSchema = z
	.looseObject({
		name: z.string().meta({ title: "name" }),
		in: z.literal("body").meta({ title: "in" }),
		description: z.string().optional().meta({ title: "description" }),
		required: z.boolean().optional().meta({ title: "required" }),
		schema: SchemaObjectSchema.describe(
			"REQUIRED. The schema defining the type.",
		).meta({
			title: "schema",
		}),
	})
	
	.meta({ title: "BodyParameter" });

const NonBodyParameterObjectSchema = SimpleParameterBaseSchema.meta({
	title: "NonBodyParameter",
});

export const ParameterObjectSchema = z
	.union([
		BodyParameterObjectSchema,
		NonBodyParameterObjectSchema,
		ReferenceObjectSchema,
	])
	.meta({ title: "Parameter" });

export type ParameterObject = z.infer<typeof ParameterObjectSchema>;
