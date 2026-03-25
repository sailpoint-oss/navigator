import z from "zod";
import { BaseSchemaObjectSchema } from "./base";

export const NullSchemaObject = BaseSchemaObjectSchema.extend({
	type: z.literal("null"),
}).meta({
	title: "Null Object",
	description: "A Null Object defines the shape of a null value.",
	examples: [{ type: "null" }],
});

export type NullSchema = z.infer<typeof NullSchemaObject>;
