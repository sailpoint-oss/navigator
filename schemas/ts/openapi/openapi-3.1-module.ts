/**
 * OpenAPI 3.1 Schema Module
 *
 * This file is intentionally thin: it composes and re-exports the atomic
 * schema modules under `engine/schemas/3.1/`.
 *
 * @module engine/schemas/openapi-3.1-module
 */

export { CallbackSchema as Callback31Schema } from "./3.1/callback.js";
export { ComponentsObjectSchema as Components31Schema } from "./3.1/components.js";
export { ContactObjectSchema as Contact31Schema } from "./3.1/contact.js";
// Typed Schema Object helpers (used by the central schema index exports)
export { ArraySchemaObject as ArraySchema31 } from "./3.1/data-types/array.js";
export { BooleanSchemaObject as BooleanSchema31 } from "./3.1/data-types/boolean.js";
export { IntegerSchemaObject as IntegerSchema31 } from "./3.1/data-types/integer.js";
export { NullSchemaObject as NullSchema31 } from "./3.1/data-types/null.js";
export { NumberSchemaObject as NumberSchema31 } from "./3.1/data-types/number.js";
export { ObjectSchemaObject as ObjectSchema31 } from "./3.1/data-types/object.js";
export { StringObjectSchema as StringSchema31 } from "./3.1/data-types/string.js";
export { DiscriminatorObjectSchema as Discriminator31Schema } from "./3.1/discriminator.js";
export { EncodingObjectSchema as Encoding31Schema } from "./3.1/encoding.js";
export { ExampleSchema as Example31Schema } from "./3.1/example.js";
export { ExternalDocumentationObjectSchema as ExternalDocumentation31Schema } from "./3.1/external-documentation.js";
export { HeaderSchema as Header31Schema } from "./3.1/header.js";
export { InfoObjectSchema as Info31Schema } from "./3.1/info.js";
export { LicenseObjectSchema as License31Schema } from "./3.1/license.js";
export { LinkSchema as Link31Schema } from "./3.1/link.js";
export { MediaTypeObjectSchema as MediaType31Schema } from "./3.1/media-type.js";
export { OAuthFlowObjectSchema as OAuthFlow31Schema } from "./3.1/oauth-flow.js";
export { OAuthFlowsObjectSchema as OAuthFlows31Schema } from "./3.1/oauth-flows.js";
export { OpenAPIObjectSchema as OpenAPI31Schema } from "./3.1/openapi.js";
export { OperationObjectSchema as Operation31Schema } from "./3.1/operation.js";
export { ParameterSchema as Parameter31Schema } from "./3.1/parameter.js";
export { PathItemObjectSchema as PathItem31Schema } from "./3.1/path-item.js";
export { PathsObjectSchema as Paths31Schema } from "./3.1/paths.js";
export { RequestBodySchema as RequestBody31Schema } from "./3.1/request-body.js";
export { ResponseSchema as Response31Schema } from "./3.1/response.js";
export { ResponsesObjectSchema as Responses31Schema } from "./3.1/responses.js";
export { SchemaObjectSchema as SchemaObject31Schema } from "./3.1/schema.js";
export { SecurityRequirementObjectSchema as SecurityRequirement31Schema } from "./3.1/security-requirement.js";
export { SecuritySchemeSchema as SecurityScheme31Schema } from "./3.1/security-scheme.js";
export { ServerObjectSchema as Server31Schema } from "./3.1/server.js";
export { ServerVariableObjectSchema as ServerVariable31Schema } from "./3.1/server-variable.js";
export { TagObjectSchema as Tag31Schema } from "./3.1/tag.js";
export { XMLObjectSchema as XML31Schema } from "./3.1/xml.js";
