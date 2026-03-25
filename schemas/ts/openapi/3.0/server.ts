import z from "zod";
import { ServerVariableObjectSchema } from "./server-variable";

export const ServerObjectSchema = z
	.looseObject({
		url: z
			.string()
			.meta({
				title: "url",
				examples: [
					"https://api.example.com/v1",
					"https://{environment}.api.example.com",
					"{scheme}://{host}:{port}/api",
					"/api/v1",
				],
			})
			.describe(
				"A URL to the target host. Supports server variables in {braces}. May be relative.",
			),
		description: z
			.string()
			.meta({
				title: "description",
				examples: ["Production server", "Staging server", "Development server"],
			})
			.describe("An optional description of the server.")
			.optional(),
		variables: z
			.record(z.string(), ServerVariableObjectSchema)
			.meta({
				title: "variables",
				examples: [
					{
						environment: { default: "prod", enum: ["prod", "staging", "dev"] },
					},
				],
			})
			.describe("A map of server variables for URL template substitution.")
			.optional(),
	})
	.meta({
		title: "Server",
		description: "An object representing a Server.",
		examples: [
			{
				url: "https://api.example.com/v1",
				description: "Production server",
			},
		],
	});

export type ServerObject = z.infer<typeof ServerObjectSchema>;


