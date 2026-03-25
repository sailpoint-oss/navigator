/**
 * Zod schema for Telescope configuration files.
 * Config file location: .telescope/config.yaml
 */

import { z } from "zod";

/**
 * Severity levels for rule diagnostics.
 * - off: Disable the rule entirely
 * - error: Critical issues that must be fixed
 * - warn: Best practice violations that should be addressed
 * - info: Informational suggestions for improvement
 * - hint: Minor style suggestions
 */
export const SeveritySchema = z
	.union([
		z.literal("off"),
		z.literal("error"),
		z.literal("warn"),
		z.literal("info"),
		z.literal("hint"),
	])
	.meta({ title: "Severity" })
	.describe("Rule severity level");

export type Severity = "off" | "error" | "warn" | "info" | "hint";

// Schema for OpenAPI rule configuration
export const OpenAPIRuleConfigSchema = z
	.looseObject({
		rule: z
			.string()
			.meta({ title: "rule" })
			.describe(
				"Path to custom OpenAPI rule file. Resolved from .telescope/rules/ or workspace root.",
			),
		patterns: z
			.array(z.string())
			.meta({ title: "patterns" })
			.describe(
				"Glob patterns to match files. If omitted, uses OpenAPI.patterns. Use ! prefix to exclude.",
			)
			.optional(),
	})
	
	.meta({ title: "OpenAPIRuleConfig" })
	.describe(
		"Custom OpenAPI rule. Requires rule (file path). Optional pattern (array of glob patterns). Invalid properties cause errors.",
	);

export type OpenAPIRuleConfig = z.infer<typeof OpenAPIRuleConfigSchema>;

// Schema for AdditionalValidation labeled group
export const AdditionalValidationGroupSchema = z
	.looseObject({
		patterns: z
			.array(z.string())
			.meta({ title: "patterns" })
			.describe(
				"Glob patterns for files. Use ! prefix to exclude. Default: all YAML/JSON files.",
			)
			.optional(),
		schemas: z
			.array(
				z.looseObject({
					schema: z
						.string()
						.meta({ title: "schema" })
						.describe(
							"Path to schema file. Supports JSON Schema (.json) and Zod Schema (.ts). Resolved from .telescope/schemas/ or workspace root.",
						),
					patterns: z
						.array(z.string())
						.meta({ title: "patterns" })
						.describe(
							"Glob patterns for files. If omitted, uses group patterns. Valid properties: schema (required), pattern (optional).",
						)
						.optional(),
				}),
			)
			.meta({ title: "schemas" })
			.describe(
				"JSON schema validations. Each requires schema (file path) and optionally pattern (array of glob patterns).",
			)
			.optional(),
		rules: z
			.array(
				z.looseObject({
					rule: z
						.string()
						.meta({ title: "rule" })
						.describe(
							"Path to custom generic rule file. Resolved from .telescope/rules/ or workspace root.",
						),
					patterns: z
						.array(z.string())
						.meta({ title: "patterns" })
						.describe(
							"Glob patterns for files. If omitted, uses group patterns. Valid properties: rule (required), pattern (optional).",
						)
						.optional(),
				}),
			)
			.meta({ title: "rules" })
			.describe(
				"Custom generic rules. Each requires rule (file path) and optionally pattern (array of glob patterns).",
			)
			.optional(),
	})
	
	.meta({ title: "AdditionalValidationGroup" })
	.describe(
		"Additional validation group. Fields: patterns (glob array), schemas (schema array), rules (rule array). All optional.",
	);

export type AdditionalValidationGroup = z.infer<
	typeof AdditionalValidationGroupSchema
>;

// Telescope configuration schema
export const TelescopeConfigSchema = z
	.looseObject({
		openapi: z
			.looseObject({
				targetVersion: z
					.union([z.literal("3.0"), z.literal("3.1"), z.literal("3.2")])
					.optional()
					.meta({ title: "targetVersion" })
					.describe(
						"OpenAPI version to use for fragment file validation. When set, standalone files (schemas, path items, etc.) are validated against the specified version's rules. Defaults to 3.1.",
					),
				patterns: z
					.array(z.string())
					.meta({ title: "patterns" })
					.describe(
						"Glob patterns for files. Use ! prefix to exclude. Example: **/*.yaml or !**/node_modules/**",
					)
					.optional(),
				sailpoint: z
					.boolean()
					.meta({ title: "sailpoint" })
					.describe(
						"Enable SailPoint-specific rules. When true, adds rules for x-sailpoint-api, user levels, security requirements, etc.",
					)
					.optional(),
				rules: z
					.array(OpenAPIRuleConfigSchema)
					.meta({ title: "rules" })
					.describe(
						"Custom OpenAPI rules. Each rule requires rule (path to file) and optionally pattern.",
					)
					.optional(),
				rulesOverrides: z
					.record(z.string(), SeveritySchema)
					.meta({ title: "rulesOverrides" })
					.describe(
						"Rule overrides. Map rule IDs to severity (error, warn, info, hint, off).",
					)
					.optional(),
				extensions: z
					.looseObject({
						schemas: z
							.array(z.string())
							.meta({ title: "schemas" })
							.describe(
								"Custom OpenAPI extension schema files. Paths are resolved from .telescope/extensions/ by default.",
							)
							.optional(),
						required: z
							.array(z.string())
							.meta({ title: "required" })
							.describe(
								"Extension names that are required. Can include both builtin and custom extensions. Example: ['x-company-auth', 'x-speakeasy-entity']",
							)
							.optional(),
					})
					
					.meta({ title: "extensions" })
					.describe(
						"OpenAPI extension configuration. schemas: custom extension files, required: extension names that must be present.",
					)
					.optional(),
			})
			
			.meta({ title: "openapi" })
			.describe(
				"OpenAPI configuration. All fields are optional.\n\nFields: patterns, sailpoint, rules, rulesOverrides, extensions",
			)
			.optional(),
		additionalValidation: z
			.record(z.string(), AdditionalValidationGroupSchema)
			.meta({ title: "additionalValidation" })
			.describe(
				"Additional validation groups (labeled). Each group can have patterns, schemas, and rules.",
			)
			.optional(),
		lsp: z
			.looseObject({
				schemaValidation: z
					.looseObject({
						mode: z
							.literal("go")
							.meta({ title: "mode" })
							.describe(
								"Schema validation backend. Telescope uses the Go validator.",
							)
							.optional(),
					})
					.meta({ title: "schemaValidation" })
					.describe("LSP schema validation settings.")
					.optional(),
			})
			.meta({ title: "lsp" })
			.describe("LSP runtime settings.")
			.optional(),
	})
	
	.meta({ title: "TelescopeConfig" })
	.describe(
		"Telescope configuration file. Location: .telescope/config.yaml. Top-level keys: OpenAPI, AdditionalValidation. All optional.",
	);

export type TelescopeConfig = z.infer<typeof TelescopeConfigSchema>;
