import z from "zod";
import { ReferenceObjectSchema } from "../openapi-base";

export const ExampleObjectSchema = z
	.looseObject({
		summary: z
			.string()
			.meta({ title: "summary" })
			.describe("Short description for the example.")
			.optional(),
		description: z
			.string()
			.meta({ title: "description" })
			.describe("Long description for the example.")
			.optional(),
		value: z
			.unknown()
			.meta({ title: "value" })
			.describe("Embedded literal example.")
			.optional(),
		externalValue: z
			.url()
			.meta({ title: "externalValue" })
			.describe("A URL that points to the literal example.")
			.optional(),
		dataValue: z
			.unknown()
			.meta({ title: "dataValue" })
			.describe(
				"The data value of the example before serialization. Mutually exclusive with value and externalValue.",
			)
			.optional(),
		serializedValue: z
			.string()
			.meta({ title: "serializedValue" })
			.describe(
				"The serialized representation of the example. Mutually exclusive with value and externalValue.",
			)
			.optional(),
	})
	.meta({
		title: "Example",
		description: "Example Object",
		examples: [{ summary: "A sample", value: { id: 1, name: "Example" } }],
	});

export const ExampleSchema = z
	.union([ReferenceObjectSchema, ExampleObjectSchema])
	.meta({
		title: "Example",
		description: "Example Object",
		examples: [{ summary: "A sample", value: { id: 1, name: "Example" } }],
	});

export type ExampleObject = z.infer<typeof ExampleObjectSchema>;
export type Example = z.infer<typeof ExampleSchema>;
