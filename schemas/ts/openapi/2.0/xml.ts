import z from "zod";

export const XmlObjectSchema = z
	.looseObject({
		name: z.string().optional().meta({ title: "name" }),
		namespace: z.string().optional().meta({ title: "namespace" }),
		prefix: z.string().optional().meta({ title: "prefix" }),
		attribute: z.boolean().optional().meta({ title: "attribute" }),
		wrapped: z.boolean().optional().meta({ title: "wrapped" }),
	})
	.meta({ title: "XML" });

export type XmlObject = z.infer<typeof XmlObjectSchema>;


