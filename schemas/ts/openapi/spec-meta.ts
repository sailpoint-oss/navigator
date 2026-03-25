import type { z } from "zod";

export type SpecVersion = "2.0" | "3.0" | "3.1" | "3.2";

const SPEC_BASE_URL: Record<SpecVersion, string> = {
	"2.0": "https://swagger.io/specification/v2/",
	// Canonical HTML documents for the selected patch versions:
	"3.0": "https://spec.openapis.org/oas/v3.0.4.html",
	"3.1": "https://spec.openapis.org/oas/v3.1.2.html",
	"3.2": "https://spec.openapis.org/oas/v3.2.0.html",
};

export function specLink(
	version: SpecVersion,
	anchor: string,
): {
	url: string;
	anchor: string;
} {
	const base = SPEC_BASE_URL[version];
	const normalized = anchor.startsWith("#") ? anchor.slice(1) : anchor;
	return {
		url: `${base}#${normalized}`,
		anchor: normalized,
	};
}

export function withSpec<T extends z.ZodType>(
	schema: T,
	version: SpecVersion,
	anchor: string,
): T {
	const link = specLink(version, anchor);
	// biome-ignore lint/suspicious/noExplicitAny: zod meta is untyped in many places
	const anySchema: any = schema as any;
	return anySchema.meta({
		...(anySchema._def?.meta ?? {}),
		spec: { version, ...link },
	});
}
