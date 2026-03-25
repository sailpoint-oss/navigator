import z from "zod";
import { PathItemObjectSchema } from "./path-item";

export const PathsObjectSchema = z
	.record(z.string(), z.union([PathItemObjectSchema, z.unknown()]))
	.meta({ title: "Paths" });

export type PathsObject = z.infer<typeof PathsObjectSchema>;
