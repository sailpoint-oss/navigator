import z from "zod";
import {
	CookieParameterStyleSchema,
	HeaderParameterStyleSchema,
	PathParameterStyleSchema,
	QueryParameterStyleSchema,
} from "../data-types";
import { ReferenceObjectSchema } from "../openapi-base";
import { ExampleSchema } from "./example";
import { MediaTypeObjectSchema } from "./media-type";
import { SchemaObjectSchema } from "./schema";

const QueryParameterObjectSchema = z
	.looseObject({
		name: z
			.string()
			.meta({
				title: "name",
				examples: ["page", "limit", "filter", "sort", "status"],
			})
			.describe("The name of the query parameter. Case-sensitive."),
		in: z.literal("query").meta({ title: "in", examples: ["query"] }),
		description: z
			.string()
			.meta({
				title: "description",
				examples: [
					"Page number for pagination",
					"Maximum items to return",
				],
			})
			.describe(
				"A description of the parameter. CommonMark syntax MAY be used.",
			)
			.optional(),
		required: z
			.boolean()
			.optional()
			.describe("Whether the parameter is required.")
			.meta({ title: "required", examples: [true, false] }),
		deprecated: z
			.boolean()
			.optional()
			.describe("Marks the parameter as deprecated.")
			.meta({ title: "deprecated", examples: [true] }),
		allowEmptyValue: z
			.boolean()
			.optional()
			.describe("Allows empty values (?param=). Deprecated in favor of schema-level validation.")
			.meta({ title: "allowEmptyValue", examples: [true] }),
		style: QueryParameterStyleSchema.describe(
			"Serialization style. Default: 'form'.",
		).optional(),
		explode: z
			.boolean()
			.meta({ title: "explode", examples: [true, false] })
			.describe(
				"For arrays/objects. When true, each value gets its own parameter.",
			)
			.optional(),
		allowReserved: z
			.boolean()
			.optional()
			.describe(
				"When true, allows reserved characters without percent-encoding.",
			)
			.meta({ title: "allowReserved", examples: [true] }),
		get schema() {
			return SchemaObjectSchema.optional().meta({ title: "schema" });
		},
		example: z
			.unknown()
			.meta({ title: "example", examples: [123, "abc", ["a", "b"]] })
			.describe("Example value. Mutually exclusive with 'examples'.")
			.optional(),
		get examples() {
			return z
				.record(z.string(), ExampleSchema)
				.meta({ title: "examples" })
				.describe(
					"Multiple named examples. Mutually exclusive with 'example'.",
				)
				.optional();
		},
		content: z
			.record(z.string(), MediaTypeObjectSchema)
			.meta({ title: "content" })
			.describe("For complex parameters. Mutually exclusive with 'schema'.")
			.optional(),
	})
	.meta({
		title: "Query Parameter",
		description: "A query string parameter (?name=value).",
		examples: [
			{
				name: "page",
				in: "query",
				schema: { type: "integer" },
			},
		],
	});

const PathParameterObjectSchema = z
	.looseObject({
		name: z
			.string()
			.meta({
				title: "name",
				examples: ["id", "userId", "slug"],
			})
			.describe(
				"The name of the path parameter. Must match a template expression in the path.",
			),
		in: z.literal("path").meta({ title: "in", examples: ["path"] }),
		description: z
			.string()
			.meta({
				title: "description",
				examples: ["Unique identifier", "Resource slug"],
			})
			.describe(
				"A description of the parameter. CommonMark syntax MAY be used.",
			)
			.optional(),
		required: z
			.literal(true)
			.describe("Path parameters must always be required.")
			.meta({ title: "required" }),
		deprecated: z
			.boolean()
			.optional()
			.describe("Marks the parameter as deprecated.")
			.meta({ title: "deprecated", examples: [true] }),
		style: PathParameterStyleSchema.describe(
			"Serialization style. Default: 'simple'.",
		).optional(),
		explode: z
			.boolean()
			.meta({ title: "explode", examples: [true, false] })
			.describe(
				"For arrays/objects. When true, each value gets its own parameter.",
			)
			.optional(),
		get schema() {
			return SchemaObjectSchema.optional().meta({ title: "schema" });
		},
		example: z
			.unknown()
			.meta({ title: "example", examples: [123, "abc"] })
			.describe("Example value. Mutually exclusive with 'examples'.")
			.optional(),
		get examples() {
			return z
				.record(z.string(), ExampleSchema)
				.meta({ title: "examples" })
				.describe(
					"Multiple named examples. Mutually exclusive with 'example'.",
				)
				.optional();
		},
		content: z
			.record(z.string(), MediaTypeObjectSchema)
			.meta({ title: "content" })
			.describe("For complex parameters. Mutually exclusive with 'schema'.")
			.optional(),
	})
	.meta({
		title: "Path Parameter",
		description: "A path parameter (/users/{id}).",
		examples: [
			{
				name: "id",
				in: "path",
				required: true,
				schema: { type: "integer" },
			},
		],
	});

