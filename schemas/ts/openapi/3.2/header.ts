import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";
import { ExampleSchema } from "./example";
import { SchemaObjectSchema } from "./schema";

export const HeaderObjectSchema = z
	.looseObject({
		description: z
			.string()
			.meta({ title: "description" })
			.describe("A brief description of the parameter.")
			.optional(),
		required: z
			.boolean()
			.optional()
			.describe("Determines whether this parameter is mandatory.")
			.meta({ title: "required" }),
		deprecated: z
			.boolean()
			.optional()
			.describe("Specifies that a parameter is deprecated.")
			.meta({ title: "deprecated" }),
		allowEmptyValue: z
			.boolean()
			.optional()
			.describe("Sets the ability to pass empty-valued parameters.")
			.meta({ title: "allowEmptyValue" }),
		style: z
			.literal("simple")
			.meta({ title: "style" })
			.describe("Describes how the parameter value will be serialized.")
			.optional(),
		explode: z.boolean().optional().meta({ title: "explode" }),
		allowReserved: z
			.boolean()
			.optional()
			.meta({ title: "allowReserved" }),
		schema: SchemaObjectSchema.optional().meta({ title: "schema" }),
		example: z.unknown().optional().meta({ title: "example" }),
		examples: z
			.record(z.string(), ExampleSchema)
			.meta({ title: "examples" })
			.optional(),
		content: z
			.record(z.string(), z.unknown())
			.meta({ title: "content" })
			.optional(),
	})
	.meta({
		title: "Header",
		description:
			"The Header Object follows the structure of the Parameter Object.",
		examples: [
			{ description: "Rate limit remaining", schema: { type: "integer" } },
		],
	});

export const HeaderSchema = z
	.union([ReferenceObjectSchema, HeaderObjectSchema])
	.meta({
		title: "Header",
		description:
			"The Header Object follows the structure of the Parameter Object.",
	});

export type HeaderObject = z.infer<typeof HeaderObjectSchema>;
export type Header = z.infer<typeof HeaderSchema>;
