package navigator

import (
	"fmt"
	"strconv"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
)

func issueFromMetaOutputUnit(u *jsonschema.OutputUnit, fragment bool, idx *Index) Issue {
	ptr := u.InstanceLocation
	code, msg := metaFormatErrorKind(u.Error.Kind, ptr, fragment, idx)
	rng := FileStartRange
	if idx != nil && idx.IsArazzo() && idx.Arazzo != nil {
		rng = idx.Arazzo.Loc.Range
	} else if idx != nil && idx.Document != nil {
		rng = idx.Document.Loc.Range
	}
	return Issue{
		Code:     code,
		Message:  msg,
		Pointer:  ptr,
		Range:    rng,
		Severity: SeverityError,
		Category: CategoryMeta,
	}
}

func issueFromLeafError(le leafValidationError, fragment bool, idx *Index) Issue {
	// Use raw segments for context inference to avoid double-escaping issues
	ctx := inferOpenAPIContextFromSegments(le.Segments, fragment, idx)
	code, msg := metaFormatErrorKindWithContext(le.ErrorKind, ctx)
	rng := FileStartRange
	if idx != nil && idx.IsArazzo() && idx.Arazzo != nil {
		rng = idx.Arazzo.Loc.Range
	} else if idx != nil && idx.Document != nil {
		rng = idx.Document.Loc.Range
	}
	return Issue{
		Code:     code,
		Message:  msg,
		Pointer:  le.InstanceLocation,
		Range:    rng,
		Severity: SeverityError,
		Category: CategoryMeta,
	}
}

// inferOpenAPIContextFromSegments takes raw InstanceLocation segments (unescaped)
// and produces a human-readable context string.
func inferOpenAPIContextFromSegments(segs []string, fragment bool, idx *Index) string {
	if len(segs) == 0 {
		if fragment {
			return "this fragment"
		}
		if idx != nil && idx.IsArazzo() {
			return "the Arazzo document root"
		}
		return "the document root"
	}

	if idx != nil && idx.IsArazzo() {
		return inferArazzoContext(segs)
	}

	switch segs[0] {
	case "info":
		return "Info Object"
	case "externalDocs":
		return "External Documentation Object"
	case "servers":
		if len(segs) >= 2 {
			return fmt.Sprintf("Server Object (index %s)", segs[1])
		}
		return "the servers array"
	case "security":
		if len(segs) >= 2 {
			return fmt.Sprintf("Security Requirement (index %s)", segs[1])
		}
		return "the security array"
	case "tags":
		if len(segs) >= 2 {
			return fmt.Sprintf("Tag Object (index %s)", segs[1])
		}
		return "the tags array"
	case "paths":
		return inferPathContextFromSegments(segs[1:])
	case "webhooks":
		if len(segs) >= 2 {
			return fmt.Sprintf("Webhook '%s'", segs[1])
		}
		return "webhooks"
	case "components":
		return inferComponentContext(segs[1:])
	}

	if fragment {
		return "this fragment"
	}
	return "this document"
}

func inferPathContextFromSegments(segs []string) string {
	if len(segs) == 0 {
		return "paths"
	}
	pathTemplate := segs[0] // already unescaped
	if len(segs) == 1 {
		return fmt.Sprintf("Path Item '%s'", pathTemplate)
	}
	method := segs[1]
	if httpMethods[method] {
		ctx := fmt.Sprintf("Operation '%s %s'", strings.ToUpper(method), pathTemplate)
		if len(segs) >= 3 {
			return inferOperationSubContext(ctx, segs[2:])
		}
		return ctx
	}
	if method == "parameters" {
		if len(segs) >= 3 {
			return fmt.Sprintf("Parameter (index %s) in Path Item '%s'", segs[2], pathTemplate)
		}
		return fmt.Sprintf("parameters in Path Item '%s'", pathTemplate)
	}
	return fmt.Sprintf("Path Item '%s'", pathTemplate)
}

