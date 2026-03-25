import z from "zod";
import { CallbackSchema } from "./callback";
import { ExternalDocumentationObjectSchema } from "./external-documentation";
import { ParameterSchema } from "./parameter";
import { RequestBodySchema } from "./request-body";
import { ResponsesObjectSchema } from "./responses";
import { SecurityRequirementObjectSchema } from "./security-requirement";
import { ServerObjectSchema } from "./server";

export const OperationObjectSchema = z
	.looseObject({
		tags: z
			.array(z.string())
			.meta({
				title: "tags",
				examples: [["pets"], ["users", "authentication"], ["orders", "store"]],
			})
			.describe("A list of tags for API documentation grouping.")
			.optional(),
		summary: z
			.string()
			.meta({ title: "summary", examples: ["List all pets"] })
			.describe("A short summary of the operation. Keep under ~120 characters.")
			.optional(),
		description: z
			.string()
			.meta({ title: "description" })
			.describe("A verbose description. CommonMark syntax MAY be used.")
			.optional(),
		operationId: z
			.string()
			.meta({ title: "operationId", examples: ["listPets", "createUser"] })
			.describe(
				"Unique identifier for the operation. Used for code generation.",
			)
			.optional(),
		parameters: z
			.array(ParameterSchema)
			.meta({ title: "parameters" })
			.describe(
				"Parameters for this operation, combined with path-level parameters.",
			)
			.optional(),
		requestBody: RequestBodySchema.optional().meta({ title: "requestBody" }),
		responses: ResponsesObjectSchema.meta({ title: "responses" }).describe(
			"REQUIRED. The responses for this operation.",
		),
		get callbacks() {
			return z
				.record(z.string(), CallbackSchema)
				.meta({ title: "callbacks" })
				.describe("Webhooks/callbacks triggered by this operation.")
				.optional();
		},
		deprecated: z
			.boolean()
			.optional()
			.meta({
				title: "deprecated",
				description: "Marks this operation as deprecated. Defaults to false.",
				examples: [true, false],
			}),
		security: z.array(SecurityRequirementObjectSchema).optional().meta({
			title: "security",
			description:
				"Security requirements for this operation. Overrides root-level security.",
		}),
		servers: z.array(ServerObjectSchema).optional().meta({ title: "servers" }),
		externalDocs: ExternalDocumentationObjectSchema.optional(),
	})
	.meta({
		title: "Operation",
		description: "Describes a single API operation on a path.",
		examples: [
			{
				summary: "List pets",
				operationId: "listPets",
				responses: { "200": { description: "OK" } },
			},
		],
	});

export type OperationObject = z.infer<typeof OperationObjectSchema>;
