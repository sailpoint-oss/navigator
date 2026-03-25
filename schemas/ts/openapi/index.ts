/**
 * Central export point for all OpenAPI 3.x schema definitions.
 *
 * These schemas are built using Zod with .meta() and .describe()
 * for proper JSON Schema generation via z.toJSONSchema().
 *
 * Supports OpenAPI 3.0, 3.1, and 3.2 specifications with version-specific schemas.
 */

import type { z } from "zod";

// =============================================================================
// Telescope Configuration Schema
// =============================================================================
export {
	type AdditionalValidationGroup,
	AdditionalValidationGroupSchema,
	type OpenAPIRuleConfig,
	OpenAPIRuleConfigSchema,
	type Severity,
	SeveritySchema,
	type TelescopeConfig,
	TelescopeConfigSchema,
} from "./config-schema.js";
// =============================================================================
// OpenAPI 2.0 (Swagger 2.0) Schemas
// =============================================================================
export {
	Contact2Schema,
	ExternalDocs2Schema,
	Info2Schema,
	OpenAPI2Schema,
	Operation2Schema,
	Parameter2Schema,
	PathItem2Schema,
	Paths2Schema,
	Response2Schema,
	Responses2Schema,
	SchemaObject2Schema,
	Tag2Schema,
} from "./openapi-2.0-module.js";

export type Contact2 = z.infer<
	typeof import("./openapi-2.0-module.js").Contact2Schema
>;
export type ExternalDocs2 = z.infer<
	typeof import("./openapi-2.0-module.js").ExternalDocs2Schema
>;
export type Info2 = z.infer<
	typeof import("./openapi-2.0-module.js").Info2Schema
>;
export type OpenAPI2 = z.infer<
	typeof import("./openapi-2.0-module.js").OpenAPI2Schema
>;
export type Operation2 = z.infer<
	typeof import("./openapi-2.0-module.js").Operation2Schema
>;
export type Parameter2 = z.infer<
	typeof import("./openapi-2.0-module.js").Parameter2Schema
>;
export type PathItem2 = z.infer<
	typeof import("./openapi-2.0-module.js").PathItem2Schema
>;
export type Paths2 = z.infer<
	typeof import("./openapi-2.0-module.js").Paths2Schema
>;
export type Response2 = z.infer<
	typeof import("./openapi-2.0-module.js").Response2Schema
>;
export type Responses2 = z.infer<
	typeof import("./openapi-2.0-module.js").Responses2Schema
>;
export type SchemaObject2 = z.infer<
	typeof import("./openapi-2.0-module.js").SchemaObject2Schema
>;
export type Tag2 = z.infer<typeof import("./openapi-2.0-module.js").Tag2Schema>;

// =============================================================================
// OpenAPI 3.0 Schemas
// =============================================================================
export {
	ArraySchema30,
	BooleanSchema30,
	Callback30Schema,
	Components30Schema,
	// Types
	// Individual schemas
	Contact30Schema,
	Discriminator30Schema,
	Encoding30Schema,
	Example30Schema,
	ExternalDocumentation30Schema,
	Header30Schema,
	Info30Schema,
	IntegerSchema30,
	License30Schema,
	Link30Schema,
	MediaType30Schema,
	NumberSchema30,
	OAuthFlow30Schema,
	OAuthFlows30Schema,
	ObjectSchema30,
	OpenAPI3Schema,
	OpenAPI30Schema,
	Operation30Schema,
	Parameter30Schema,
	PathItem30Schema,
	Paths30Schema,
	RequestBody30Schema,
	Response30Schema,
	Responses30Schema,
	SchemaObject30Schema,
	SecurityRequirement30Schema,
	SecurityScheme30Schema,
	Server30Schema,
	ServerVariable30Schema,
	StringSchema30,
	Tag30Schema,
	XML30Schema,
} from "./openapi-3.0-module.js";

export type Callback30 = z.infer<
	typeof import("./openapi-3.0-module.js").Callback30Schema
>;
export type Components30 = z.infer<
	typeof import("./openapi-3.0-module.js").Components30Schema
>;
export type Contact30 = z.infer<
	typeof import("./openapi-3.0-module.js").Contact30Schema
>;
export type Discriminator30 = z.infer<
	typeof import("./openapi-3.0-module.js").Discriminator30Schema
>;
export type Encoding30 = z.infer<
	typeof import("./openapi-3.0-module.js").Encoding30Schema
