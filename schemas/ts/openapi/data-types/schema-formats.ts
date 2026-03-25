/**
 * Type-Specific Format Schemas for OpenAPI
 *
 * This module defines format enums that are constrained to valid values
 * for each JSON Schema / OpenAPI data type. These provide:
 * - Better LSP autocomplete suggestions
 * - Accurate validation of format values
 * - Rich hover documentation with examples
 *
 * @module engine/schemas/data-types/format-schemas
 */
import { z } from "zod";

// ============================================================================
// String Formats
// ============================================================================

/**
 * Valid format values for string type schemas.
 *
 * Includes JSON Schema formats (RFC draft) and OpenAPI-specific formats.
 * @see https://json-schema.org/understanding-json-schema/reference/string#built-in-formats
 * @see https://spec.openapis.org/oas/v3.1.0#data-types
 */
export const StringFormatSchema = z
	.enum([
		// Date and Time formats (RFC 3339)
		"date-time", // Full date-time: 2024-01-15T09:30:00Z
		"date", // Full date: 2024-01-15
		"time", // Full time: 09:30:00Z
		"duration", // Duration: P3Y6M4DT12H30M5S

		// Email formats (RFC 5321/6531)
		"email", // Email address: user@example.com
		"idn-email", // Internationalized email: 用户@例子.广告

		// Hostname formats (RFC 1123/5890)
		"hostname", // Internet hostname: api.example.com
		"idn-hostname", // Internationalized hostname: api.例え.jp

		// IP Address formats (RFC 2673/4291)
		"ipv4", // IPv4 address: 192.168.1.1
		"ipv6", // IPv6 address: 2001:db8::1

		// URI formats (RFC 3986/3987)
		"uri", // Full URI: https://example.com/path?query=value
		"uri-reference", // URI reference (can be relative): /path/to/resource
		"iri", // Internationalized URI
		"iri-reference", // Internationalized URI reference
		"uri-template", // URI template (RFC 6570): /users/{id}

		// Identifier formats
		"uuid", // UUID: 550e8400-e29b-41d4-a716-446655440000

		// JSON Pointer formats (RFC 6901)
		"json-pointer", // JSON Pointer: /foo/bar/0
		"relative-json-pointer", // Relative JSON Pointer: 1/foo

		// Pattern format
		"regex", // Regular expression pattern

		// OpenAPI-specific formats
		"password", // Hint for UIs to obscure input
		"byte", // Base64-encoded data
		"binary", // Binary data (octet stream)
	])
	.meta({
		title: "format",
		examples: [
			"date-time",
			"date",
			"email",
			"uuid",
			"uri",
			"hostname",
			"ipv4",
			"byte",
			"binary",
			"password",
		],
	})
	.describe(
		"Semantic format hint for string values. Common formats: date-time, email, uuid, uri, hostname, ipv4, byte, binary.",
	);

export type StringFormat = z.infer<typeof StringFormatSchema>;

// ============================================================================
// Number Formats
// ============================================================================

/**
 * Valid format values for number type schemas.
 *
 * @see https://spec.openapis.org/oas/v3.1.0#data-types
 */
export const NumberFormatSchema = z
	.enum([
		"float", // 32-bit single-precision IEEE 754
		"double", // 64-bit double-precision IEEE 754 (default)
	])
	.meta({
		title: "format",
		examples: ["float", "double"],
	})
	.describe(
		"Number precision format. Use 'float' for 32-bit single-precision or 'double' (default) for 64-bit double-precision IEEE 754.",
	);

export type NumberFormat = z.infer<typeof NumberFormatSchema>;

// ============================================================================
// Integer Formats
// ============================================================================

/**
 * Valid format values for integer type schemas.
 *
 * @see https://spec.openapis.org/oas/v3.1.0#data-types
 */
export const IntegerFormatSchema = z
	.enum([
		"int32", // Signed 32-bit integer (-2,147,483,648 to 2,147,483,647)
		"int64", // Signed 64-bit integer (-9,223,372,036,854,775,808 to 9,223,372,036,854,775,807)
	])
	.meta({
		title: "format",
		examples: ["int32", "int64"],
	})
	.describe(
		"Integer size format. Use 'int32' for signed 32-bit (-2B to 2B) or 'int64' for signed 64-bit integers.",
	);

export type IntegerFormat = z.infer<typeof IntegerFormatSchema>;

// ============================================================================
// Combined Schema Format (for generic schema contexts)
// ============================================================================

/**
 * Union of all valid format values across all types.
 * Use this when the schema type is not yet known.
 */
export const AnyFormatSchema = z
	.union([StringFormatSchema, NumberFormatSchema, IntegerFormatSchema])
	.meta({
		title: "format",
		examples: ["date-time", "email", "uuid", "int32", "int64", "float"],
	})
	.describe(
		"Schema format hint. Valid values depend on the schema type: string (date-time, email, uuid, uri, etc.), integer (int32, int64), number (float, double).",
	);

export type AnyFormat = z.infer<typeof AnyFormatSchema>;

// ============================================================================
// Parameter Style Schemas
// ============================================================================

/**
 * Valid style values for path parameters.
 */
export const PathParameterStyleSchema = z
	.enum([
		"simple", // Default. Comma-separated values: /users/3,4,5
		"matrix", // Semicolon-prefixed: /users;id=3;id=4;id=5
		"label", // Dot-prefixed: /users/.3.4.5
	])
	.meta({
		title: "style",
		examples: ["simple", "matrix", "label"],
	})
	.describe(
		"Serialization style for path parameters. 'simple' (default): comma-separated. 'matrix': semicolon-prefixed. 'label': dot-prefixed.",
	);

