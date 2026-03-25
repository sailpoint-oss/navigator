/**
 * JSON Schema 2020-12 (Core) - Zod schema fragments
 *
 * OpenAPI 3.1+ uses JSON Schema 2020-12 as its base dialect for Schema Objects.
 * This module provides reusable keyword schemas for OAS 3.1/3.2 Schema Objects.
 *
 * Note: This is intentionally the “core keyword surface”; OpenAPI Schema Objects
 * add additional fields (e.g. discriminator, xml) and apply context-specific rules.
 */

import { z } from "zod";

export const JsonSchema2020_12CoreFields = {
	$schema: z
		.string()
		.optional()
		.meta({
			title: "$schema",
			description:
				"Identifies the JSON Schema dialect used to interpret this schema (a URI).",
			examples: ["https://json-schema.org/draft/2020-12/schema"],
		}),
	$id: z
		.string()
		.optional()
		.meta({
			title: "$id",
			description:
				"Sets the base URI for relative references within this schema resource.",
			examples: [
				"https://example.com/schemas/Pet",
				"urn:uuid:123e4567-e89b-12d3-a456-426614174000",
			],
		}),
	$ref: z
		.string()
		.optional()
		.meta({
			title: "$ref",
			description:
				"JSON Schema reference to another schema (by URI), replacing the current schema at this location.",
			examples: ["#/components/schemas/Pet", "https://example.com/schemas/Pet"],
		}),
	$comment: z
		.string()
		.optional()
		.meta({
			title: "$comment",
			description:
				"Non-functional comment intended for schema authors; tooling MAY display it.",
			examples: ["Used for legacy clients; do not remove"],
		}),
	title: z
		.string()
		.optional()
		.meta({
			title: "title",
			description: "Short, human-readable label for this schema.",
			examples: ["Pet", "ErrorResponse"],
		}),
	description: z
		.string()
		.optional()
		.meta({
			title: "description",
			description: "Human-readable explanation of this schema.",
			examples: ["A pet available for adoption."],
		}),
	default: z
		.unknown()
		.optional()
		.meta({
			title: "default",
			description:
				"Default value used as a hint for consumers; does not affect validation.",
			examples: ["unknown", 0, false, null],
		}),
	examples: z
		.array(z.unknown())
		.optional()
		.meta({
			title: "examples",
			description:
				"Example instance values for this schema; does not affect validation.",
			examples: [[{ id: 1, name: "Fido" }]],
		}),
	const: z
		.unknown()
		.optional()
		.meta({
			title: "const",
			description: "Require the instance to be exactly this single value.",
			examples: ["fixed-value", 123, true],
		}),
	enum: z
		.array(z.unknown())
		.optional()
		.meta({
			title: "enum",
			description: "Require the instance to be equal to one of these values.",
			examples: [["small", "medium", "large"]],
		}),
	allOf: z
		.array(z.unknown())
		.optional()
		.meta({
			title: "allOf",
			description: "Require the instance to validate against all subschemas.",
			examples: [[{ $ref: "#/components/schemas/Base" }, { type: "object" }]],
		}),
	anyOf: z
		.array(z.unknown())
		.optional()
		.meta({
			title: "anyOf",
			description:
				"Require the instance to validate against at least one subschema.",
			examples: [[{ type: "string" }, { type: "number" }]],
		}),
	oneOf: z
		.array(z.unknown())
		.optional()
		.meta({
			title: "oneOf",
			description:
				"Require the instance to validate against exactly one subschema.",
			examples: [
				[
					{ $ref: "#/components/schemas/Cat" },
					{ $ref: "#/components/schemas/Dog" },
				],
			],
		}),
	not: z
		.unknown()
		.optional()
		.meta({
			title: "not",
			description:
				"Require the instance to NOT validate against the given subschema.",
			examples: [{ type: "null" }],
		}),
	if: z
		.unknown()
		.optional()
		.meta({
			title: "if",
			description:
				"Conditional validation: if the instance matches this schema, apply `then`, otherwise apply `else` (when present).",
			examples: [{ properties: { country: { const: "US" } } }],
		}),
	// biome-ignore lint/suspicious/noThenProperty: then is a perfectly valid object key, come on biome
	then: z
		.unknown()
		.optional()
		.meta({
			title: "then",
			description: "Subschema applied when `if` matches.",
			examples: [{ required: ["zipCode"] }],
		}),
	else: z
		.unknown()
		.optional()
		.meta({
			title: "else",
			description: "Subschema applied when `if` does not match.",
			examples: [{ not: { required: ["zipCode"] } }],
		}),
	type: z
		.union([z.string(), z.array(z.string())])
		.optional()
		.meta({
			title: "type",
			description:
				"Primitive type constraint (string or array of strings). JSON Schema types: null, boolean, object, array, number, integer, string.",
			examples: ["string", ["string", "null"]],
		}),
	properties: z
		.record(z.string(), z.unknown())
		.optional()
		.meta({
			title: "properties",
			description:
				"Object property schemas keyed by property name (applies when instance is an object).",
			examples: [{ id: { type: "integer" }, name: { type: "string" } }],
		}),
	items: z
		.unknown()
		.optional()
		.meta({
			title: "items",
			description:
				"Array item schema (applies when instance is an array). Can be a schema (or a tuple form in full JSON Schema).",
			examples: [{ type: "string" }, { $ref: "#/components/schemas/Pet" }],
		}),
	additionalProperties: z
		.union([z.unknown(), z.boolean()])
		.optional()
		.meta({
			title: "additionalProperties",
			description:
				"Controls properties not listed in `properties`: boolean to allow/disallow, or a schema to validate them.",
			examples: [true, false, { type: "string" }],
		}),
	required: z
		.array(z.string())
		.optional()
		.meta({
			title: "required",
			description: "List of required property names (for object instances).",
			examples: [["id", "name"]],
		}),
};
