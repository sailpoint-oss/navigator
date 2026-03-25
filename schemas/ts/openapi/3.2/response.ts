import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { HeaderSchema } from "./header";
import { LinkSchema } from "./link";
import { MediaTypeObjectSchema } from "./media-type";

export const ResponseObjectSchema = z
	.looseObject({
		description: z
			.string()
			.meta({ title: "description" })
			.describe("A description of the response."),
		headers: z
			.record(z.string(), HeaderSchema)
			.meta({ title: "headers" })
			.optional(),
		content: z
			.record(z.string(), MediaTypeObjectSchema)
			.meta({ title: "content" })
			.optional(),
		links: z.record(z.string(), LinkSchema).meta({ title: "links" }).optional(),
	})
	.meta({
		title: "Response",
		description: "Describes a single response from an API Operation.",
		examples: [
			{
				description: "Successful response",
				content: { "application/json": { schema: { type: "object" } } },
			},
		],
	});

export const ResponseSchema = z
	.union([ReferenceObjectSchema, ResponseObjectSchema])
	.meta({
		title: "Response",
		description: "Describes a single response from an API Operation.",
	});

export type ResponseObject = z.infer<typeof ResponseObjectSchema>;
export type Response = z.infer<typeof ResponseSchema>;