// metaFormatErrorKindWithContext formats an error using a pre-computed context string.
func metaFormatErrorKindWithContext(ek jsonschema.ErrorKind, ctx string) (code string, msg string) {
	switch k := ek.(type) {
	case *kind.Required:
		return "meta.required", metaMsgRequired(k, ctx)
	case *kind.Type:
		return "meta.type", metaMsgType(k, ctx)
	case *kind.Enum:
		return "meta.enum", fmt.Sprintf("Value is not valid in %s; must be one of the values permitted by the schema.", ctx)
	case *kind.OneOf:
		if len(k.Subschemas) == 0 {
			return "meta.one-of", fmt.Sprintf("Value does not match any allowed alternative in %s.", ctx)
		}
		return "meta.one-of", fmt.Sprintf("Value matches more than one alternative in %s, which is ambiguous.", ctx)
	case *kind.AnyOf:
		return "meta.any-of", fmt.Sprintf("Value does not satisfy any of the allowed options in %s.", ctx)
	case *kind.AllOf:
		return "meta.all-of", fmt.Sprintf("Value fails one or more combined requirements in %s.", ctx)
	case *kind.AdditionalProperties:
		return "meta.additional-property", metaMsgAdditionalProperties(k, ctx)
	case *kind.Format:
		return "meta.format", fmt.Sprintf("Expected valid %s format in %s.", k.Want, ctx)
	case *kind.Const:
		return "meta.const", fmt.Sprintf("Value must match the constant required by %s.", ctx)
	case *kind.Not:
		return "meta.not", fmt.Sprintf("Value must not match the forbidden pattern in %s.", ctx)
	case *kind.PropertyNames:
		return "meta.property-names", fmt.Sprintf("Property name %q is not allowed in %s.", k.Property, ctx)
	case *kind.MinProperties, *kind.MaxProperties:
		return "meta.property-count", fmt.Sprintf("Object does not meet property-count rules in %s.", ctx)
	case *kind.MinItems, *kind.MaxItems, *kind.AdditionalItems:
		return "meta.items", fmt.Sprintf("Array does not meet item rules in %s.", ctx)
	case *kind.UniqueItems:
		return "meta.unique-items", fmt.Sprintf("Array contains duplicate items in %s.", ctx)
	case *kind.Contains:
		return "meta.contains", fmt.Sprintf("Array must contain at least one matching element in %s.", ctx)
	case *kind.MinLength, *kind.MaxLength, *kind.Pattern:
		return "meta.string", fmt.Sprintf("String does not satisfy length or pattern rules in %s.", ctx)
	case *kind.Minimum, *kind.Maximum, *kind.ExclusiveMinimum, *kind.ExclusiveMaximum, *kind.MultipleOf:
		return "meta.number", fmt.Sprintf("Number is out of range in %s.", ctx)
	case *kind.DependentRequired, *kind.Dependency:
		return "meta.dependency", fmt.Sprintf("Missing related properties required by %s.", ctx)
	case *kind.InvalidJsonValue:
		return "meta.instance-type", fmt.Sprintf("Value has an invalid type for %s.", ctx)
	case *kind.FalseSchema:
		return "meta.false-schema", fmt.Sprintf("Value is not permitted in %s.", ctx)
	case *kind.Reference:
		return "meta.ref", fmt.Sprintf("Validation failed in %s (see nested causes).", ctx)
	case *kind.Group, *kind.Schema:
		return "meta.group", fmt.Sprintf("Value failed meta-schema validation in %s.", ctx)
	default:
		return "meta.catchall", fmt.Sprintf("Value does not satisfy the published meta-schema for %s. (keyword: %s)", ctx, metaKeywordHint(ek))
	}
}

