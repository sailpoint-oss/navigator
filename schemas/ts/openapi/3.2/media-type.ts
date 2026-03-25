import z from "zod";
import { EncodingObjectSchema } from "./encoding";
import { ExampleSchema } from "./example";
import { SchemaObjectSchema } from "./schema";

export const MediaTypeObjectSchema = z
	.looseObject({
		schema: SchemaObjectSchema.optional().meta({ title: "schema" }),
		example: z.unknown().optional().meta({ title: "example" }),
		examples: z
			.record(z.string(), ExampleSchema)
			.meta({ title: "examples" })
			.optional(),
		encoding: z
			.record(z.string(), EncodingObjectSchema)
			.meta({ title: "encoding" })
			.optional(),
		itemSchema: SchemaObjectSchema.meta({ title: "itemSchema" })
			.describe("Schema for individual items in streaming responses.")
			.optional(),
		itemEncoding: z
			.record(z.string(), EncodingObjectSchema)
			.meta({ title: "itemEncoding" })
			.describe("Encoding for individual items in streaming responses.")
			.optional(),
	})
	.meta({
		title: "MediaType",
		description:
			"Each Media Type Object provides schema and examples for the media type identified by its key.",
		examples: [{ schema: { type: "object" }, example: { id: 1 } }],
	});

export type MediaTypeObject = z.infer<typeof MediaTypeObjectSchema>;