>;
export type Example30 = z.infer<
	typeof import("./openapi-3.0-module.js").Example30Schema
>;
export type ExternalDocumentation30 = z.infer<
	typeof import("./openapi-3.0-module.js").ExternalDocumentation30Schema
>;
export type Header30 = z.infer<
	typeof import("./openapi-3.0-module.js").Header30Schema
>;
export type Info30 = z.infer<
	typeof import("./openapi-3.0-module.js").Info30Schema
>;
export type License30 = z.infer<
	typeof import("./openapi-3.0-module.js").License30Schema
>;
export type Link30 = z.infer<
	typeof import("./openapi-3.0-module.js").Link30Schema
>;
export type MediaType30 = z.infer<
	typeof import("./openapi-3.0-module.js").MediaType30Schema
>;
export type OAuthFlow30 = z.infer<
	typeof import("./openapi-3.0-module.js").OAuthFlow30Schema
>;
export type OAuthFlows30 = z.infer<
	typeof import("./openapi-3.0-module.js").OAuthFlows30Schema
>;
export type OpenAPI3 = z.infer<
	typeof import("./openapi-3.0-module.js").OpenAPI3Schema
>;
export type OpenAPI30 = z.infer<
	typeof import("./openapi-3.0-module.js").OpenAPI30Schema
>;
export type Operation30 = z.infer<
	typeof import("./openapi-3.0-module.js").Operation30Schema
>;
export type Parameter30 = z.infer<
	typeof import("./openapi-3.0-module.js").Parameter30Schema
>;
export type PathItem30 = z.infer<
	typeof import("./openapi-3.0-module.js").PathItem30Schema
>;
export type Paths30 = z.infer<
	typeof import("./openapi-3.0-module.js").Paths30Schema
>;
export type RequestBody30 = z.infer<
	typeof import("./openapi-3.0-module.js").RequestBody30Schema
>;
export type Response30 = z.infer<
	typeof import("./openapi-3.0-module.js").Response30Schema
>;
export type Responses30 = z.infer<
	typeof import("./openapi-3.0-module.js").Responses30Schema
>;
export type SchemaObject30 = z.infer<
	typeof import("./openapi-3.0-module.js").SchemaObject30Schema
>;
export type SecurityRequirement30 = z.infer<
	typeof import("./openapi-3.0-module.js").SecurityRequirement30Schema
>;
export type SecurityScheme30 = z.infer<
	typeof import("./openapi-3.0-module.js").SecurityScheme30Schema
>;
export type Server30 = z.infer<
	typeof import("./openapi-3.0-module.js").Server30Schema
>;
export type ServerVariable30 = z.infer<
	typeof import("./openapi-3.0-module.js").ServerVariable30Schema
>;
export type Tag30 = z.infer<
	typeof import("./openapi-3.0-module.js").Tag30Schema
>;
export type XML30 = z.infer<
	typeof import("./openapi-3.0-module.js").XML30Schema
>;
// =============================================================================
// OpenAPI 3.1 Schemas
// =============================================================================
export {
	ArraySchema31,
	BooleanSchema31,
	Callback31Schema,
	Components31Schema,
	// Types
	// Individual schemas
	Contact31Schema,
	Discriminator31Schema,
	Encoding31Schema,
	Example31Schema,
	ExternalDocumentation31Schema,
	Header31Schema,
	Info31Schema,
	IntegerSchema31,
	License31Schema,
	Link31Schema,
	MediaType31Schema,
	NullSchema31,
	NumberSchema31,
	OAuthFlow31Schema,
	OAuthFlows31Schema,
	ObjectSchema31,
	OpenAPI31Schema,
	Operation31Schema,
	Parameter31Schema,
	PathItem31Schema,
	Paths31Schema,
	RequestBody31Schema,
	Response31Schema,
	Responses31Schema,
	SchemaObject31Schema,
	SecurityRequirement31Schema,
	SecurityScheme31Schema,
	Server31Schema,
	ServerVariable31Schema,
	StringSchema31,
	Tag31Schema,
	XML31Schema,
} from "./openapi-3.1-module.js";

export type Callback31 = z.infer<
	typeof import("./openapi-3.1-module.js").Callback31Schema
>;
export type Components31 = z.infer<
	typeof import("./openapi-3.1-module.js").Components31Schema
>;
export type Contact31 = z.infer<
	typeof import("./openapi-3.1-module.js").Contact31Schema
