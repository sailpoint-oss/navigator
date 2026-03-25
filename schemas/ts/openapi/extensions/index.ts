/**
 * OpenAPI Extension Schema Registry
 *
 * Provides the defineExtension helper and extension schema compilation.
 * Extensions are validated when present in OpenAPI documents.
 */

import { z } from "zod";
import type {
	CompiledExtension,
	ExtensionRegistry,
	ExtensionSchemaMeta,
	ExtensionScope,
	ExtensionValidationError,
} from "./types.js";

// Re-export types
export type {
	CompiledExtension,
	ExtensionRegistry,
	ExtensionSchemaMeta,
	ExtensionScope,
	ExtensionValidationError,
} from "./types.js";

/**
 * Define an OpenAPI extension schema.
 *
 * Extensions are validated when they appear at the specified scopes.
 * To mark extensions as required, use the `required` array in the telescope config.
 *
 * @param meta - Extension metadata including name, scope, description, and schema
 * @returns The extension metadata (validated)
 *
 * @example
 * ```typescript
 * import { defineExtension } from "telescope/extensions";
 *
 * export default defineExtension({
 *   name: "x-company-auth",
 *   scope: ["operation"],
 *   description: "Authentication level required for this operation",
 *   schema: (z) => z.union([
 *     z.literal("public"),
 *     z.literal("authenticated"),
 *     z.literal("admin")
 *   ])
 * });
 * ```
 */
export function defineExtension(
	meta: ExtensionSchemaMeta,
): ExtensionSchemaMeta {
	// Validate extension name starts with "x-"
	if (!meta.name.startsWith("x-")) {
		throw new Error(
			`Extension name "${meta.name}" must start with "x-". ` +
				`Example: "x-${meta.name}"`,
		);
	}

	// Validate reserved prefixes
	const lowerName = meta.name.toLowerCase();
	if (lowerName.startsWith("x-oai-") || lowerName.startsWith("x-oas-")) {
		throw new Error(
			`Extension name "${meta.name}" uses a reserved prefix. ` +
				`Names starting with "x-oai-" or "x-oas-" are reserved by OpenAPI.`,
		);
	}

	// Validate scope is not empty
	if (!meta.scope || meta.scope.length === 0) {
		throw new Error(
			`Extension "${meta.name}" must specify at least one scope.`,
		);
	}

	return meta;
}

/**
 * Compile an extension definition into a ready-to-use compiled extension.
 *
 * @param meta - Extension metadata
 * @returns Compiled extension with built schema
 */
export function compileExtension(meta: ExtensionSchemaMeta): CompiledExtension {
	return {
		name: meta.name,
		scope: meta.scope,
		description: meta.description,
		url: meta.url,
		compiledSchema: meta.schema(z),
	};
}

/**
 * Create an empty extension registry.
 */
export function createExtensionRegistry(): ExtensionRegistry {
	return {
		root: [],
		info: [],
		paths: [],
		pathItem: [],
		operation: [],
		parameter: [],
		schema: [],
		response: [],
		requestBody: [],
		components: [],
		header: [],
		mediaType: [],
		securityScheme: [],
		tag: [],
		server: [],
		any: [],
	};
}

/**
 * Register a compiled extension into the registry.
 *
 * @param registry - The extension registry
 * @param extension - The compiled extension to register
 */
export function registerExtension(
	registry: ExtensionRegistry,
	extension: CompiledExtension,
): void {
	for (const scope of extension.scope) {
		registry[scope].push(extension);
	}
}

/**
 * Get all extensions applicable to a specific scope.
 * Returns extensions registered for that scope plus any registered for "any".
 *
 * @param registry - The extension registry
 * @param scope - The scope to get extensions for
 * @returns Array of compiled extensions
 */
export function getExtensionsForScope(
	registry: ExtensionRegistry,
	scope: Exclude<ExtensionScope, "any">,
): CompiledExtension[] {
	return [...registry[scope], ...registry.any];
}

/**
 * Build a complete extension registry from an array of extension definitions.
 *
 * @param extensions - Array of extension metadata
 * @returns Populated extension registry
 */
