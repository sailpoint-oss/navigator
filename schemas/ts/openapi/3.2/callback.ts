import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { PathItemObjectSchema } from "./path-item";

export const CallbackSchema: z.ZodType = z
	.union([
		ReferenceObjectSchema,
		z.record(z.string(), PathItemObjectSchema).meta({
			title: "Callback",
			description:
				"A map of possible out-of-band callbacks related to the parent operation.",
		}),
	])
	.meta({
		title: "Callback",
		description:
			"A map of possible out-of-band callbacks related to the parent operation.",
		examples: [
			{
				"{$request.body#/callbackUrl}": {
					post: { responses: { "200": { description: "OK" } } },
				},
			},
		],
	});

export type Callback = z.infer<typeof CallbackSchema>;
