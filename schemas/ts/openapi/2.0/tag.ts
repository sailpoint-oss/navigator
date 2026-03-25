import z from "zod";
import { ExternalDocumentationObjectSchema } from "./external-documentation";

export const TagObjectSchema = z
	.looseObject({
		name: z.string().describe("REQUIRED. The name of the tag.").meta({ title: "name" }),
		description: z.string().optional().meta({ title: "description" }),
		externalDocs: ExternalDocumentationObjectSchema.optional().meta({
			title: "externalDocs",
		}),
	})
	.meta({ title: "Tag" });

export type TagObject = z.infer<typeof TagObjectSchema>;


