import z from "zod";
import { ApiKeyLocationSchema, HttpAuthSchemeSchema } from "../data-types";
import { ReferenceObjectSchema } from "../openapi-base";
import { OAuthFlowsObjectSchema } from "./oauth-flows";

const ApiKeySecuritySchemeObjectSchema = z
	.looseObject({
		type: z.literal("apiKey").meta({ title: "type", examples: ["apiKey"] }),
		name: z
			.string()
			.meta({
				title: "name",
				examples: ["X-API-Key", "api_key", "Authorization", "X-Auth-Token"],
			})
			.describe("The name of the header, query, or cookie parameter."),
		in: ApiKeyLocationSchema.describe(
			"Location of the API key: 'header' (most common), 'query', or 'cookie'.",
		),
		description: z
			.string()
			.meta({ title: "description" })
			.describe("A description of the security scheme.")
			.optional(),
	})
	.meta({
		title: "ApiKeySecurityScheme",
		description: "API Key security scheme.",
		examples: [{ type: "apiKey", name: "X-API-Key", in: "header" }],
	});

const HttpSecuritySchemeObjectSchema = z
	.looseObject({
		type: z.literal("http").meta({ title: "type", examples: ["http"] }),
		scheme: HttpAuthSchemeSchema.describe(
			"HTTP auth scheme (IANA registered). Common: 'bearer', 'basic'.",
		),
		bearerFormat: z
			.string()
			.meta({ title: "bearerFormat", examples: ["JWT", "opaque", "Bearer"] })
			.describe(
				"Hint for bearer token format. Only applicable when scheme='bearer'.",
			)
			.optional(),
		description: z
			.string()
			.meta({ title: "description" })
			.describe("A description of the security scheme.")
			.optional(),
	})
	.meta({
		title: "HttpSecurityScheme",
		description: "HTTP authentication security scheme (Basic, Bearer, etc.).",
		examples: [{ type: "http", scheme: "bearer", bearerFormat: "JWT" }],
	});

const MutualTLSSecuritySchemeObjectSchema = z
	.looseObject({
		type: z.literal("mutualTLS").meta({ title: "type" }),
		description: z.string().optional().meta({ title: "description" }),
	})
	.meta({
		title: "MutualTLSSecurityScheme",
		description: "Mutual TLS security scheme",
		examples: [
			{ type: "mutualTLS", description: "Client certificate authentication" },
		],
	});

const OAuth2SecuritySchemeObjectSchema = z
	.looseObject({
		type: z.literal("oauth2").meta({ title: "type" }),
		flows: OAuthFlowsObjectSchema.meta({ title: "flows" }),
		description: z.string().optional().meta({ title: "description" }),
	})
	.meta({
		title: "OAuth2SecurityScheme",
		description: "OAuth2 security scheme",
	});

const OpenIdConnectSecuritySchemeObjectSchema = z
	.looseObject({
		type: z.literal("openIdConnect").meta({ title: "type" }),
		openIdConnectUrl: z
			.url()
			.meta({ title: "openIdConnectUrl" })
			.describe("OpenID Connect URL to discover OAuth2 configuration values."),
		description: z.string().optional().meta({ title: "description" }),
	})
	.meta({
		title: "OpenIdConnectSecurityScheme",
		description: "OpenID Connect security scheme",
	});

export const SecuritySchemeSchema = z
	.union([
		ReferenceObjectSchema,
		ApiKeySecuritySchemeObjectSchema,
		HttpSecuritySchemeObjectSchema,
		MutualTLSSecuritySchemeObjectSchema,
		OAuth2SecuritySchemeObjectSchema,
		OpenIdConnectSecuritySchemeObjectSchema,
	])
	.meta({
		title: "SecurityScheme",
		description:
			"Defines a security scheme that can be used by the operations.",
		examples: [
			{ type: "apiKey", name: "X-API-Key", in: "header" },
			{ type: "http", scheme: "bearer" },
		],
	});

export type SecurityScheme = z.infer<typeof SecuritySchemeSchema>;
