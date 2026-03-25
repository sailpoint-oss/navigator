import z from "zod";
import { OAuthFlowObjectSchema } from "./oauth-flow";

export const OAuthFlowsObjectSchema = z
	.looseObject({
		implicit: OAuthFlowObjectSchema.optional().meta({ title: "implicit" }),
		password: OAuthFlowObjectSchema.optional().meta({ title: "password" }),
		clientCredentials: OAuthFlowObjectSchema.optional().meta({
			title: "clientCredentials",
		}),
		authorizationCode: OAuthFlowObjectSchema.optional().meta({
			title: "authorizationCode",
		}),
		device: OAuthFlowObjectSchema.optional()
			.meta({ title: "device" })
			.describe("Configuration for the Device Authorization Grant flow."),
	})
	.meta({
		title: "OAuthFlows",
		description: "Allows configuration of the supported OAuth Flows.",
		examples: [
			{
				implicit: {
					authorizationUrl: "https://auth.example.com/authorize",
					scopes: {},
				},
			},
		],
	});

export type OAuthFlowsObject = z.infer<typeof OAuthFlowsObjectSchema>;


