import z from "zod";

export const LicenseIdentifierSchema = z.string().meta({
	title: "License Identifier",
	description: "An SPDX license expression for the API.",
	examples: ["Apache-2.0", "MIT", "BSD-3-Clause", "GPL-3.0-only"],
});

export const LicenseUrlSchema = z.url().meta({
	title: "License URL",
	description: "A URL to the license used for the API.",
	examples: [
		"https://www.apache.org/licenses/LICENSE-2.0.html",
		"https://opensource.org/licenses/MIT",
	],
});

export const LicenseNameSchema = z.string().meta({
	title: "License Name",
	description: "The license name used for the API.",
	examples: ["Apache 2.0", "MIT", "BSD-3-Clause", "GPL-3.0"],
});

export const LicenseObjectSchema = z
	.looseObject({
		name: LicenseNameSchema,
		url: LicenseUrlSchema.optional(),
		identifier: LicenseIdentifierSchema.optional(),
	})
	.superRefine((value, ctx) => {
		// OpenAPI 3.1: url and identifier are mutually exclusive, but both are optional.
		if (value.url && value.identifier) {
			ctx.addIssue({
				code: "custom",
				message: "License Object: 'url' and 'identifier' are mutually exclusive",
				path: ["url"],
			});
		}
	})
	.meta({
		title: "License",
		description: "License information for the exposed API.",
		examples: [
			{ name: "Apache 2.0", identifier: "Apache-2.0" },
			{ name: "MIT", url: "https://opensource.org/licenses/MIT" },
		],
	});

export type LicenseObject = z.infer<typeof LicenseObjectSchema>;
