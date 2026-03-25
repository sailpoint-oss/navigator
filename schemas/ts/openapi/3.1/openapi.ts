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
		openapi: z.enum(["3.1.0", "3.1.1", "3.1.2"]).meta({
			title: "openapi",
			description:
				"OpenAPI version. Must be '3.1.x' for OpenAPI 3.1 documents.",
			examples: ["3.1.2"],
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
		description: "Root object of an OpenAPI 3.1 document.",
	});

export type OpenAPIObject = z.infer<typeof OpenAPIObjectSchema>;
