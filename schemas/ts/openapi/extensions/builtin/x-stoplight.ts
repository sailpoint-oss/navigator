/**
 * Stoplight OpenAPI Extensions
 *
 * Stoplight platform extensions for API design and documentation.
 * @see https://stoplight.io/
 */

import { defineExtension } from "../index.js";
import type { ExtensionSchemaMeta } from "../types.js";

/**
 * x-stoplight: Stoplight platform metadata.
 */
export const xStoplight: ExtensionSchemaMeta = defineExtension({
	name: "x-stoplight",
	scope: ["root", "operation", "schema"],
	description: "Stoplight platform metadata for API design and documentation",
	url: "https://stoplight.io/",
	schema: (z) =>
		z.looseObject({
			id: z.string().optional(),
			name: z.string().optional(),
			public: z.boolean().optional(),
		}),
});

/**
 * x-internal: Mark elements as internal (shared with other vendors).
 */
export const xStoplightInternal: ExtensionSchemaMeta = defineExtension({
	name: "x-internal",
	scope: ["operation", "schema", "parameter", "pathItem"],
	description:
		"Mark an element as internal, excluding it from public documentation",
	url: "https://stoplight.io/",
	schema: (z) => z.boolean(),
});

/**
 * x-tags: Custom tags for organizing API elements.
 */
export const xTags: ExtensionSchemaMeta = defineExtension({
	name: "x-tags",
	scope: ["schema"],
	description: "Add custom tags to schemas for organization and filtering",
	url: "https://stoplight.io/",
	schema: (z) => z.array(z.string()),
});

/**
 * x-examples: Additional examples for schemas and parameters.
 */
export const xExamples: ExtensionSchemaMeta = defineExtension({
	name: "x-examples",
	scope: ["schema", "parameter", "mediaType"],
	description: "Provide additional examples for schemas and parameters",
	url: "https://stoplight.io/",
	schema: (z) => z.record(z.string(), z.unknown()),
});

/**
 * All Stoplight extensions.
 */
export const stoplightExtensions: ExtensionSchemaMeta[] = [
	xStoplight,
	xStoplightInternal,
	xTags,
	xExamples,
];
