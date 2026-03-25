/**
 * Speakeasy OpenAPI Extensions
 *
 * Speakeasy extensions for SDK generation configuration.
 * @see https://speakeasy.com/docs
 */

import { defineExtension } from "../index.js";
import type { ExtensionSchemaMeta } from "../types.js";

/**
 * x-speakeasy-entity: Marks a schema as a distinct entity for SDK generation.
 * Can be a boolean to enable/disable or a string for custom naming.
 */
export const xSpeakeasyEntity: ExtensionSchemaMeta = defineExtension({
	name: "x-speakeasy-entity",
	scope: ["schema"],
	description:
		"Marks a schema as a distinct entity for SDK generation, affecting how types are generated",
	url: "https://speakeasy.com/docs/customize-sdks/entities",
	schema: (z) => z.union([z.boolean(), z.string()]),
});

/**
 * x-speakeasy-name-override: Override the generated name for an operation or schema.
 */
export const xSpeakeasyNameOverride: ExtensionSchemaMeta = defineExtension({
	name: "x-speakeasy-name-override",
	scope: ["operation", "schema", "parameter"],
	description:
		"Override the generated name for an operation, schema, or parameter in SDKs",
	url: "https://speakeasy.com/docs/customize-sdks/method-names",
	schema: (z) => z.string().min(1),
});

/**
 * x-speakeasy-group: Group operations together in the generated SDK.
 */
export const xSpeakeasyGroup: ExtensionSchemaMeta = defineExtension({
	name: "x-speakeasy-group",
	scope: ["operation"],
	description: "Group operations together under a namespace in generated SDKs",
	url: "https://speakeasy.com/docs/customize-sdks/namespaces",
	schema: (z) => z.string().min(1),
});

/**
 * x-speakeasy-ignore: Exclude an operation or schema from SDK generation.
 */
export const xSpeakeasyIgnore: ExtensionSchemaMeta = defineExtension({
	name: "x-speakeasy-ignore",
	scope: ["operation", "schema", "parameter"],
	description: "Exclude an operation, schema, or parameter from SDK generation",
	url: "https://speakeasy.com/docs/customize-sdks/ignore",
	schema: (z) => z.boolean(),
});

/**
 * x-speakeasy-retries: Configure retry behavior for an operation.
 */
export const xSpeakeasyRetries: ExtensionSchemaMeta = defineExtension({
	name: "x-speakeasy-retries",
	scope: ["operation", "root"],
	description: "Configure automatic retry behavior for failed requests",
	url: "https://speakeasy.com/docs/customize-sdks/retries",
	schema: (z) =>
		z.looseObject({
			strategy: z.union([z.literal("backoff"), z.literal("none")]).optional(),
			backoff: z
				.looseObject({
					initialInterval: z.number().positive().optional(),
					maxInterval: z.number().positive().optional(),
					maxElapsedTime: z.number().positive().optional(),
					exponent: z.number().positive().optional(),
				})
				
				.optional(),
			statusCodes: z.array(z.union([z.number(), z.string()])).optional(),
			retryConnectionErrors: z.boolean().optional(),
		}),
});

/**
 * x-speakeasy-errors: Define custom error types for an operation.
 */
export const xSpeakeasyErrors: ExtensionSchemaMeta = defineExtension({
	name: "x-speakeasy-errors",
	scope: ["operation"],
	description: "Define which response status codes should be treated as errors",
	url: "https://speakeasy.com/docs/customize-sdks/errors",
	schema: (z) => z.array(z.union([z.number(), z.string()])),
});

/**
 * x-speakeasy-deprecation-message: Custom deprecation message.
 */
export const xSpeakeasyDeprecationMessage: ExtensionSchemaMeta =
	defineExtension({
		name: "x-speakeasy-deprecation-message",
		scope: ["operation", "schema", "parameter"],
		description: "Custom message shown when a deprecated element is used",
		schema: (z) => z.string().min(1),
	});

/**
 * x-speakeasy-pagination: Configure pagination for list operations.
 */
export const xSpeakeasyPagination: ExtensionSchemaMeta = defineExtension({
	name: "x-speakeasy-pagination",
	scope: ["operation"],
	description: "Configure automatic pagination handling for list operations",
	url: "https://speakeasy.com/docs/customize-sdks/pagination",
	schema: (z) =>
		z.looseObject({
			type: z.union([
				z.literal("offsetLimit"),
				z.literal("cursor"),
				z.literal("link"),
			]),
			inputs: z
				.array(
					z.looseObject({
						name: z.string(),
						in: z.union([z.literal("parameters"), z.literal("requestBody")]),
						type: z
							.union([
								z.literal("offset"),
								z.literal("limit"),
								z.literal("cursor"),
								z.literal("page"),
							])
							.optional(),
					}),
				)
				.optional(),
			outputs: z
				.array(
					z.looseObject({
						name: z.string(),
						in: z.union([z.literal("body"), z.literal("headers")]),
						type: z
							.union([
								z.literal("nextCursor"),
								z.literal("nextPage"),
								z.literal("results"),
							])
							.optional(),
					}),
				)
				.optional(),
		}),
});

/**
 * All Speakeasy extensions.
 */
export const speakeasyExtensions: ExtensionSchemaMeta[] = [
	xSpeakeasyEntity,
	xSpeakeasyNameOverride,
	xSpeakeasyGroup,
	xSpeakeasyIgnore,
	xSpeakeasyRetries,
	xSpeakeasyErrors,
	xSpeakeasyDeprecationMessage,
	xSpeakeasyPagination,
];
