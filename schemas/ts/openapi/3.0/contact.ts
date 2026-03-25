import z from "zod";

export const ContactObjectSchema = z
	.looseObject({
		name: z
			.string()
			.meta({
				title: "name",
				examples: ["API Support", "Developer Team", "John Smith"],
			})
			.describe("The identifying name of the contact person/organization.")
			.optional(),
		url: z
			.url()
			.meta({
				title: "url",
				examples: [
					"https://www.example.com/support",
					"https://developer.example.com",
				],
			})
			.describe(
				"The URL pointing to the contact information. Must be a valid URL.",
			)
			.optional(),
		email: z
			.email()
			.meta({
				title: "email",
				examples: ["support@example.com", "api@company.io"],
			})
			.describe("The email address of the contact person/organization.")
			.optional(),
	})
	.meta({
		title: "Contact",
		description: "Contact information for the exposed API.",
		examples: [
			{
				name: "API Support",
				url: "https://example.com/support",
				email: "support@example.com",
			},
		],
	});
