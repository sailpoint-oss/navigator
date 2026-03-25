/**
 * OpenAPI 3.0 Schema Module
 *
 * This file is intentionally thin: it composes and re-exports the atomic
 * schema modules under `engine/schemas/3.0/`.
 *
 * @module engine/schemas/openapi-3.0-module
 */

export { CallbackSchema as Callback30Schema } from "./3.0/callback.js";
export { ComponentsObjectSchema as Components30Schema } from "./3.0/components.js";
export { ContactObjectSchema as Contact30Schema } from "./3.0/contact.js";
// Typed Schema Object helpers (used by the central schema index exports)
export { ArraySchemaObject as ArraySchema30 } from "./3.0/data-types/array.js";
export { BooleanSchemaObject as BooleanSchema30 } from "./3.0/data-types/boolean.js";
export { IntegerSchemaObject as IntegerSchema30 } from "./3.0/data-types/integer.js";
export { NumberSchemaObject as NumberSchema30 } from "./3.0/data-types/number.js";
export { ObjectSchemaObject as ObjectSchema30 } from "./3.0/data-types/object.js";
export { StringObjectSchema as StringSchema30 } from "./3.0/data-types/string.js";
export { DiscriminatorObjectSchema as Discriminator30Schema } from "./3.0/discriminator.js";
export { EncodingObjectSchema as Encoding30Schema } from "./3.0/encoding.js";
export { ExampleSchema as Example30Schema } from "./3.0/example.js";
export { ExternalDocumentationObjectSchema as ExternalDocumentation30Schema } from "./3.0/external-documentation.js";
export { HeaderSchema as Header30Schema } from "./3.0/header.js";
export { InfoObjectSchema as Info30Schema } from "./3.0/info.js";
export { LicenseObjectSchema as License30Schema } from "./3.0/license.js";
export { LinkSchema as Link30Schema } from "./3.0/link.js";
export { MediaTypeObjectSchema as MediaType30Schema } from "./3.0/media-type.js";
export { OAuthFlowObjectSchema as OAuthFlow30Schema } from "./3.0/oauth-flow.js";
export { OAuthFlowsObjectSchema as OAuthFlows30Schema } from "./3.0/oauth-flows.js";
export type { OpenAPIObject as OpenAPI3 } from "./3.0/openapi.js";
// Back-compat alias: some callers treat OpenAPI 3.0 as "OpenAPI3"
export {
	OpenAPIObjectSchema as OpenAPI30Schema,
	OpenAPIObjectSchema as OpenAPI3Schema,
} from "./3.0/openapi.js";
export { OperationObjectSchema as Operation30Schema } from "./3.0/operation.js";
export { ParameterSchema as Parameter30Schema } from "./3.0/parameter.js";
export { PathItemObjectSchema as PathItem30Schema } from "./3.0/path-item.js";
export { PathsObjectSchema as Paths30Schema } from "./3.0/paths.js";
export { RequestBodySchema as RequestBody30Schema } from "./3.0/request-body.js";
export { ResponseSchema as Response30Schema } from "./3.0/response.js";
export { ResponsesObjectSchema as Responses30Schema } from "./3.0/responses.js";
export { SchemaObjectSchema as SchemaObject30Schema } from "./3.0/schema.js";
export { SecurityRequirementObjectSchema as SecurityRequirement30Schema } from "./3.0/security-requirement.js";
export { SecuritySchemeSchema as SecurityScheme30Schema } from "./3.0/security-scheme.js";
export { ServerObjectSchema as Server30Schema } from "./3.0/server.js";
export { ServerVariableObjectSchema as ServerVariable30Schema } from "./3.0/server-variable.js";
export { TagObjectSchema as Tag30Schema } from "./3.0/tag.js";
export { XMLObjectSchema as XML30Schema } from "./3.0/xml.js";
