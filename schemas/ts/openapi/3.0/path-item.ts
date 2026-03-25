import z from "zod";
import { ReferenceSchema } from "../openapi-base";
import { OperationObjectSchema } from "./operation";
import { ParameterSchema } from "./parameter";
import { ServerObjectSchema } from "./server";

export const PathItemObjectSchema: z.ZodType = z
	.looseObject({
		$ref: ReferenceSchema.optional(),
		summary: z.string().optional().meta({ title: "summary" }),
		description: z.string().optional().meta({ title: "description" }),
		get: OperationObjectSchema.optional().meta({ title: "get" }),
		put: OperationObjectSchema.optional().meta({ title: "put" }),
		post: OperationObjectSchema.optional().meta({ title: "post" }),
		delete: OperationObjectSchema.optional().meta({ title: "delete" }),
		options: OperationObjectSchema.optional().meta({ title: "options" }),
		head: OperationObjectSchema.optional().meta({ title: "head" }),
		patch: OperationObjectSchema.optional().meta({ title: "patch" }),
		trace: OperationObjectSchema.optional().meta({ title: "trace" }),
		servers: z.array(ServerObjectSchema).optional().meta({ title: "servers" }),
		parameters: z
			.array(ParameterSchema)
			.meta({ title: "parameters" })
			.optional(),
	})
	.meta({
		title: "PathItem",
		description: "Describes the operations available on a single path.",
	});

export type PathItemObject = z.infer<typeof PathItemObjectSchema>;
