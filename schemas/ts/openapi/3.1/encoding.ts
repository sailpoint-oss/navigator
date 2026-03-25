import z from "zod";
import { HeaderSchema } from "./header";

export const EncodingObjectSchema = z
	.looseObject({
		contentType: z
			.string()
			.meta({ title: "contentType" })
			.describe("The Content-Type for encoding a specific property.")
			.optional(),
		headers: z
			.record(z.string(), HeaderSchema)
			.meta({ title: "headers" })
			.optional(),
		style: z
			.union([
				z.literal("form"),
				z.literal("spaceDelimited"),
				z.literal("pipeDelimited"),
				z.literal("deepObject"),
			])
			.meta({ title: "style" })
			.optional(),
		explode: z.boolean().optional().meta({ title: "explode" }),
		allowReserved: z
			.boolean()
			.optional()
			.meta({ title: "allowReserved" }),
	})
	.meta({
		title: "Encoding",
		description:
			"A single encoding definition applied to a single schema property.",
		examples: [{ contentType: "application/json", style: "form" }],
	});

export type EncodingObject = z.infer<typeof EncodingObjectSchema>;


