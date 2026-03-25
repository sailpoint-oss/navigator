import z from "zod";
import { NumberFormatSchema } from "../../data-types";
import { BaseSchemaObjectSchema } from "./base";

export const NumberSchemaObject = BaseSchemaObjectSchema.extend({
	type: z.literal("number"),
	format: NumberFormatSchema.optional(),
	multipleOf: z
		.number()
		.optional()
		.meta({
			title: "multipleOf",
			description: "Require the number to be a multiple of this value.",
			examples: [0.5, 1, 10],
		}),
	minimum: z
		.number()
		.optional()
		.meta({
			title: "minimum",
			description: "Inclusive lower bound for the number.",
			examples: [0, 1, -10],
		}),
	maximum: z
		.number()
		.optional()
		.meta({
			title: "maximum",
			description: "Inclusive upper bound for the number.",
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
	title: "Number Object",
	description: "A Number Object defines the shape of a number value.",
	examples: [
		{ type: "number", format: "float" },
		{ type: "number", multipleOf: 0.5 },
		{ type: "number", minimum: 0, maximum: 100 },
		{ type: "number", exclusiveMinimum: 0, exclusiveMaximum: 100 },
	],
});
