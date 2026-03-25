import z from "zod";
import { ResponseSchema } from "./response";

export const ResponsesObjectSchema = z
	.record(
		z.union([z.string().regex(/^[1-5][0-9]{2}$/), z.literal("default")]),
		ResponseSchema,
	)
	.meta({
		title: "Responses",
		description: "A container for the expected responses of an operation.",
		examples: [
			{
				"200": { description: "Success" },
				"404": { description: "Not found" },
			},
		],
	});

export type ResponsesObject = z.infer<typeof ResponsesObjectSchema>;


