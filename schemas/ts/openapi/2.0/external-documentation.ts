import z from "zod";
import { UrlStringSchema } from "./primitives";

export const ExternalDocumentationObjectSchema = z
	.looseObject({
		description: z
			.string()
			.optional()
			.describe("A short description of the target documentation.")
			.meta({ title: "description" }),
		url: UrlStringSchema.describe(
			"REQUIRED. The URL for the target documentation.",
		).meta({ title: "url" }),
	})
	.meta({ title: "ExternalDocumentation" });

export type ExternalDocumentationObject = z.infer<
	typeof ExternalDocumentationObjectSchema
>;


