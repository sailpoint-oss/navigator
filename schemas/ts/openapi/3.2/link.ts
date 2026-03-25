import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ServerObjectSchema } from "./server";

export const LinkObjectSchema = z
	.looseObject({
		operationId: z
			.string()
			.meta({ title: "operationId" })
			.describe(
				"The name of an existing, resolvable OAS operation, as defined with a unique operationId.",
			)
			.optional(),
		operationRef: z
			.string()
			.meta({ title: "operationRef" })
			.describe("A relative or absolute URI reference to an OAS operation.")
			.optional(),
		parameters: z
			.record(z.string(), z.unknown())
			.meta({ title: "parameters" })
			.describe(
				"A map representing parameters to pass to an operation as specified with operationId or identified via operationRef.",
			)
			.optional(),
		requestBody: z
			.unknown()
			.meta({ title: "requestBody" })
			.describe(
				"A literal value or {expression} to use as a request body when calling the target operation.",
			)
			.optional(),
		description: z
			.string()
			.meta({ title: "description" })
			.describe("A description of the link.")
			.optional(),
		server: ServerObjectSchema.optional(),
	})
	.meta({
		title: "Link",
		description:
			"The Link object represents a possible design-time link for a response.",
		examples: [
			{
				operationId: "getUserById",
				parameters: { userId: "$response.body#/id" },
			},
		],
	});

export const LinkSchema = z
	.union([ReferenceObjectSchema, LinkObjectSchema])
	.meta({
		title: "Link",
		description:
			"The Link object represents a possible design-time link for a response.",
	});

export type LinkObject = z.infer<typeof LinkObjectSchema>;
export type Link = z.infer<typeof LinkSchema>;
