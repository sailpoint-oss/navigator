import z from "zod";
import { ContactObjectSchema } from "./contact";
import { LicenseObjectSchema } from "./license";

export const InfoObjectSchema = z
	.looseObject({
		title: z
			.string()
			.meta({
				title: "title",
				examples: [
					"Pet Store API",
					"User Management Service",
					"Payment Gateway",
				],
			})
			.describe("The title of the API."),
		version: z
			.string()
			.meta({
				title: "version",
				examples: ["1.0.0", "2.3.1", "0.1.0-beta", "1.0.0-rc.1"],
			})
			.describe(
				"The version of the API document (not the OpenAPI spec version).",
			),
		summary: z
			.string()
			.meta({
				title: "summary",
				examples: [
					"A simple pet store API",
					"Manages user accounts and authentication",
				],
			})
			.describe("A short summary of the API.")
			.optional(),
		description: z
			.string()
			.meta({
				title: "description",
				examples: [
					"A sample API for managing pets",
					"This API handles user authentication and profile management",
				],
			})
			.describe(
				"A description of the API. CommonMark syntax MAY be used for rich text.",
			)
			.optional(),
		termsOfService: z
			.url()
			.meta({
				title: "termsOfService",
				examples: ["https://example.com/terms", "https://api.example.com/tos"],
			})
			.describe("A URL to the Terms of Service for the API.")
			.optional(),
		contact: ContactObjectSchema.optional(),
		license: LicenseObjectSchema.optional(),
	})
	.meta({
		title: "Info",
		description:
			"Provides metadata about the API. REQUIRED fields: title, version.",
		examples: [
			{ title: "Pet Store API", version: "1.0.0", description: "A sample API" },
		],
	});
