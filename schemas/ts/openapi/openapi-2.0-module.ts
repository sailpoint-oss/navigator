/**
 * OpenAPI 2.0 (Swagger 2.0) Schema Module
 *
 * Thin wrapper around atomic modules under `engine/schemas/2.0/`.
 *
 * @module engine/schemas/openapi-2.0-module
 */

export { ContactObjectSchema as Contact2Schema } from "./2.0/contact.js";
export { ExternalDocumentationObjectSchema as ExternalDocs2Schema } from "./2.0/external-documentation.js";
export { HeaderObjectSchema as Header2Schema } from "./2.0/header.js";
export { InfoObjectSchema as Info2Schema } from "./2.0/info.js";
export { LicenseObjectSchema as License2Schema } from "./2.0/license.js";
export { OperationObjectSchema as Operation2Schema } from "./2.0/operation.js";
export { ParameterObjectSchema as Parameter2Schema } from "./2.0/parameter.js";
export { PathItemObjectSchema as PathItem2Schema } from "./2.0/path-item.js";
export { PathsObjectSchema as Paths2Schema } from "./2.0/paths.js";
export { ResponseSchema as Response2Schema } from "./2.0/response.js";
export { ResponsesObjectSchema as Responses2Schema } from "./2.0/responses.js";
export {
	ItemsObjectSchema as ItemsObject2Schema,
	SchemaObjectSchema as SchemaObject2Schema,
} from "./2.0/schema.js";
export { SwaggerObjectSchema as OpenAPI2Schema } from "./2.0/swagger.js";
export { TagObjectSchema as Tag2Schema } from "./2.0/tag.js";
