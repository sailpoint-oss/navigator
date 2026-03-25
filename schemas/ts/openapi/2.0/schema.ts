import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ExternalDocumentationObjectSchema } from "./external-documentation";
import { XmlObjectSchema } from "./xml";

// Forward declarations
export type SchemaObject = z.infer<typeof SchemaObjectSchema>;
export type ItemsObject = z.infer<typeof ItemsObjectSchema>;

export const ItemsObjectSchema = z
	.looseObject({
		$ref: ReferenceObjectSchema.shape.$ref.optional(),
		type: z
			.enum(["string", "number", "integer", "boolean", "array", "object"])
			.optional()
			.meta({ title: "type" }),
		format: z.string().optional().meta({ title: "format" }),
		get items() {
			return ItemsObjectSchema.optional().meta({ title: "items" });
		},
		collectionFormat: z
			.enum(["csv", "ssv", "tsv", "pipes"])
			.optional()
			.meta({ title: "collectionFormat" }),
		default: z.unknown().optional().meta({ title: "default" }),
		maximum: z.number().optional().meta({ title: "maximum" }),
		exclusiveMaximum: z
			.boolean()
			.optional()
			.meta({ title: "exclusiveMaximum" }),
		minimum: z.number().optional().meta({ title: "minimum" }),
		exclusiveMinimum: z
			.boolean()
			.optional()
			.meta({ title: "exclusiveMinimum" }),
		maxLength: z.number().int().min(0).optional().meta({ title: "maxLength" }),
		minLength: z.number().int().min(0).optional().meta({ title: "minLength" }),
		pattern: z.string().optional().meta({ title: "pattern" }),
		maxItems: z.number().int().min(0).optional().meta({ title: "maxItems" }),
		minItems: z.number().int().min(0).optional().meta({ title: "minItems" }),
		uniqueItems: z.boolean().optional().meta({ title: "uniqueItems" }),
		enum: z.array(z.unknown()).optional().meta({ title: "enum" }),
		multipleOf: z.number().optional().meta({ title: "multipleOf" }),
	})
	.meta({ title: "Items" });

export const SchemaObjectSchema = z
	.looseObject({
		$ref: ReferenceObjectSchema.shape.$ref.optional().meta({ title: "$ref" }),
		title: z.string().optional().meta({ title: "title" }),
		description: z.string().optional().meta({ title: "description" }),
		type: z
			.enum(["string", "number", "integer", "boolean", "array", "object"])
			.optional()
			.meta({ title: "type" }),
		format: z.string().optional().meta({ title: "format" }),
		required: z.array(z.string()).optional().meta({ title: "required" }),
		get properties() {
			return z
				.record(z.string(), SchemaObjectSchema)
				.optional()
				.meta({ title: "properties" });
		},
		get additionalProperties() {
			return z
				.union([z.boolean(), SchemaObjectSchema])
				.optional()
				.meta({ title: "additionalProperties" });
		},
		get items() {
			return ItemsObjectSchema.optional().meta({ title: "items" });
		},
		get allOf() {
			return z.array(SchemaObjectSchema).optional().meta({ title: "allOf" });
		},
		default: z.unknown().optional().meta({ title: "default" }),
		discriminator: z.string().optional().meta({ title: "discriminator" }),
		readOnly: z.boolean().optional().meta({ title: "readOnly" }),
		xml: XmlObjectSchema.optional().meta({ title: "xml" }),
		externalDocs: ExternalDocumentationObjectSchema.optional().meta({
			title: "externalDocs",
		}),
		example: z.unknown().optional().meta({ title: "example" }),
		maximum: z.number().optional().meta({ title: "maximum" }),
		exclusiveMaximum: z
			.boolean()
			.optional()
			.meta({ title: "exclusiveMaximum" }),
		minimum: z.number().optional().meta({ title: "minimum" }),
		exclusiveMinimum: z
			.boolean()
			.optional()
			.meta({ title: "exclusiveMinimum" }),
		maxLength: z.number().int().min(0).optional().meta({ title: "maxLength" }),
		minLength: z.number().int().min(0).optional().meta({ title: "minLength" }),
		pattern: z.string().optional().meta({ title: "pattern" }),
		maxItems: z.number().int().min(0).optional().meta({ title: "maxItems" }),
		minItems: z.number().int().min(0).optional().meta({ title: "minItems" }),
		uniqueItems: z.boolean().optional().meta({ title: "uniqueItems" }),
		maxProperties: z
			.number()
			.int()
			.min(0)
			.optional()
			.meta({ title: "maxProperties" }),
		minProperties: z
			.number()
			.int()
			.min(0)
			.optional()
			.meta({ title: "minProperties" }),
		enum: z.array(z.unknown()).optional().meta({ title: "enum" }),
		multipleOf: z.number().optional().meta({ title: "multipleOf" }),
	})
	.meta({ title: "Schema" });