func metaFormatErrorKind(ek jsonschema.ErrorKind, instPtr string, fragment bool, idx *Index) (code string, msg string) {
	ctx := inferOpenAPIContext(instPtr, fragment, idx)

	switch k := ek.(type) {
	case *kind.Required:
		return "meta.required", metaMsgRequired(k, ctx)
	case *kind.Type:
		return "meta.type", metaMsgType(k, ctx)
	case *kind.Enum:
		return "meta.enum", fmt.Sprintf("Value is not valid in %s; must be one of the values permitted by the schema.", ctx)
	case *kind.OneOf:
		if len(k.Subschemas) == 0 {
			return "meta.one-of", fmt.Sprintf("Value does not match any allowed alternative in %s.", ctx)
		}
		return "meta.one-of", fmt.Sprintf("Value matches more than one alternative in %s, which is ambiguous.", ctx)
	case *kind.AnyOf:
		return "meta.any-of", fmt.Sprintf("Value does not satisfy any of the allowed options in %s.", ctx)
	case *kind.AllOf:
		return "meta.all-of", fmt.Sprintf("Value fails one or more combined requirements in %s.", ctx)
	case *kind.AdditionalProperties:
		return "meta.additional-property", metaMsgAdditionalProperties(k, ctx)
	case *kind.Format:
		return "meta.format", fmt.Sprintf("Expected valid %s format in %s.", k.Want, ctx)
	case *kind.Const:
		return "meta.const", fmt.Sprintf("Value must match the constant required by %s.", ctx)
	case *kind.Not:
		return "meta.not", fmt.Sprintf("Value must not match the forbidden pattern in %s.", ctx)
	case *kind.PropertyNames:
		return "meta.property-names", fmt.Sprintf("Property name %q is not allowed in %s.", k.Property, ctx)
	case *kind.MinProperties, *kind.MaxProperties:
		return "meta.property-count", fmt.Sprintf("Object does not meet property-count rules in %s.", ctx)
	case *kind.MinItems, *kind.MaxItems, *kind.AdditionalItems:
		return "meta.items", fmt.Sprintf("Array does not meet item rules in %s.", ctx)
	case *kind.UniqueItems:
		return "meta.unique-items", fmt.Sprintf("Array contains duplicate items in %s.", ctx)
	case *kind.Contains:
		return "meta.contains", fmt.Sprintf("Array must contain at least one matching element in %s.", ctx)
	case *kind.MinLength, *kind.MaxLength, *kind.Pattern:
		return "meta.string", fmt.Sprintf("String does not satisfy length or pattern rules in %s.", ctx)
	case *kind.Minimum, *kind.Maximum, *kind.ExclusiveMinimum, *kind.ExclusiveMaximum, *kind.MultipleOf:
		return "meta.number", fmt.Sprintf("Number is out of range in %s.", ctx)
	case *kind.DependentRequired, *kind.Dependency:
		return "meta.dependency", fmt.Sprintf("Missing related properties required by %s.", ctx)
	case *kind.InvalidJsonValue:
		return "meta.instance-type", fmt.Sprintf("Value has an invalid type for %s.", ctx)
	case *kind.FalseSchema:
		return "meta.false-schema", fmt.Sprintf("Value is not permitted in %s.", ctx)
	case *kind.Reference:
		return "meta.ref", fmt.Sprintf("Validation failed in %s (see nested causes).", ctx)
	case *kind.Group, *kind.Schema:
		return "meta.group", fmt.Sprintf("Value failed meta-schema validation in %s.", ctx)
	default:
		return "meta.catchall", fmt.Sprintf("Value does not satisfy the published meta-schema for %s. (keyword: %s)", ctx, metaKeywordHint(ek))
	}
}

func metaKeywordHint(ek jsonschema.ErrorKind) string {
	if ek == nil {
		return "unknown"
	}
	kp := ek.KeywordPath()
	if len(kp) == 0 {
		return strings.TrimPrefix(fmt.Sprintf("%T", ek), "*kind.")
	}
	return strings.Join(kp, ".")
}

func metaMsgRequired(k *kind.Required, ctx string) string {
	if len(k.Missing) == 1 {
		return fmt.Sprintf("%s requires property '%s'.", ctx, k.Missing[0])
	}
	return fmt.Sprintf("%s is missing required properties: %s.", ctx, joinSingleQuoted(k.Missing))
}

