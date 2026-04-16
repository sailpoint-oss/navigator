/**
 * OpenAPI Base Schemas - Shared across all OpenAPI 3.x versions
 *
 * These schemas are IDENTICAL across OpenAPI 3.0, 3.1, and 3.2.
 * Version-specific schemas are defined in their respective modules.
 *
 * @module engine/schemas/openapi-base
 */
import { z } from "zod";

// ============================================================================
// Contact Object
// ============================================================================

/**
 * Contact Object - Contact information for the exposed API.
 * Identical across all OpenAPI 3.x versions.
 */
export const ContactSchema = z
	.looseObject({
		name: z
			.string()
			.meta({
				title: "name",
				description: "The identifying name of the contact person/organization.",
				examples: ["API Support", "Developer Team", "John Smith"],
			})
			.optional(),
		url: z
			.url()
			.meta({
				title: "url",
				description:
					"The URL pointing to the contact information. Must be a valid URL.",
				examples: [
					"https://www.example.com/support",
					"https://developer.example.com",
				],
			})
			.optional(),
		email: z
			.email()
			.meta({
				title: "email",
				description:
					"The email address of the contact person/organization. Must be a valid email format.",
				examples: ["support@example.com", "api@company.io", "dev-team@org.com"],
			})
			.optional(),
	})
	
	.meta({
		title: "Contact",
		description:
			"Contact information for the exposed API. Provides ways for API consumers to reach the maintainers.",
		examples: [
			{
				name: "API Support",
				url: "https://example.com/support",
				email: "support@example.com",
			},
		],
	});

export type Contact = z.infer<typeof ContactSchema>;

// ============================================================================
// License Object
// ============================================================================

export const LicenseIdentifierSchema = z.string().meta({
	title: "License Identifier",
	description: "An SPDX license expression for the API.",
	examples: ["Apache-2.0", "MIT", "BSD-3-Clause", "GPL-3.0-only"],
});

export const LicenseUrlSchema = z.url().meta({
	title: "License URL",
	description: "A URL to the license used for the API.",
	examples: [
		"https://www.apache.org/licenses/LICENSE-2.0.html",
		"https://opensource.org/licenses/MIT",
	],
});

export const LicenseNameSchema = z.string().meta({
	title: "License Name",
	description: "The license name used for the API.",
	examples: ["Apache 2.0", "MIT", "BSD-3-Clause", "GPL-3.0"],
});

/**
 * License Object - License information for the exposed API.
 * Note: 3.0 only has name+url; 3.1+ adds identifier (mutually exclusive with url).
 * This base schema uses the 3.1+ superset shape.
 */
export const LicenseObjectSchema = z
	.looseObject({
		name: LicenseNameSchema,
		url: LicenseUrlSchema.optional(),
		identifier: LicenseIdentifierSchema.optional(),
	})
	.superRefine((value, ctx) => {
		if (value.url && value.identifier) {
			ctx.addIssue({
				code: "custom",
				message: "License Object: 'url' and 'identifier' are mutually exclusive",
				path: ["url"],
			});
		}
	})
	.meta({
		title: "License",
		description: "License information for the exposed API.",
		examples: [
			{ name: "Apache 2.0", identifier: "Apache-2.0" },
			{ name: "MIT", url: "https://opensource.org/licenses/MIT" },
		],
	});

export type LicenseObject = z.infer<typeof LicenseObjectSchema>;

// Back-compat aliases (central schema index expects these names)
export const LicenseSchema = LicenseObjectSchema;
export type License = LicenseObject;

// ============================================================================
// External Documentation Object
// ============================================================================

/**
 * External Documentation Object.
 * Identical across all OpenAPI 3.x versions.
 */
export const ExternalDocumentationSchema = z
	.looseObject({
		description: z
			.string()
			.describe(
				"A description of the target documentation. CommonMark syntax MAY be used for rich text.",
			)
			.optional()
			.meta({
				title: "description",
				examples: [
					"Find more info here",
					"Complete API documentation",
					"Authentication guide",
				],
			}),
		url: z
			.url()
			.describe("REQUIRED. The URL for the target documentation.")
			.meta({
				title: "url",
				examples: [
					"https://docs.example.com",
					"https://wiki.example.com/api-guide",
					"https://github.com/org/repo/wiki",
			],
		}),
	})
	
	.meta({
		title: "ExternalDocumentation",
		description:
			"Allows referencing an external resource for extended documentation.",
		examples: [
			{ description: "Find more info here", url: "https://docs.example.com" },
		],
	});

