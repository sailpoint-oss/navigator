import z from "zod";

export const LicenseObjectSchema = z
	.looseObject({
		name: z
			.string()
			.describe("REQUIRED. The license name used for the API.")
			.meta({ title: "name", examples: ["Apache 2.0"] }),
		url: z
			.string()
			.optional()
			.describe("A URL to the license used for the API. MUST be a URL.")
			.meta({
				title: "url",
				examples: ["http://www.apache.org/licenses/LICENSE-2.0.html"],
			}),
	})
	.meta({ title: "License" });

export type LicenseObject = z.infer<typeof LicenseObjectSchema>;


