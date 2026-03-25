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
		summary: z
			.string()
			.meta({
				title: "summary",
				examples: ["Pet operations", "User management", "Order processing"],
			})
			.describe("A short summary of the tag. CommonMark syntax MAY be used.")
			.optional(),
		parent: z
			.string()
			.meta({
				title: "parent",
				examples: ["store", "users", "admin"],
			})
			.describe("The name of the parent tag for hierarchical organization.")
			.optional(),
		kind: z
			.union([z.literal("nav"), z.literal("badge"), z.literal("audience")])
			.meta({
				title: "kind",
				examples: ["nav", "badge", "audience"],
			})
			.describe("Tag classification for UI and grouping.")
			.optional(),
	})
	.meta({
		title: "Tag",
		description:
			"Tag metadata. OpenAPI 3.2 adds 'parent', 'kind', and 'summary' for hierarchical tags.",
		examples: [
			{ name: "pets", description: "Everything about your Pets", kind: "nav" },
		],
	});

export type TagObject = z.infer<typeof TagObjectSchema>;