export function buildExtensionRegistry(
	extensions: ExtensionSchemaMeta[],
): ExtensionRegistry {
	const registry = createExtensionRegistry();

	for (const meta of extensions) {
		const compiled = compileExtension(meta);
		registerExtension(registry, compiled);
	}

	return registry;
}

/**
 * Validate an extension value against a compiled extension schema.
 *
 * @param extension - The compiled extension
 * @param value - The value to validate
 * @returns Validation result with success flag and any errors
 */
export function validateExtensionValue(
	extension: CompiledExtension,
	value: unknown,
): { success: boolean; errors?: ExtensionValidationError[] } {
	const result = extension.compiledSchema.safeParse(value);

	if (result.success) {
		return { success: true };
	}

	// Convert Zod errors to ExtensionValidationError format
	const errors: ExtensionValidationError[] = [];
	for (const issue of result.error.issues) {
		errors.push({
			path: `/${issue.path.join("/")}`,
			message: issue.message,
			value: value,
		});
	}

	return { success: false, errors };
}

/**
 * Diagnostic for a missing required extension.
 */
export interface ExtensionDiagnostic {
	extensionName: string;
	scope: ExtensionScope;
	message: string;
	description: string;
	url?: string;
}

/**
 * Check for missing required extensions at a specific scope.
 *
 * @param registry - The extension registry
 * @param scope - The scope being checked
 * @param obj - The object at that scope
 * @param requiredExtensions - Array of extension names that are required
 * @returns Array of diagnostics for missing required extensions
 */
export function checkRequiredExtensions(
	registry: ExtensionRegistry,
	scope: Exclude<ExtensionScope, "any">,
	obj: Record<string, unknown>,
	requiredExtensions: string[] = [],
): ExtensionDiagnostic[] {
	const diagnostics: ExtensionDiagnostic[] = [];
	const extensions = getExtensionsForScope(registry, scope);

	// Create a set of required extension names for quick lookup
	const requiredSet = new Set(requiredExtensions);

	for (const extension of extensions) {
		// Check if this extension is in the required list AND is missing
		if (requiredSet.has(extension.name) && !(extension.name in obj)) {
			// Warn if required extension has scope "any" (insane configuration)
			if (extension.scope.includes("any")) {
				console.warn(
					`⚠️  Extension "${extension.name}" is marked as required but has scope "any". ` +
						`This requires the extension on EVERY OpenAPI object of this type, which is almost certainly unintended.`,
				);
			}

			diagnostics.push({
				extensionName: extension.name,
				scope,
				message: `Required extension "${extension.name}" is missing`,
				description: extension.description,
				url: extension.url,
			});
		}
	}

	return diagnostics;
}

/**
 * Validate all extensions present in an object at a specific scope.
 *
 * @param registry - The extension registry
 * @param scope - The scope being validated
 * @param obj - The object at that scope
 * @param requiredExtensions - Array of extension names that are required
 * @returns Array of validation diagnostics
 */
export function validateExtensionsAtScope(
	registry: ExtensionRegistry,
	scope: Exclude<ExtensionScope, "any">,
	obj: Record<string, unknown>,
	requiredExtensions: string[] = [],
): ExtensionDiagnostic[] {
	const diagnostics: ExtensionDiagnostic[] = [];
	const extensions = getExtensionsForScope(registry, scope);

	// Create a map for quick lookup
	const extensionMap = new Map(extensions.map((ext) => [ext.name, ext]));

	// Check each x-* key in the object
	for (const [key, value] of Object.entries(obj)) {
		if (!key.startsWith("x-")) continue;

		const extension = extensionMap.get(key);
		if (extension) {
			// Validate against the schema
			const result = validateExtensionValue(extension, value);
			if (!result.success && result.errors) {
				for (const error of result.errors) {
					diagnostics.push({
						extensionName: extension.name,
						scope,
						message: `Invalid value for "${extension.name}": ${error.message}`,
						description: extension.description,
						url: extension.url,
					});
				}
			}
		}
		// Unknown extensions are allowed (just not validated)
	}

	// Also check for required extensions
	diagnostics.push(
		...checkRequiredExtensions(registry, scope, obj, requiredExtensions),
	);

	return diagnostics;
}
