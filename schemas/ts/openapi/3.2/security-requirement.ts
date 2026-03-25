import z from "zod";

export const SecurityRequirementObjectSchema = z
	.record(z.string(), z.array(z.string()))
	.meta({
		title: "SecurityRequirement",
		description: "Lists the required security schemes for this operation.",
		examples: [{ api_key: [] }, { oauth2: ["read:pets", "write:pets"] }],
	});

export type SecurityRequirementObject = z.infer<
	typeof SecurityRequirementObjectSchema
>;


