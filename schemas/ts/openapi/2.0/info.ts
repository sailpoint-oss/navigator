import z from "zod";
import { ContactObjectSchema } from "./contact";
import { LicenseObjectSchema } from "./license";

export const InfoObjectSchema = z
	.looseObject({
		title: z
			.string()
			.describe("REQUIRED. The title of the application.")
			.meta({ title: "title", examples: ["Swagger Sample App"] }),
		description: z
			.string()
			.optional()
			.describe(
				"A short description of the application. GFM syntax can be used for rich text representation.",
			)
			.meta({ title: "description" }),
		termsOfService: z
			.string()
			.optional()
			.describe("The Terms of Service for the API.")
			.meta({ title: "termsOfService", examples: ["http://swagger.io/terms/"] }),
		contact: ContactObjectSchema.optional().meta({ title: "contact" }),
		license: LicenseObjectSchema.optional().meta({ title: "license" }),
		version: z
			.string()
			.describe(
				"REQUIRED. Provides the version of the application API (not the specification version).",
			)
			.meta({ title: "version", examples: ["1.0.1"] }),
	})
	.meta({ title: "Info" });

export type InfoObject = z.infer<typeof InfoObjectSchema>;