func metaMsgType(k *kind.Type, ctx string) string {
	want := strings.Join(k.Want, " or ")
	return fmt.Sprintf("Expected %s in %s, got %s.", want, ctx, k.Got)
}

func metaMsgAdditionalProperties(k *kind.AdditionalProperties, ctx string) string {
	if len(k.Properties) == 1 {
		return fmt.Sprintf("'%s' is not a valid property in %s.", k.Properties[0], ctx)
	}
	return fmt.Sprintf("Properties %s are not valid in %s.", joinSingleQuoted(k.Properties), ctx)
}

// joinSingleQuoted formats strings as 'a', 'b', 'c'.
func joinSingleQuoted(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("'" + strs[0] + "'")
	for _, s := range strs[1:] {
		b.WriteString(", '")
		b.WriteString(s)
		b.WriteString("'")
	}
	return b.String()
}

// joinQuotedMeta formats strings as "a", "b", "c" (for backward compat).
func joinQuotedMeta(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(strconv.Quote(strs[0]))
	for _, s := range strs[1:] {
		b.WriteString(", ")
		b.WriteString(strconv.Quote(s))
	}
	return b.String()
}

func collectMetaOutputUnits(out *jsonschema.OutputUnit, acc *[]*jsonschema.OutputUnit) {
	if out == nil {
		return
	}
	if out.Error != nil {
		*acc = append(*acc, out)
	}
	for i := range out.Errors {
		collectMetaOutputUnits(&out.Errors[i], acc)
	}
}

// leafValidationError holds a flattened leaf error from the ValidationError tree.
type leafValidationError struct {
	InstanceLocation string   // JSON Pointer (escaped)
	Segments         []string // raw path segments (unescaped)
	ErrorKind        jsonschema.ErrorKind
}

// collectLeafValidationErrors walks the ValidationError tree and returns only
// actionable leaf errors — errors that have no children, or whose ErrorKind
// is a specific constraint violation (AdditionalProperties, Required, Type, etc.)
// rather than a structural wrapper (Reference, Group, Schema).
// Results are deduplicated by (pointer, code, message-key).
func collectLeafValidationErrors(ve *jsonschema.ValidationError) []leafValidationError {
	if ve == nil {
		return nil
	}
	var results []leafValidationError
	walkValidationErrors(ve, &results)
	return deduplicateLeafErrors(results)
}

func deduplicateLeafErrors(errors []leafValidationError) []leafValidationError {
	seen := make(map[string]bool, len(errors))
	out := make([]leafValidationError, 0, len(errors))
	for _, e := range errors {
		// Build a dedup key from pointer + error kind type + error detail
		key := e.InstanceLocation + "|" + fmt.Sprintf("%T", e.ErrorKind) + "|" + errorDetail(e.ErrorKind)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, e)
	}
	return out
}

func errorDetail(ek jsonschema.ErrorKind) string {
	switch k := ek.(type) {
	case *kind.AdditionalProperties:
		return strings.Join(k.Properties, ",")
	case *kind.Required:
		return strings.Join(k.Missing, ",")
	case *kind.Type:
		return k.Got
	case *kind.Const:
		return "const"
	default:
		return ""
	}
}

func escapePointerSegments(segs []string) string {
	if len(segs) == 0 {
		return ""
	}
	var b strings.Builder
	for _, seg := range segs {
		b.WriteByte('/')
		// JSON Pointer escaping: ~ → ~0, / → ~1
		seg = strings.ReplaceAll(seg, "~", "~0")
		seg = strings.ReplaceAll(seg, "/", "~1")
		b.WriteString(seg)
	}
	return b.String()
}

