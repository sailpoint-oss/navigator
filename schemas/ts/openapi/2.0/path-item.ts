import z from "zod";
import { ReferenceObjectSchema, ReferenceSchema } from "../openapi-base";
import { OperationObjectSchema } from "./operation";
import { ParameterObjectSchema } from "./parameter";

export const PathItemObjectSchema = z
	.looseObject({
		$ref: ReferenceSchema,
		get: OperationObjectSchema.optional().meta({ title: "get" }),
		put: OperationObjectSchema.optional().meta({ title: "put" }),
		post: OperationObjectSchema.optional().meta({ title: "post" }),
		delete: OperationObjectSchema.optional().meta({ title: "delete" }),
		options: OperationObjectSchema.optional().meta({ title: "options" }),
		head: OperationObjectSchema.optional().meta({ title: "head" }),
		patch: OperationObjectSchema.optional().meta({ title: "patch" }),
		parameters: z
			.array(z.union([ParameterObjectSchema, ReferenceObjectSchema]))
			.optional()
			.meta({ title: "parameters" }),
	})
	.meta({ title: "PathItem" });

export type PathItemObject = z.infer<typeof PathItemObjectSchema>;
