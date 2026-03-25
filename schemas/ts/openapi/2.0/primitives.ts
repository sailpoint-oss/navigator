import z from "zod";

export const UrlStringSchema = z
	.string()
	.meta({ title: "UrlString", examples: ["https://example.com/docs"] })
	.describe("A URL string.");

export const MimeTypeStringSchema = z
	.string()
	.meta({ title: "MimeTypeString", examples: ["application/json"] })
	.describe("A MIME type string (RFC 6838).");