func walkValidationErrors(ve *jsonschema.ValidationError, acc *[]leafValidationError) {
	if ve == nil {
		return
	}

	// Structural wrappers: skip and recurse into children.
	switch ve.ErrorKind.(type) {
	case *kind.Reference, *kind.Group, *kind.Schema:
		if len(ve.Causes) > 0 {
			for _, cause := range ve.Causes {
				walkValidationErrors(cause, acc)
			}
			return
		}
	}

	// Composition keywords (oneOf, anyOf): pick the BEST matching alternative
	// (fewest errors) and only report its errors. This avoids contradictory
	// messages from alternatives that don't match the document.
	switch ve.ErrorKind.(type) {
	case *kind.OneOf, *kind.AnyOf:
		if len(ve.Causes) > 0 {
			walkBestAlternative(ve.Causes, acc)
			return
		}
	}

	// allOf: all branches must pass, so recurse into all children.
	if _, ok := ve.ErrorKind.(*kind.AllOf); ok && len(ve.Causes) > 0 {
		for _, cause := range ve.Causes {
			walkValidationErrors(cause, acc)
		}
		return
	}

	le := leafValidationError{
		InstanceLocation: escapePointerSegments(ve.InstanceLocation),
		Segments:         ve.InstanceLocation,
		ErrorKind:        ve.ErrorKind,
	}

	if len(ve.Causes) == 0 {
		*acc = append(*acc, le)
		return
	}

	// Other non-leaf nodes: recurse into children.
	for _, cause := range ve.Causes {
		walkValidationErrors(cause, acc)
	}
}

// walkBestAlternative picks the oneOf/anyOf alternative with the fewest leaf
// errors (closest match) and only reports its errors. This prevents
// contradictory messages like "'$ref' is not valid" alongside "requires '$ref'".
func walkBestAlternative(causes []*jsonschema.ValidationError, acc *[]leafValidationError) {
	if len(causes) == 0 {
		return
	}

	var bestErrors []leafValidationError
	bestCount := -1

	for _, cause := range causes {
		var trial []leafValidationError
		walkValidationErrors(cause, &trial)
		if bestCount < 0 || len(trial) < bestCount {
			bestErrors = trial
			bestCount = len(trial)
		}
	}

	*acc = append(*acc, bestErrors...)
}

// ---------------------------------------------------------------------------
// OpenAPI context inference from JSON Pointer
// ---------------------------------------------------------------------------

var httpMethods = map[string]bool{
	"get": true, "put": true, "post": true, "delete": true,
	"patch": true, "options": true, "head": true, "trace": true,
}

// inferOpenAPIContext maps a JSON Pointer to a human-readable OpenAPI context
// string like "Schema Object 'Pet'" or "Operation 'GET /pets'".
func inferOpenAPIContext(instPtr string, fragment bool, idx *Index) string {
	segs := splitPointerSegments(instPtr)
	if len(segs) == 0 {
		if fragment {
			return "this fragment"
		}
		if idx != nil && idx.IsArazzo() {
			return "the Arazzo document root"
		}
		return "the document root"
	}

	// Arazzo documents have different structure
	if idx != nil && idx.IsArazzo() {
		return inferArazzoContext(segs)
	}

	switch segs[0] {
	case "info":
		return "Info Object"
	case "externalDocs":
		return "External Documentation Object"
	case "servers":
		if len(segs) >= 2 {
			return fmt.Sprintf("Server Object (index %s)", segs[1])
		}
		return "the servers array"
	case "security":
		if len(segs) >= 2 {
			return fmt.Sprintf("Security Requirement (index %s)", segs[1])
		}
		return "the security array"
	case "tags":
		if len(segs) >= 2 {
			return fmt.Sprintf("Tag Object (index %s)", segs[1])
		}
		return "the tags array"
	case "paths":
		return inferPathContext(segs[1:])
	case "webhooks":
		if len(segs) >= 2 {
			return fmt.Sprintf("Webhook '%s'", unescapePointerSegment(segs[1]))
		}
		return "webhooks"
	case "components":
		return inferComponentContext(segs[1:])
	}

	if fragment {
		return "this fragment"
	}
	return "this document"
}