export type ExternalDocumentation = z.infer<typeof ExternalDocumentationSchema>;

// ============================================
// Reference Objects
// ============================================

export const InternalRefSchema = z
	.string()
	.regex(/^#.*/)
	.meta({
		title: "Internal JSON Pointer Reference",
		description: "Internal JSON Pointer reference",
		examples: [
			"#/components/schemas/User",
			"#/components/schemas/Error",
			"#/components/responses/NotFound",
			"#/components/parameters/PageParam",
			"#/components/requestBodies/UserInput",
		],
	});

export type InternalRef = z.infer<typeof InternalRefSchema>;

export const UrlRefSchema = z
	.string()
	.regex(/^https?:\/\//)
	.meta({
		title: "URL Reference",
		description: "URL reference",
		examples: [
			"https://example.com/schemas/Pet.yaml",
			"https://example.com/schemas/User.json",
			"https://example.com/responses/Not-Found-Error.yaml",
			"https://example.com/parameters.yaml#/PageParam",
			"https://example.com/requestBodies.yaml#/requestBodies/UserInput",
		],
	});

export type UrlRef = z.infer<typeof UrlRefSchema>;

export const FileRefSchema = z
	.string()
	.regex(/^[^#\s]+(\.ya?ml|\.json)(#\/.*)?$/i)
	.meta({
		title: "Relative File Reference",
		description: "Relative file reference",
		examples: [
			"./schemas/Pet.yaml",
			"../common/types.yaml",
			"schemas/Pet.yaml",
			"./v2/components/parameters.yaml#/components/parameters/LimitParam",
			"../common/types.yaml#/properties/id",
			"schemas/Pet.yaml#/properties/name",
		],
	});

export type FileRef = z.infer<typeof FileRefSchema>;

export const ReferenceSchema = z
	.union([InternalRefSchema, UrlRefSchema, FileRefSchema])
	.meta({
		title: "Reference",
		description: "A reference to another component in the specification.",
		examples: [
			"#/components/schemas/Pet",
			"https://example.com/schemas/Pet.yaml",
			"./schemas/Pet.yaml",
		],
	});

export type Reference = z.infer<typeof ReferenceSchema>;

export const ReferenceObjectSchema = z
	.looseObject({
		$ref: ReferenceSchema,
		summary: z
			.string()
			.meta({
				title: "summary",
				description:
					"A short summary which by default SHOULD override that of the referenced component.",
			})
			.optional(),
		description: z
			.string()
			.meta({
				title: "description",
				description:
					"A description which by default SHOULD override that of the referenced component.",
			})
			.optional(),
	})
	
	.meta({
		title: "Reference Object",
		description: "A reference to another component in the specification.",
		examples: [
			{ $ref: "#/components/schemas/Pet" },
			{
				$ref: "./schemas/User.yaml",
				summary: "A user",
				description: "A user is a type of person.",
			},
			{
				$ref: "definitions.yaml#/Pet",
				summary: "A pet",
				description: "A pet is a type of animal.",
			},
			{
				$ref: "https://example.com/schemas/Pet.yaml",
				summary: "A pet",
				description: "A pet is a type of animal.",
			},
		],
	});

export type ReferenceObject = z.infer<typeof ReferenceObjectSchema>;

// ============================================================================
// Security Requirement Object
// ============================================================================

/**
 * Security Requirement Object.
 * Identical across all OpenAPI 3.x versions.
 */
export const SecurityRequirementSchema = z
	.record(z.string(), z.array(z.string()))
	.meta({
		title: "SecurityRequirement",
		description:
			"Security requirement object. Keys are security scheme names from components/securitySchemes. Values are arrays of required scopes (empty for non-OAuth2).",
		examples: [
			{ api_key: [] },
			{ bearerAuth: [] },
			{ oauth2: ["read:pets", "write:pets"] },
			{ basicAuth: [], api_key: [] },
		],
	});

export type SecurityRequirement = z.infer<typeof SecurityRequirementSchema>;

// ============================================================================
// XML Object
// ============================================================================

/**
 * XML Object.
 * Identical across all OpenAPI 3.x versions.
 */
export const XMLSchema = z
	.looseObject({
		name: z
			.string()
			.describe(
				"Replaces the name of the element/attribute used for the described schema property.",
			)
			.optional()
			.meta({
				title: "name",
				examples: ["animal", "item", "user-data"],
			}),
		namespace: z
			.url()
			.describe("The URI of the namespace definition.")
			.optional()
			.meta({
				title: "namespace",
				examples: ["http://example.com/schema/pet", "urn:example:animals"],
			}),
		prefix: z
			.string()
			.describe("The prefix to be used for the name.")
			.optional()
			.meta({
				title: "prefix",
				examples: ["smp", "ns1", "pet"],
			}),
		attribute: z
			.boolean()
			.optional()
			.describe(
				"When true, the property is serialized as an XML attribute instead of an element. Default: false.",
			)
			.meta({
				title: "attribute",
				examples: [true, false],
			}),
		wrapped: z
			.boolean()
			.optional()
			.describe(
				"Only for arrays. When true, wraps the array in an outer element. Default: false.",
			)
			.meta({
				title: "wrapped",
				examples: [true, false],
			}),
	})
	
	.meta({
		title: "XML",
		description:
			"Metadata for fine-tuned XML serialization. Controls element names, namespaces, prefixes, and array wrapping.",
		examples: [
			{ name: "animal", prefix: "pet", namespace: "http://example.com/pet" },
			{ wrapped: true, name: "tags" },
		],
	});

export type XML = z.infer<typeof XMLSchema>;

// ============================================================================
// OAuth Flow Object (Base)
// ============================================================================

/**
 * OAuth Flow Object (Base).
 * Identical across all OpenAPI 3.x versions.
 */
export const OAuthFlowSchema = z
	.looseObject({
		authorizationUrl: z
			.url()
			.describe(
				"REQUIRED for implicit and authorizationCode flows. The authorization URL for this OAuth2 flow.",
			)
			.optional()
			.meta({
				title: "authorizationUrl",
				examples: [
					"https://auth.example.com/oauth/authorize",
					"https://login.example.com/oauth2/v2.0/authorize",
				],
			}),
		tokenUrl: z
			.url()
			.describe(
				"REQUIRED for password, clientCredentials, and authorizationCode flows. The token URL for this OAuth2 flow.",
			)
			.optional()
			.meta({
				title: "tokenUrl",
				examples: [
					"https://auth.example.com/oauth/token",
					"https://login.example.com/oauth2/v2.0/token",
				],
			}),
		refreshUrl: z
			.url()
			.describe("The URL to be used for obtaining refresh tokens.")
			.optional()
			.meta({
				title: "refreshUrl",
				examples: ["https://auth.example.com/oauth/refresh"],
			}),
		scopes: z
			.record(z.string(), z.string())
			.describe(
				"REQUIRED. Map of scope names to descriptions. Keys are scope names, values are descriptions.",
			)
			.meta({
				title: "scopes",
				examples: [
					{
						"read:pets": "Read access to pets",
						"write:pets": "Write access to pets",
					},
					{
						"user:email": "Access user email",
						"user:profile": "Access user profile",
					},
				],
			}),
	})
	
	.meta({
		title: "OAuthFlow",
		description:
			"Configuration for an OAuth 2.0 flow. Required fields depend on the flow type (implicit, password, clientCredentials, authorizationCode).",
		examples: [
			{
				authorizationUrl: "https://auth.example.com/oauth/authorize",
				tokenUrl: "https://auth.example.com/oauth/token",
				scopes: {
					"read:pets": "Read access to pets",
					"write:pets": "Write access to pets",
				},
			},
		],
	});

export type OAuthFlow = z.infer<typeof OAuthFlowSchema>;
