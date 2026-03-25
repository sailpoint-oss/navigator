import z from "zod";

export const DiscriminatorObjectSchema = z
	.looseObject({
		propertyName: z
			.string()
			.meta({ title: "propertyName" })
			.describe(
				"The name of the property in the payload that will hold the discriminator value.",
			),
		mapping: z
			.record(z.string(), z.string())
			.meta({ title: "mapping" })
			.describe(
				"An object to hold mappings between payload values and schema names or references.",
			)
			.optional(),
	})
	.meta({
		title: "Discriminator",
		description:
			"When request bodies or response payloads may be one of a number of different schemas, a discriminator object can be used to aid in serialization, deserialization, and validation.",
		examples: [
			{ propertyName: "petType", mapping: { dog: "#/components/schemas/Dog" } },
		],
	});

export type DiscriminatorObject = z.infer<typeof DiscriminatorObjectSchema>;


