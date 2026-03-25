import z from "zod";

export const OAuthFlowObjectSchema = z
	.looseObject({
		authorizationUrl: z
			.url()
			.meta({ title: "authorizationUrl" })
			.describe("The authorization URL to be used for this flow.")
			.optional(),
		tokenUrl: z
			.url()
			.meta({ title: "tokenUrl" })
			.describe("The token URL to be used for this flow.")
			.optional(),
		refreshUrl: z
			.url()
			.meta({ title: "refreshUrl" })
			.describe("The URL to be used for obtaining refresh tokens.")
			.optional(),
		scopes: z
			.record(z.string(), z.string())
			.meta({ title: "scopes" })
			.describe("The available scopes for the OAuth2 security scheme."),
	})
	.meta({
		title: "OAuthFlow",
		description: "Configuration details for a supported OAuth Flow.",
		examples: [
			{
				authorizationUrl: "https://auth.example.com/authorize",
				tokenUrl: "https://auth.example.com/token",
				scopes: { "read:pets": "Read pets" },
			},
		],
	});

export type OAuthFlowObject = z.infer<typeof OAuthFlowObjectSchema>;