/**
 * Valid style values for query parameters.
 */
export const QueryParameterStyleSchema = z
	.enum([
		"form", // Default. Ampersand-separated: ?id=3&id=4&id=5
		"spaceDelimited", // Space-separated (encoded): ?id=3%204%205
		"pipeDelimited", // Pipe-separated: ?id=3|4|5
		"deepObject", // Nested objects: ?filter[status]=active&filter[limit]=10
	])
	.meta({
		title: "style",
		examples: ["form", "spaceDelimited", "pipeDelimited", "deepObject"],
	})
	.describe(
		"Serialization style for query parameters. 'form' (default): ampersand-separated. 'deepObject': for nested objects.",
	);

/**
 * Valid style values for header parameters.
 */
export const HeaderParameterStyleSchema = z
	.literal("simple")
	.meta({
		title: "style",
		examples: ["simple"],
	})
	.describe(
		"Serialization style for header parameters. Only 'simple' (comma-separated) is valid.",
	);

/**
 * Valid style values for cookie parameters.
 */
export const CookieParameterStyleSchema = z
	.literal("form")
	.meta({
		title: "style",
		examples: ["form"],
	})
	.describe("Serialization style for cookie parameters. Only 'form' is valid.");

/**
 * All parameter styles combined.
 */
export const ParameterStyleSchema = z
	.enum([
		"matrix",
		"label",
		"form",
		"simple",
		"spaceDelimited",
		"pipeDelimited",
		"deepObject",
	])
	.meta({
		title: "style",
		examples: ["simple", "form", "matrix", "label", "deepObject"],
	})
	.describe(
		"Serialization style for parameter values. Valid options depend on parameter location (in).",
	);

export type ParameterStyle = z.infer<typeof ParameterStyleSchema>;

// ============================================================================
// Encoding Style Schema
// ============================================================================

/**
 * Valid style values for encoding objects in request bodies.
 */
export const EncodingStyleSchema = z
	.enum([
		"form", // Default for application/x-www-form-urlencoded
		"spaceDelimited", // Space-separated values
		"pipeDelimited", // Pipe-separated values
		"deepObject", // Nested object serialization
	])
	.meta({
		title: "style",
		examples: ["form", "spaceDelimited", "pipeDelimited", "deepObject"],
	})
	.describe(
		"Serialization style for encoding properties in multipart or form-urlencoded request bodies.",
	);

export type EncodingStyle = z.infer<typeof EncodingStyleSchema>;

// ============================================================================
// Parameter Location Schema
// ============================================================================

/**
 * Valid locations for parameters.
 */
export const ParameterLocationSchema = z
	.enum([
		"query", // Query string parameter: ?name=value
		"header", // HTTP header: X-Custom-Header: value
		"path", // Path parameter: /users/{id}
		"cookie", // Cookie parameter
	])
	.meta({
		title: "in",
		examples: ["query", "path", "header", "cookie"],
	})
	.describe(
		"Location of the parameter. 'path' parameters are required by default.",
	);

export type ParameterLocation = z.infer<typeof ParameterLocationSchema>;

// ============================================================================
// Security Scheme Type Schema
// ============================================================================

/**
 * Valid security scheme types.
 */
export const SecuritySchemeTypeSchema = z
	.enum([
		"apiKey", // API key in header, query, or cookie
		"http", // HTTP authentication (Basic, Bearer, etc.)
		"oauth2", // OAuth 2.0 flows
		"openIdConnect", // OpenID Connect Discovery
		"mutualTLS", // Mutual TLS (3.1+)
	])
	.meta({
		title: "type",
		examples: ["apiKey", "http", "oauth2", "openIdConnect"],
	})
	.describe(
		"Type of security scheme. 'apiKey': API key. 'http': HTTP auth (Bearer, Basic). 'oauth2': OAuth 2.0. 'openIdConnect': OIDC.",
	);

export type SecuritySchemeType = z.infer<typeof SecuritySchemeTypeSchema>;

// ============================================================================
// API Key Location Schema
// ============================================================================

/**
 * Valid locations for API key security schemes.
 */
export const ApiKeyLocationSchema = z
	.enum([
		"query", // API key in query string
		"header", // API key in header (most common)
		"cookie", // API key in cookie
	])
	.meta({
		title: "in",
		examples: ["header", "query", "cookie"],
	})
	.describe(
		"Location of the API key. 'header' is most common (e.g., X-API-Key header).",
	);

export type ApiKeyLocation = z.infer<typeof ApiKeyLocationSchema>;

// ============================================================================
// HTTP Auth Scheme Schema
// ============================================================================

/**
 * Common HTTP authentication schemes.
 * This is not exhaustive - custom schemes are allowed.
 */
export const HttpAuthSchemeSchema = z
	.string()
	.meta({
		title: "scheme",
		examples: [
			"bearer",
			"basic",
			"digest",
			"hoba",
			"mutual",
			"negotiate",
			"oauth",
			"scram-sha-1",
			"scram-sha-256",
			"vapid",
		],
	})
	.describe(
		"HTTP authentication scheme name (IANA registered). Common: 'bearer' (JWT/tokens), 'basic' (username:password).",
	);

export type HttpAuthScheme = z.infer<typeof HttpAuthSchemeSchema>;
