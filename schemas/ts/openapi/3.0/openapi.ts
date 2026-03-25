import z from "zod";
import { ComponentsObjectSchema } from "./components";
import { ExternalDocumentationObjectSchema } from "./external-documentation";
import { InfoObjectSchema } from "./info";
import { PathsObjectSchema } from "./paths";
import { SecurityRequirementObjectSchema } from "./security-requirement";
import { ServerObjectSchema } from "./server";
import { TagObjectSchema } from "./tag";

export const OpenAPIObjectSchema: z.ZodType = z
	.looseObject({
		openapi: z
			.string()
			.regex(/^3\.0\.\d+$/)
			.meta({
				title: "openapi",
				description:
					"OpenAPI version. Must be '3.0.x' for OpenAPI 3.0 documents.",
				examples: ["3.0.4"],
			}),
		info: InfoObjectSchema.meta({
			title: "info",
			description: "API metadata.",
		}),
		servers: z.array(ServerObjectSchema).optional().meta({ title: "servers" }),
		paths: PathsObjectSchema.meta({
			title: "paths",
			description: "Available paths and operations for the API.",
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
		description: "Root object of an OpenAPI 3.0 document.",
	});

export type OpenAPIObject = z.infer<typeof OpenAPIObjectSchema>;
