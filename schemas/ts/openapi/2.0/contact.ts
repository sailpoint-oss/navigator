import z from "zod";

export const ContactObjectSchema = z
	.looseObject({
		name: z
			.string()
			.optional()
			.describe("The identifying name of the contact person/organization.")
			.meta({ title: "name", examples: ["API Support"] }),
		url: z
			.string()
			.optional()
			.describe("The URL pointing to the contact information. MUST be a URL.")
			.meta({ title: "url", examples: ["http://www.swagger.io/support"] }),
		email: z
			.string()
			.optional()
			.describe("The email address of the contact person/organization.")
			.meta({ title: "email", examples: ["support@swagger.io"] }),
	})
	.meta({ title: "Contact" });

export type ContactObject = z.infer<typeof ContactObjectSchema>;


