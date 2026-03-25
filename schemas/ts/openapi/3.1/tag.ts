import z from "zod";
import { ExternalDocumentationObjectSchema } from "./external-documentation";

export const TagObjectSchema = z
	.looseObject({
		name: z
			.string()
			.meta({
				title: "name",
				examples: ["pets", "users", "orders", "authentication", "admin"],
			})
			.describe("The name of the tag."),
		description: z
			.string()
			.meta({
				title: "description",
				examples: [
					"Everything about your Pets",
					"Operations for user management",
					"Access to orders",
				],
			})
			.describe("A description for the tag. CommonMark syntax MAY be used.")
			.optional(),
		externalDocs: ExternalDocumentationObjectSchema.optional(),
	})
	.meta({
		title: "Tag",
		description: "Tag metadata.",
		examples: [{ name: "pets", description: "Everything about your Pets" }],
	});

export type TagObject = z.infer<typeof TagObjectSchema>;


