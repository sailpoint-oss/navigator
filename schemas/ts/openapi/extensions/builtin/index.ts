/**
 * Builtin OpenAPI Extensions
 *
 * All builtin extension schemas from popular vendors.
 * These are automatically available and validated when present in OpenAPI documents.
 */

import type { ExtensionSchemaMeta } from "../types.js";
import { redoclyExtensions } from "./x-redocly.js";
import { scalarExtensions } from "./x-scalar.js";
import { speakeasyExtensions } from "./x-speakeasy.js";
import { stoplightExtensions } from "./x-stoplight.js";

// Re-export individual extension sets
export { redoclyExtensions } from "./x-redocly.js";
export { scalarExtensions } from "./x-scalar.js";
export { speakeasyExtensions } from "./x-speakeasy.js";
export { stoplightExtensions } from "./x-stoplight.js";

/**
 * All builtin extensions from all vendors.
 * This is the complete set of extensions that are automatically validated.
 */
export const builtinExtensions: ExtensionSchemaMeta[] = [
	...speakeasyExtensions,
	...redoclyExtensions,
	...scalarExtensions,
	...stoplightExtensions,
];

/**
 * Map of extension names to their schemas for quick lookup.
 */
export const builtinExtensionsMap = new Map<string, ExtensionSchemaMeta>(
	builtinExtensions.map((ext) => [ext.name, ext]),
);
