import z from "zod";

export const XMLObjectSchema = z
	.looseObject({
		name: z
			.string()
			.meta({ title: "name" })
			.describe(
				"Replaces the name of the element/attribute used for the described schema property.",
			)
			.optional(),
		namespace: z
			.url()
			.meta({ title: "namespace" })
			.describe("The URI of the namespace definition.")
			.optional(),
		prefix: z
			.string()
			.meta({ title: "prefix" })
			.describe("The prefix to be used for the name.")
			.optional(),
		attribute: z
			.boolean()
			.optional()
			.describe(
				"Declares whether the property definition translates to an attribute instead of an element.",
			)
			.meta({ title: "attribute" }),
		wrapped: z
			.boolean()
			.optional()
			.describe(
				"May be used only for an array definition. Signifies whether the array is wrapped or not.",
			)
			.meta({ title: "wrapped" }),
	})
	.meta({
		title: "XML",
		description:
			"A metadata object that allows for more fine-tuned XML model definitions.",
		examples: [{ name: "animal", wrapped: true }],
	});

export type XMLObject = z.infer<typeof XMLObjectSchema>;