>;
export type Discriminator31 = z.infer<
	typeof import("./openapi-3.1-module.js").Discriminator31Schema
>;
export type Encoding31 = z.infer<
	typeof import("./openapi-3.1-module.js").Encoding31Schema
>;
export type Example31 = z.infer<
	typeof import("./openapi-3.1-module.js").Example31Schema
>;
export type ExternalDocumentation31 = z.infer<
	typeof import("./openapi-3.1-module.js").ExternalDocumentation31Schema
>;
export type Header31 = z.infer<
	typeof import("./openapi-3.1-module.js").Header31Schema
>;
export type Info31 = z.infer<
	typeof import("./openapi-3.1-module.js").Info31Schema
>;
export type License31 = z.infer<
	typeof import("./openapi-3.1-module.js").License31Schema
>;
export type Link31 = z.infer<
	typeof import("./openapi-3.1-module.js").Link31Schema
>;
export type MediaType31 = z.infer<
	typeof import("./openapi-3.1-module.js").MediaType31Schema
>;
export type OAuthFlow31 = z.infer<
	typeof import("./openapi-3.1-module.js").OAuthFlow31Schema
>;
export type OAuthFlows31 = z.infer<
	typeof import("./openapi-3.1-module.js").OAuthFlows31Schema
>;
export type OpenAPI31 = z.infer<
	typeof import("./openapi-3.1-module.js").OpenAPI31Schema
>;
export type Operation31 = z.infer<
	typeof import("./openapi-3.1-module.js").Operation31Schema
>;
export type Parameter31 = z.infer<
	typeof import("./openapi-3.1-module.js").Parameter31Schema
>;
export type PathItem31 = z.infer<
	typeof import("./openapi-3.1-module.js").PathItem31Schema
>;
export type Paths31 = z.infer<
	typeof import("./openapi-3.1-module.js").Paths31Schema
>;
export type RequestBody31 = z.infer<
	typeof import("./openapi-3.1-module.js").RequestBody31Schema
>;
export type Response31 = z.infer<
	typeof import("./openapi-3.1-module.js").Response31Schema
>;
export type Responses31 = z.infer<
	typeof import("./openapi-3.1-module.js").Responses31Schema
>;
export type SchemaObject31 = z.infer<
	typeof import("./openapi-3.1-module.js").SchemaObject31Schema
>;
export type SecurityRequirement31 = z.infer<
	typeof import("./openapi-3.1-module.js").SecurityRequirement31Schema
>;
export type SecurityScheme31 = z.infer<
	typeof import("./openapi-3.1-module.js").SecurityScheme31Schema
>;
export type Server31 = z.infer<
	typeof import("./openapi-3.1-module.js").Server31Schema
>;
export type ServerVariable31 = z.infer<
	typeof import("./openapi-3.1-module.js").ServerVariable31Schema
>;
export type Tag31 = z.infer<
	typeof import("./openapi-3.1-module.js").Tag31Schema
>;
export type XML31 = z.infer<
	typeof import("./openapi-3.1-module.js").XML31Schema
>;

// =============================================================================
// OpenAPI 3.2 Schemas
// =============================================================================
export {
	Callback32Schema,
	Components32Schema,
	Example32Schema,
	Header32Schema,
	Link32Schema,
	OpenAPI32Schema,
	Operation32Schema,
	Parameter32Schema,
	PathItem32Schema,
	RequestBody32Schema,
	Response32Schema,
	SchemaObject32Schema,
	SecurityScheme32Schema,
} from "./openapi-3.2-module.js";

export type OpenAPI32 = z.infer<
	typeof import("./openapi-3.2-module.js").OpenAPI32Schema
>;

// =============================================================================
// Base Schemas (shared primitives - Contact, License, etc.)
// =============================================================================
export {
	type Contact,
	ContactSchema,
	type ExternalDocumentation,
	ExternalDocumentationSchema,
	type FileRef,
	FileRefSchema,
	type InternalRef,
	InternalRefSchema,
	type License,
	LicenseSchema,
	type OAuthFlow,
	OAuthFlowSchema,
	type Reference,
	type ReferenceObject,
	ReferenceObjectSchema,
	ReferenceSchema,
	type SecurityRequirement,
	SecurityRequirementSchema,
	type UrlRef,
	UrlRefSchema,
	type XML,
	XMLSchema,
} from "./openapi-base.js";