const HeaderParameterObjectSchema = z
	.looseObject({
		name: z
			.string()
			.meta({
				title: "name",
				examples: [
					"X-Request-ID",
					"X-Correlation-ID",
					"Accept-Language",
				],
			})
			.describe("The name of the header. Case-insensitive."),
		in: z.literal("header").meta({ title: "in", examples: ["header"] }),
		description: z
			.string()
			.meta({
				title: "description",
				examples: ["Request correlation ID", "Preferred language"],
			})
			.describe(
				"A description of the parameter. CommonMark syntax MAY be used.",
			)
			.optional(),
		required: z
			.boolean()
			.optional()
			.describe("Whether the header is required.")
			.meta({ title: "required", examples: [true, false] }),
		deprecated: z
			.boolean()
			.optional()
			.describe("Marks the parameter as deprecated.")
			.meta({ title: "deprecated", examples: [true] }),
		style: HeaderParameterStyleSchema.describe(
			"Serialization style. Only 'simple' is valid for headers.",
		).optional(),
		explode: z
			.boolean()
			.meta({ title: "explode", examples: [true, false] })
			.describe(
				"For arrays/objects. When true, each value gets its own parameter.",
			)
			.optional(),
		get schema() {
			return SchemaObjectSchema.optional().meta({ title: "schema" });
		},
		example: z
			.unknown()
			.meta({
				title: "example",
				examples: ["en-US", "550e8400-e29b-41d4-a716-446655440000"],
			})
			.describe("Example value. Mutually exclusive with 'examples'.")
			.optional(),
		get examples() {
			return z
				.record(z.string(), ExampleSchema)
				.meta({ title: "examples" })
				.describe(
					"Multiple named examples. Mutually exclusive with 'example'.",
				)
				.optional();
		},
		content: z
			.record(z.string(), MediaTypeObjectSchema)
			.meta({ title: "content" })
			.describe("For complex parameters. Mutually exclusive with 'schema'.")
			.optional(),
	})
	.meta({
		title: "Header Parameter",
		description: "An HTTP header parameter.",
		examples: [
			{
				name: "X-Request-ID",
				in: "header",
				schema: { type: "string", format: "uuid" },
			},
		],
	});

const CookieParameterObjectSchema = z
	.looseObject({
		name: z
			.string()
			.meta({
				title: "name",
				examples: ["session_id", "csrf_token", "debug"],
			})
			.describe("The name of the cookie."),
		in: z.literal("cookie").meta({ title: "in", examples: ["cookie"] }),
		description: z
			.string()
			.meta({
				title: "description",
				examples: ["Session identifier", "CSRF protection token"],
			})
			.describe(
				"A description of the parameter. CommonMark syntax MAY be used.",
			)
			.optional(),
		required: z
			.boolean()
			.optional()
			.describe("Whether the cookie is required.")
			.meta({ title: "required", examples: [true, false] }),
		deprecated: z
			.boolean()
			.optional()
			.describe("Marks the parameter as deprecated.")
			.meta({ title: "deprecated", examples: [true] }),
		style: CookieParameterStyleSchema.describe(
			"Serialization style. Only 'form' is valid for cookies.",
		).optional(),
		explode: z
			.boolean()
			.meta({ title: "explode", examples: [true, false] })
			.describe(
				"For arrays/objects. When true, each value gets its own parameter.",
			)
			.optional(),
		get schema() {
			return SchemaObjectSchema.optional().meta({ title: "schema" });
		},
		example: z
			.unknown()
			.meta({ title: "example", examples: ["abc123", true] })
			.describe("Example value. Mutually exclusive with 'examples'.")
			.optional(),
		get examples() {
			return z
				.record(z.string(), ExampleSchema)
				.meta({ title: "examples" })
				.describe(
					"Multiple named examples. Mutually exclusive with 'example'.",
				)
				.optional();
		},
		content: z
			.record(z.string(), MediaTypeObjectSchema)
			.meta({ title: "content" })
			.describe("For complex parameters. Mutually exclusive with 'schema'.")
			.optional(),
	})
	.meta({
		title: "Cookie Parameter",
		description: "A cookie parameter.",
		examples: [
			{
				name: "session_id",
				in: "cookie",
				schema: { type: "string" },
			},
		],
	});

export const ParameterSchema = z
	.union([
		ReferenceObjectSchema,
		QueryParameterObjectSchema,
		PathParameterObjectSchema,
		HeaderParameterObjectSchema,
		CookieParameterObjectSchema,
	])
	.meta({
		title: "Parameter",
		description: "Describes a single operation parameter.",
	});

export type Parameter = z.infer<typeof ParameterSchema>;
