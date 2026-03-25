import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { HeaderObjectSchema } from "./header";
import { SchemaObjectSchema } from "./schema";

export const ResponseObjectSchema = z
	.looseObject({
		description: z
			.string()
			.describe("REQUIRED. A short description of the response.")
			.meta({ title: "description" }),
		schema: SchemaObjectSchema.optional().meta({ title: "schema" }),
		headers: z
			.record(z.string(), HeaderObjectSchema)
			.optional()
			.meta({ title: "headers" }),
		examples: z
			.record(z.string(), z.unknown())
			.optional()
			.meta({ title: "examples" }),
	})
	.meta({ title: "Response" });

export const ResponseSchema = z
	.union([ResponseObjectSchema, ReferenceObjectSchema])
	.meta({ title: "Response" });

export type ResponseObject = z.infer<typeof ResponseObjectSchema>;
export type Response = z.infer<typeof ResponseSchema>;
