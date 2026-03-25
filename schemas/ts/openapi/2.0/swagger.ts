import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ExternalDocumentationObjectSchema } from "./external-documentation";
import { InfoObjectSchema } from "./info";
import { ParameterObjectSchema } from "./parameter";
import { PathsObjectSchema } from "./paths";
import { MimeTypeStringSchema } from "./primitives";
import { ResponseSchema } from "./response";
import { SchemaObjectSchema } from "./schema";
import { SecurityRequirementObjectSchema } from "./security-requirement";
import { TagObjectSchema } from "./tag";

export const DefinitionsObjectSchema = z
	.record(z.string(), SchemaObjectSchema)
	.meta({ title: "Definitions" });

export const ParametersDefinitionsObjectSchema = z
	.record(z.string(), z.union([ParameterObjectSchema, ReferenceObjectSchema]))
	.meta({ title: "ParametersDefinitions" });

export const ResponsesDefinitionsObjectSchema = z
	.record(z.string(), z.union([ResponseSchema, ReferenceObjectSchema]))
	.meta({ title: "ResponsesDefinitions" });

export const SwaggerObjectSchema = z
	.looseObject({
		swagger: z
			.literal("2.0")
			.describe(
				"REQUIRED. Specifies the Swagger Specification version being used.",
			)
			.meta({ title: "swagger", examples: ["2.0"] }),
		info: InfoObjectSchema.describe(
			"REQUIRED. Provides metadata about the API.",
		).meta({
			title: "info",
		}),
		host: z
			.string()
			.optional()
			.describe(
				"The host (name or ip) serving the API. MAY include a port. Does not include scheme or basePath.",
			)
			.meta({ title: "host" }),
		basePath: z
			.string()
			.optional()
			.describe("The base path on which the API is served, relative to host.")
			.meta({ title: "basePath" }),
		schemes: z
			.array(z.enum(["http", "https", "ws", "wss"]))
			.optional()
			.meta({ title: "schemes" }),
		consumes: z
			.array(MimeTypeStringSchema)
			.optional()
			.meta({ title: "consumes" }),
		produces: z
			.array(MimeTypeStringSchema)
			.optional()
			.meta({ title: "produces" }),
		paths: PathsObjectSchema.describe(
			"REQUIRED. The available paths and operations.",
		).meta({
			title: "paths",
		}),
		definitions: DefinitionsObjectSchema.optional().meta({
			title: "definitions",
		}),
		parameters: ParametersDefinitionsObjectSchema.optional().meta({
			title: "parameters",
		}),
		responses: ResponsesDefinitionsObjectSchema.optional().meta({
			title: "responses",
		}),
		security: z
			.array(SecurityRequirementObjectSchema)
			.optional()
			.meta({ title: "security" }),
		tags: z.array(TagObjectSchema).optional().meta({ title: "tags" }),
		externalDocs: ExternalDocumentationObjectSchema.optional().meta({
			title: "externalDocs",
		}),
	})
	.meta({ title: "Swagger" });

export type SwaggerObject = z.infer<typeof SwaggerObjectSchema>;