func inferPathContext(segs []string) string {
	if len(segs) == 0 {
		return "paths"
	}
	pathTemplate := unescapePointerSegment(segs[0])
	if len(segs) == 1 {
		return fmt.Sprintf("Path Item '%s'", pathTemplate)
	}
	method := segs[1]
	if httpMethods[method] {
		ctx := fmt.Sprintf("Operation '%s %s'", strings.ToUpper(method), pathTemplate)
		if len(segs) >= 3 {
			return inferOperationSubContext(ctx, segs[2:])
		}
		return ctx
	}
	if method == "parameters" {
		if len(segs) >= 3 {
			return fmt.Sprintf("Parameter (index %s) in Path Item '%s'", segs[2], pathTemplate)
		}
		return fmt.Sprintf("parameters in Path Item '%s'", pathTemplate)
	}
	return fmt.Sprintf("Path Item '%s'", pathTemplate)
}

func inferOperationSubContext(opCtx string, segs []string) string {
	switch segs[0] {
	case "parameters":
		if len(segs) >= 2 {
			return fmt.Sprintf("Parameter (index %s) in %s", segs[1], opCtx)
		}
		return fmt.Sprintf("parameters in %s", opCtx)
	case "responses":
		if len(segs) >= 2 {
			return fmt.Sprintf("Response '%s' in %s", segs[1], opCtx)
		}
		return fmt.Sprintf("responses in %s", opCtx)
	case "requestBody":
		return fmt.Sprintf("Request Body in %s", opCtx)
	case "security":
		return fmt.Sprintf("security in %s", opCtx)
	case "callbacks":
		if len(segs) >= 2 {
			return fmt.Sprintf("Callback '%s' in %s", segs[1], opCtx)
		}
		return fmt.Sprintf("callbacks in %s", opCtx)
	}
	return opCtx
}

func inferComponentContext(segs []string) string {
	if len(segs) == 0 {
		return "components"
	}
	kind := segs[0]
	if len(segs) < 2 {
		return fmt.Sprintf("components/%s", kind)
	}
	name := segs[1]
	switch kind {
	case "schemas":
		return fmt.Sprintf("Schema Object '%s'", name)
	case "responses":
		return fmt.Sprintf("Response Object '%s'", name)
	case "parameters":
		return fmt.Sprintf("Parameter Object '%s'", name)
	case "requestBodies":
		return fmt.Sprintf("Request Body '%s'", name)
	case "headers":
		return fmt.Sprintf("Header Object '%s'", name)
	case "securitySchemes":
		return fmt.Sprintf("Security Scheme '%s'", name)
	case "links":
		return fmt.Sprintf("Link Object '%s'", name)
	case "callbacks":
		return fmt.Sprintf("Callback Object '%s'", name)
	case "examples":
		return fmt.Sprintf("Example Object '%s'", name)
	case "pathItems":
		return fmt.Sprintf("Path Item Object '%s'", name)
	}
	return fmt.Sprintf("components/%s/%s", kind, name)
}

func inferArazzoContext(segs []string) string {
	switch segs[0] {
	case "info":
		return "Arazzo Info Object"
	case "sourceDescriptions":
		if len(segs) >= 2 {
			return fmt.Sprintf("Source Description (index %s)", segs[1])
		}
		return "sourceDescriptions"
	case "workflows":
		if len(segs) >= 2 {
			return fmt.Sprintf("Workflow (index %s)", segs[1])
		}
		return "workflows"
	}
	return "this Arazzo document"
}

// splitPointerSegments splits a JSON Pointer into its path segments.
func splitPointerSegments(ptr string) []string {
	if ptr == "" || ptr == "/" {
		return nil
	}
	return strings.Split(strings.TrimPrefix(ptr, "/"), "/")
}

// unescapePointerSegment reverses JSON Pointer escaping (~1 → /, ~0 → ~).
func unescapePointerSegment(s string) string {
	s = strings.ReplaceAll(s, "~1", "/")
	s = strings.ReplaceAll(s, "~0", "~")
	return s
}
