import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { MediaTypeObjectSchema } from "./media-type";

export const RequestBodyObjectSchema = z
	.looseObject({
		description: z.string().optional().meta({ title: "description" }),
		content: z
			.record(z.string(), MediaTypeObjectSchema)
			.meta({ title: "content" }),
		required: z.boolean().optional().meta({ title: "required" }),
	})
	.meta({
		title: "RequestBody",
		description: "Describes a single request body.",
		examples: [
			{
				content: { "application/json": { schema: { type: "object" } } },
				required: true,
			},
		],
	});

export const RequestBodySchema = z
	.union([ReferenceObjectSchema, RequestBodyObjectSchema])
	.meta({
		title: "RequestBody",
		description: "Describes a single request body.",
	});

export type RequestBodyObject = z.infer<typeof RequestBodyObjectSchema>;
export type RequestBody = z.infer<typeof RequestBodySchema>;
