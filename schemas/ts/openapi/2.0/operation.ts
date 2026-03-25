import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ExternalDocumentationObjectSchema } from "./external-documentation";
import { ParameterObjectSchema } from "./parameter";
import { MimeTypeStringSchema } from "./primitives";
import { ResponsesObjectSchema } from "./responses";
import { SecurityRequirementObjectSchema } from "./security-requirement";

export const OperationObjectSchema = z
	.looseObject({
		tags: z.array(z.string()).optional().meta({ title: "tags" }),
		summary: z.string().optional().meta({ title: "summary" }),
		description: z.string().optional().meta({ title: "description" }),
		externalDocs: ExternalDocumentationObjectSchema.optional().meta({
			title: "externalDocs",
		}),
		operationId: z.string().optional().meta({ title: "operationId" }),
		consumes: z
			.array(MimeTypeStringSchema)
			.optional()
			.meta({ title: "consumes" }),
		produces: z
			.array(MimeTypeStringSchema)
			.optional()
			.meta({ title: "produces" }),
		parameters: z
			.array(z.union([ParameterObjectSchema, ReferenceObjectSchema]))
			.optional()
			.meta({ title: "parameters" }),
		responses: ResponsesObjectSchema.describe(
			"REQUIRED. The responses for this operation.",
		).meta({
			title: "responses",
		}),
		schemes: z
			.array(z.enum(["http", "https", "ws", "wss"]))
			.optional()
			.meta({ title: "schemes" }),
		deprecated: z.boolean().optional().meta({ title: "deprecated" }),
		security: z
			.array(SecurityRequirementObjectSchema)
			.optional()
			.meta({ title: "security" }),
	})
	.meta({ title: "Operation" });

export type OperationObject = z.infer<typeof OperationObjectSchema>;
