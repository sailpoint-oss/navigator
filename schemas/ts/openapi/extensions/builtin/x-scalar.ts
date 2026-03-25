/**
 * Scalar OpenAPI Extensions
 *
 * Scalar extensions for API documentation and client configuration.
 * @see https://github.com/scalar/scalar
 */

import { defineExtension } from "../index.js";
import type { ExtensionSchemaMeta } from "../types.js";

/**
 * x-scalar-environments: Define multiple named environments for API clients.
 */
export const xScalarEnvironments: ExtensionSchemaMeta = defineExtension({
	name: "x-scalar-environments",
	scope: ["root"],
	description:
		"Define multiple named environments (production, staging, etc.) with their own variables",
	url: "https://github.com/scalar/scalar/blob/main/documentation/openapi.md",
	schema: (z) =>
		z.record(
			z.string(),
			z.looseObject({
				description: z.string().optional(),
				color: z.string().optional(),
				variables: z
					.record(
						z.string(),
						z.union([
							z.string(),
							z.looseObject({
								description: z.string().optional(),
								default: z.string(),
							}),
						]),
					)
					.optional(),
			}),
		),
});

/**
 * x-scalar-active-environment: Specify the default active environment.
 */
export const xScalarActiveEnvironment: ExtensionSchemaMeta = defineExtension({
	name: "x-scalar-active-environment",
	scope: ["root"],
	description:
		"Specify which environment from x-scalar-environments should be active by default",
	url: "https://github.com/scalar/scalar/blob/main/documentation/openapi.md",
	schema: (z) => z.string().min(1),
});

/**
 * x-scalar-sdk-installation: Custom SDK installation instructions.
 */
export const xScalarSdkInstallation: ExtensionSchemaMeta = defineExtension({
	name: "x-scalar-sdk-installation",
	scope: ["info"],
	description:
		"Define custom SDK installation instructions for different programming languages",
	url: "https://github.com/scalar/scalar/blob/main/documentation/openapi.md",
	schema: (z) =>
		z.array(
			z.looseObject({
				lang: z.string(),
				description: z.string().optional(),
				source: z.string().optional(),
			}),
		),
});

/**
 * x-scalar-stability: Indicate endpoint stability level.
 */
export const xScalarStability: ExtensionSchemaMeta = defineExtension({
	name: "x-scalar-stability",
	scope: ["operation"],
	description:
		"Indicate the stability level of an endpoint (stable, experimental, deprecated)",
	url: "https://github.com/scalar/scalar/blob/main/documentation/openapi.md",
	schema: (z) =>
		z.union([
			z.literal("stable"),
			z.literal("experimental"),
			z.literal("deprecated"),
			z.literal("alpha"),
			z.literal("beta"),
		]),
});

/**
 * x-displayName: Override display name for tags (shared with Redocly).
 */
export const xScalarDisplayName: ExtensionSchemaMeta = defineExtension({
	name: "x-displayName",
	scope: ["tag"],
	description:
		"Override the display name of a tag for more user-friendly documentation",
	url: "https://github.com/scalar/scalar/blob/main/documentation/openapi.md",
	schema: (z) => z.string().min(1),
});

/**
 * x-codeSamples: Add code samples to operations (shared with Redocly).
 */
export const xScalarCodeSamples: ExtensionSchemaMeta = defineExtension({
	name: "x-codeSamples",
	scope: ["operation"],
	description:
		"Provide code samples for various programming languages and SDKs",
	url: "https://github.com/scalar/scalar/blob/main/documentation/openapi.md",
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
 * x-scalar-icon: Custom icon for the API reference.
 */
export const xScalarIcon: ExtensionSchemaMeta = defineExtension({
	name: "x-scalar-icon",
	scope: ["root", "tag"],
	description: "Specify a custom icon for the API reference or specific tags",
	url: "https://github.com/scalar/scalar/blob/main/packages/openapi-types/README.md",
	schema: (z) => z.string().min(1),
});

/**
 * All Scalar extensions.
 */
export const scalarExtensions: ExtensionSchemaMeta[] = [
	xScalarEnvironments,
	xScalarActiveEnvironment,
	xScalarSdkInstallation,
	xScalarStability,
	xScalarDisplayName,
	xScalarCodeSamples,
	xScalarIcon,
];
