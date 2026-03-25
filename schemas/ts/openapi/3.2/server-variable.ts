import z from "zod";

export const ServerVariableObjectSchema = z
	.looseObject({
		enum: z
			.array(z.string())
			.meta({
				title: "enum",
				examples: [["https", "http"], ["prod", "staging", "dev"], ["8080", "443"]],
			})
			.describe("An enumeration of valid string values for this variable.")
			.optional(),
		default: z
			.string()
			.meta({
				title: "default",
				examples: ["https", "api.example.com", "443", "v1"],
			})
			.describe("The default value for substitution if not supplied."),
		description: z
			.string()
			.meta({
				title: "description",
				examples: ["The API protocol", "Server environment", "Port number"],
			})
			.describe("An optional description for the server variable.")
			.optional(),
	})
	.meta({
		title: "ServerVariable",
		description:
			"A server variable for URL template substitution. Use {variableName} in the server URL.",
		examples: [{ default: "https", enum: ["https", "http"] }],
	});

export type ServerVariableObject = z.infer<typeof ServerVariableObjectSchema>;


