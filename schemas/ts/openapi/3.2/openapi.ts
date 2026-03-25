import z from "zod";
import { ComponentsObjectSchema } from "./components";
import { ExternalDocumentationObjectSchema } from "./external-documentation";
import { InfoObjectSchema } from "./info";
import { PathItemObjectSchema } from "./path-item";
import { PathsObjectSchema } from "./paths";
import { SecurityRequirementObjectSchema } from "./security-requirement";
import { ServerObjectSchema } from "./server";
import { TagObjectSchema } from "./tag";

export const OpenAPIObjectSchema: z.ZodType = z
	.looseObject({
		openapi: z
			.string()
			.regex(/^3\.2\.\d+$/)
			.meta({
				title: "openapi",
				description:
					"OpenAPI version. Must be '3.2.x' for OpenAPI 3.2 documents.",
				examples: ["3.2.0"],
			}),
		info: InfoObjectSchema.meta({
			title: "info",
			description: "API metadata.",
		}),
		jsonSchemaDialect: z
			.url()
			.meta({ title: "jsonSchemaDialect" })
			.describe("The default JSON Schema dialect for Schema Objects.")
			.optional(),
		servers: z.array(ServerObjectSchema).optional().meta({ title: "servers" }),
		paths: PathsObjectSchema.optional().meta({
			title: "paths",
			description:
				"Available paths and operations. Optional if webhooks is provided.",
		}),
		webhooks: z.record(z.string(), PathItemObjectSchema).optional().meta({
			title: "webhooks",
			description: "Incoming webhooks that the API can receive.",
		}),
		components: ComponentsObjectSchema.optional().meta({
			title: "components",
			description: "Reusable schemas, parameters, responses, etc.",
		}),
		security: z
			.array(SecurityRequirementObjectSchema)
			.optional()
			.meta({ title: "security" }),
		tags: z.array(TagObjectSchema).optional().meta({ title: "tags" }),
		externalDocs: ExternalDocumentationObjectSchema.optional(),
	})
	.meta({
		title: "OpenAPI",
		description: "Root object of an OpenAPI 3.2 document.",
	});

export type OpenAPIObject = z.infer<typeof OpenAPIObjectSchema>;
