import z from "zod";
import { DiscriminatorObjectSchema } from "../discriminator";
import { ExternalDocumentationObjectSchema } from "../external-documentation";
import { XMLObjectSchema } from "../xml";

export const BaseSchemaObjectSchema = z.looseObject({
	title: z
		.string()
		.meta({ title: "title", examples: ["Pet", "User", "Order", "Error"] })
		.describe("A title for the schema, used for documentation.")
		.optional(),
	description: z
		.string()
		.meta({
			title: "description",
			examples: ["A pet in the store", "User account information"],
		})
		.describe("A description of the schema. CommonMark syntax MAY be used.")
		.optional(),
	default: z
		.unknown()
		.meta({ title: "default", examples: ["default value", 0, false, null] })
		.describe("The default value for this schema.")
		.optional(),
	example: z
		.unknown()
		.meta({ title: "example", examples: ["John Doe", 42, { id: 1 }] })
		.describe("An example value for this schema.")
		.optional(),
	enum: z
		.array(z.unknown())
		.meta({ title: "enum" })
		.describe("An array of valid values for this schema.")
		.optional(),
	discriminator: DiscriminatorObjectSchema.optional().meta({ title: "discriminator" }),
	xml: XMLObjectSchema.optional().meta({ title: "xml" }),
	externalDocs: ExternalDocumentationObjectSchema.optional().meta({
		title: "externalDocs",
	}),
	readOnly: z
		.boolean()
		.meta({ title: "readOnly", examples: [true] })
		.describe(
			"When true, the property is only returned in responses, not accepted in requests.",
		)
		.optional(),
	writeOnly: z
		.boolean()
		.meta({ title: "writeOnly", examples: [true] })
		.describe(
			"When true, the property is only accepted in requests, not returned in responses.",
		)
		.optional(),
	deprecated: z
		.boolean()
		.meta({ title: "deprecated", examples: [true] })
		.describe("When true, indicates this schema is deprecated and should be avoided.")
		.optional(),
	nullable: z
		.boolean()
		.meta({ title: "nullable", examples: [true, false] })
		.describe("OpenAPI 3.0: Set to true to allow null values.")
		.optional(),
});

export type BaseSchemaObject = z.infer<typeof BaseSchemaObjectSchema>;


