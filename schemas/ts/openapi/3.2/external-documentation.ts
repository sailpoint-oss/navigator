import z from "zod";

export const ExternalDocumentationObjectSchema = z
	.looseObject({
		description: z
			.string()
			.meta({ title: "description" })
			.describe("A short description of the target documentation.")
			.optional(),
		url: z
			.url()
			.meta({ title: "url" })
			.describe("The URL for the target documentation."),
	})
	.meta({
		title: "ExternalDocumentation",
		description:
			"Allows referencing an external resource for extended documentation.",
		examples: [
			{ description: "Find more info here", url: "https://docs.example.com" },
		],
	});

export type ExternalDocumentationObject = z.infer<
	typeof ExternalDocumentationObjectSchema
>;


