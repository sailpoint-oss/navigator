import z from "zod";
import { PathItemObjectSchema } from "./path-item";

export const PathsObjectSchema = z
	.record(
		z.string().refine((k) => k.startsWith("/") || k.startsWith("x-"), {
			message: "Paths Object keys must start with '/' (path template) or 'x-'.",
		}),
		PathItemObjectSchema,
	)
	.meta({
		title: "Paths",
		description: "Holds the relative paths to the individual endpoints.",
		examples: [{ "/pets": { get: { responses: { "200": { description: "OK" } } } } }],
	});

export type PathsObject = z.infer<typeof PathsObjectSchema>;


