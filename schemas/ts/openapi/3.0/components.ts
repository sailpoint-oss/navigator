import z from "zod";
import { CallbackSchema } from "./callback";
import { ExampleSchema } from "./example";
import { HeaderSchema } from "./header";
import { LinkSchema } from "./link";
import { ParameterSchema } from "./parameter";
import { RequestBodySchema } from "./request-body";
import { ResponseSchema } from "./response";
import { SchemaObjectSchema } from "./schema";
import { SecuritySchemeSchema } from "./security-scheme";

export const ComponentsObjectSchema = z
	.looseObject({
		schemas: z.record(z.string(), SchemaObjectSchema).optional(),
		responses: z.record(z.string(), ResponseSchema).optional(),
		parameters: z.record(z.string(), ParameterSchema).optional(),
		examples: z.record(z.string(), ExampleSchema).optional(),
		requestBodies: z.record(z.string(), RequestBodySchema).optional(),
		headers: z.record(z.string(), HeaderSchema).optional(),
		securitySchemes: z.record(z.string(), SecuritySchemeSchema).optional(),
		links: z.record(z.string(), LinkSchema).optional(),
		callbacks: z.record(z.string(), CallbackSchema).optional(),
	})
	.meta({
		title: "Components",
		description:
			"Holds a set of reusable objects for different aspects of the OAS.",
		examples: [{ schemas: { Pet: { type: "object" } } }],
	});

export type ComponentsObject = z.infer<typeof ComponentsObjectSchema>;
