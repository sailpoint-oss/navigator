import z from "zod";
import { ResponseSchema } from "./response";

export const ResponsesObjectSchema = z
	.looseObject({
		default: ResponseSchema.optional().meta({ title: "default" }),
	})
	.meta({ title: "Responses" });

export type ResponsesObject = z.infer<typeof ResponsesObjectSchema>;


