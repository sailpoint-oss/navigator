import z from "zod";
import { StringFormatSchema } from "../../data-types";
import { BaseSchemaObjectSchema } from "./base";

export const StringObjectSchema = BaseSchemaObjectSchema.extend({
	type: z.literal("string"),
	format: StringFormatSchema.optional(),
	pattern: z
		.string()
		.meta({
			title: "Pattern",
			description: "A regular expression pattern the string must match.",
			examples: ["^[a-zA-Z0-9]+$", "^\\d{3}-\\d{2}-\\d{4}$"],
		})
		.optional(),
	minLength: z
		.number()
		.int()
		.min(0)
		.optional()
		.meta({
			title: "Minimum Length",
			description: "Minimum string length.",
			examples: [0, 1, 10],
		}),
	maxLength: z
		.number()
		.int()
		.min(0)
		.optional()
		.meta({
			title: "Maximum Length",
			description: "Maximum string length.",
			examples: [255, 1024],
		}),
	contentEncoding: z
		.string()
		.optional()
		.meta({
			title: "contentEncoding",
			description: "Encoding for binary string content.",
			examples: ["base64", "base64url"],
		}),
	contentMediaType: z
		.string()
		.optional()
		.meta({
			title: "contentMediaType",
			description: "Media type of the encoded string content.",
			examples: ["application/octet-stream", "image/png"],
		}),
}).meta({
	title: "String Object",
	description: "A String Object defines the shape of a string value.",
	examples: [
		{ type: "string", format: "email" },
		{ type: "string", pattern: "^[a-zA-Z0-9]+$" },
		{ type: "string", minLength: 10, maxLength: 100 },
	],
});

export type StringObject = z.infer<typeof StringObjectSchema>;
