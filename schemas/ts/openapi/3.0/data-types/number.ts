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
		.boolean()
		.optional()
		.meta({
			title: "exclusiveMinimum",
			description:
				"When true, the minimum value is exclusive (value must be > minimum).",
			examples: [true, false],
		}),
	exclusiveMaximum: z
		.boolean()
		.optional()
		.meta({
			title: "exclusiveMaximum",
			description:
				"When true, the maximum value is exclusive (value must be < maximum).",
			examples: [true, false],
		}),
}).meta({
	title: "Number Object",
	description: "A Number Object defines the shape of a number value.",
	examples: [
		{ type: "number", format: "float" },
		{ type: "number", multipleOf: 0.5 },
		{ type: "number", minimum: 0, maximum: 100 },
		{ type: "number", minimum: 0, exclusiveMinimum: true },
	],
});
