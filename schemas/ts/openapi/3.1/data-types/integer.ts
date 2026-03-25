import z from "zod";
import { IntegerFormatSchema } from "../../data-types";
import { BaseSchemaObjectSchema } from "./base";

export const IntegerSchemaObject = BaseSchemaObjectSchema.extend({
	type: z.literal("integer"),
	format: IntegerFormatSchema.optional(),
	multipleOf: z
		.number()
		.optional()
		.meta({
			title: "multipleOf",
			description: "Require the integer to be a multiple of this value.",
			examples: [1, 2, 10],
		}),
	minimum: z
		.number()
		.optional()
		.meta({
			title: "minimum",
			description: "Inclusive lower bound for the integer.",
			examples: [0, 1, -10],
		}),
	maximum: z
		.number()
		.optional()
		.meta({
			title: "maximum",
			description: "Inclusive upper bound for the integer.",
			examples: [100, 1000],
		}),
	exclusiveMinimum: z
		.number()
		.optional()
		.meta({
			title: "exclusiveMinimum",
			description: "Exclusive lower bound (instance must be > this value).",
			examples: [0, 1],
		}),
	exclusiveMaximum: z
		.number()
		.optional()
		.meta({
			title: "exclusiveMaximum",
			description: "Exclusive upper bound (instance must be < this value).",
			examples: [100, 1000],
		}),
}).meta({
	title: "Integer Object",
	description: "An Integer Object defines the shape of an integer value.",
	examples: [
		{ type: "integer", format: "int32" },
		{ type: "integer", minimum: 0 },
	],
});
