import z from "zod";
import { BaseSchemaObjectSchema } from "./base";

export const BooleanSchemaObject = BaseSchemaObjectSchema.extend({
	type: z.literal("boolean"),
}).meta({
	title: "Boolean Object",
	description: "A Boolean Object defines the shape of a boolean value.",
	examples: [{ type: "boolean" }],
});

export type BooleanSchema = z.infer<typeof BooleanSchemaObject>;


