/**
 * Redocly OpenAPI Extensions
 *
 * Redocly extensions for API documentation and lifecycle management.
 * @see https://redocly.com/docs/cli/rules/
 */

import { defineExtension } from "../index.js";
import type { ExtensionSchemaMeta } from "../types.js";

/**
 * x-metadata: API lifecycle and ownership metadata.
 * Common pattern for tracking API status and team ownership.
 */
export const xMetadata: ExtensionSchemaMeta = defineExtension({
	name: "x-metadata",
	scope: ["info"],
	description:
		"API lifecycle and ownership metadata for tracking API status and team ownership",
	url: "https://redocly.com/docs/cli/custom-plugins/extended-types/",
	schema: (z) =>
		z.looseObject({
			lifecycle: z
				.union([
					z.literal("development"),
					z.literal("staging"),
					z.literal("production"),
					z.literal("deprecated"),
				])
				.optional(),
			"owner-team": z.string().optional(),
			"owner-email": z.email().optional(),
		}),
});

/**
 * x-logo: Custom logo for API documentation.
 */
export const xLogo: ExtensionSchemaMeta = defineExtension({
	name: "x-logo",
	scope: ["info"],
	description: "Custom logo configuration for API documentation rendering",
	url: "https://redocly.com/docs/api-reference-docs/specification-extensions/",
	schema: (z) =>
		z.looseObject({
			url: z.url(),
			altText: z.string().optional(),
			backgroundColor: z.string().optional(),
			href: z.url().optional(),
		}),
});

/**
 * x-tagGroups: Group tags for better documentation organization.
 */
export const xTagGroups: ExtensionSchemaMeta = defineExtension({
	name: "x-tagGroups",
	scope: ["root"],
	description: "Group tags into categories for better documentation navigation",
	url: "https://redocly.com/docs/api-reference-docs/specification-extensions/x-tag-groups/",
	schema: (z) =>
		z.array(
			z.looseObject({
				name: z.string(),
				tags: z.array(z.string()),
			}),
		),
});

/**
 * x-displayName: Override display name for tags.
 */
export const xDisplayName: ExtensionSchemaMeta = defineExtension({
	name: "x-displayName",
	scope: ["tag"],
	description:
		"Override the display name of a tag in documentation (alternative to tag name)",
	url: "https://redocly.com/docs/api-reference-docs/specification-extensions/",
	schema: (z) => z.string().min(1),
});

/**
 * x-traitTag: Mark a tag as a trait (applies to multiple operations).
 */
export const xTraitTag: ExtensionSchemaMeta = defineExtension({
	name: "x-traitTag",
	scope: ["tag"],
	description:
		"Mark a tag as a trait that describes a common behavior across operations",
	url: "https://redocly.com/docs/api-reference-docs/specification-extensions/",
	schema: (z) => z.boolean(),
});

/**
 * x-codeSamples: Add code samples to operations.
 */
export const xCodeSamples: ExtensionSchemaMeta = defineExtension({
	name: "x-codeSamples",
	scope: ["operation"],
	description: "Add code samples in various languages to demonstrate API usage",
	url: "https://redocly.com/docs/api-reference-docs/specification-extensions/x-code-samples/",
	schema: (z) =>
		z.array(
			z.looseObject({
				lang: z.string(),
				label: z.string().optional(),
				source: z.string(),
			}),
		),
});

/**
 * x-internal: Mark an operation or schema as internal (not for public docs).
 */
export const xInternal: ExtensionSchemaMeta = defineExtension({
	name: "x-internal",
	scope: ["operation", "schema", "parameter", "pathItem"],
	description:
		"Mark an element as internal, excluding it from public documentation",
	url: "https://redocly.com/docs/api-reference-docs/specification-extensions/",
	schema: (z) => z.boolean(),
});

/**
 * x-webhooks: Define webhooks (for OpenAPI 3.0 compatibility).
 */
export const xWebhooks: ExtensionSchemaMeta = defineExtension({
	name: "x-webhooks",
	scope: ["root"],
	description:
		"Define webhooks for OpenAPI 3.0 (native in 3.1+), describing callback URLs the API will call",
	url: "https://redocly.com/docs/api-reference-docs/specification-extensions/",
	schema: (z) => z.record(z.string(), z.unknown()),
});

/**
 * All Redocly extensions.
 */
export const redoclyExtensions: ExtensionSchemaMeta[] = [
	xMetadata,
	xLogo,
	xTagGroups,
	xDisplayName,
	xTraitTag,
	xCodeSamples,
	xInternal,
	xWebhooks,
];
