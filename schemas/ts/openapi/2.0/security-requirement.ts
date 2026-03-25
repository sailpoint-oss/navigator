import z from "zod";

export const SecurityRequirementObjectSchema = z
	.record(z.string(), z.array(z.string()))
	.meta({ title: "SecurityRequirement" });

export type SecurityRequirementObject = z.infer<
	typeof SecurityRequirementObjectSchema
>;


