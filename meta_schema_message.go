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

func metaFormatErrorKind(ek jsonschema.ErrorKind, instPtr string, fragment bool, idx *Index) (code string, msg string) {
	docHint := "OpenAPI document"
	if idx != nil && idx.IsArazzo() {
		docHint = "Arazzo document"
	} else if fragment {
		docHint = "OpenAPI fragment file"
	}
	at := instPtr
	if at == "" {
		at = "(document root)"
	}

	switch k := ek.(type) {
	case *kind.Required:
		return "meta.required", metaMsgRequired(k, at, docHint)
	case *kind.Type:
		return "meta.type", metaMsgType(k, at, docHint)
	case *kind.Enum:
		return "meta.enum", fmt.Sprintf("Value at %s is not allowed for this %s. It must be one of the values permitted by the API description schema.", at, docHint)
	case *kind.OneOf:
		if len(k.Subschemas) == 0 {
			return "meta.one-of", fmt.Sprintf("Value at %s does not match any allowed alternative for this %s. Adjust the structure to satisfy one of the expected shapes.", at, docHint)
		}
		return "meta.one-of", fmt.Sprintf("Value at %s matches more than one alternative in the schema for this %s, which is ambiguous. Simplify the value so only one alternative applies.", at, docHint)
	case *kind.AnyOf:
		return "meta.any-of", fmt.Sprintf("Value at %s does not satisfy any of the allowed options for this %s.", at, docHint)
	case *kind.AllOf:
		return "meta.all-of", fmt.Sprintf("Value at %s fails one or more combined requirements for this %s.", at, docHint)
	case *kind.AdditionalProperties:
		return "meta.additional-property", metaMsgAdditionalProperties(k, at, docHint)
	case *kind.Format:
		return "meta.format", fmt.Sprintf("Value at %s is not a valid %q for this %s.", at, k.Want, docHint)
	case *kind.Const:
		return "meta.const", fmt.Sprintf("Value at %s must match the constant required by this %s.", at, docHint)
	case *kind.Not:
		return "meta.not", fmt.Sprintf("Value at %s must not match the forbidden pattern for this %s.", at, docHint)
	case *kind.PropertyNames:
		return "meta.property-names", fmt.Sprintf("Property name %q at %s is not allowed for this %s.", k.Property, at, docHint)
	case *kind.MinProperties, *kind.MaxProperties:
		return "meta.property-count", fmt.Sprintf("Object at %s does not meet property-count rules for this %s.", at, docHint)
	case *kind.MinItems, *kind.MaxItems, *kind.AdditionalItems:
		return "meta.items", fmt.Sprintf("Array at %s does not meet item rules for this %s.", at, docHint)
	case *kind.UniqueItems:
		return "meta.unique-items", fmt.Sprintf("Array at %s contains duplicate items, which is not allowed for this %s.", at, docHint)
	case *kind.Contains:
		return "meta.contains", fmt.Sprintf("Array at %s must contain at least one element matching the required shape for this %s.", at, docHint)
	case *kind.MinLength, *kind.MaxLength, *kind.Pattern:
		return "meta.string", fmt.Sprintf("String at %s does not satisfy length or pattern rules for this %s.", at, docHint)
	case *kind.Minimum, *kind.Maximum, *kind.ExclusiveMinimum, *kind.ExclusiveMaximum, *kind.MultipleOf:
		return "meta.number", fmt.Sprintf("Number at %s is out of range for this %s.", at, docHint)
	case *kind.DependentRequired, *kind.Dependency:
		return "meta.dependency", fmt.Sprintf("Value at %s is missing related properties required by this %s.", at, docHint)
	case *kind.InvalidJsonValue:
		return "meta.instance-type", fmt.Sprintf("Value at %s has a type that cannot be validated as JSON data in this %s.", at, docHint)
	case *kind.FalseSchema:
		return "meta.false-schema", fmt.Sprintf("Value at %s is not permitted by this %s.", at, docHint)
	case *kind.Reference:
		return "meta.ref", fmt.Sprintf("Validation failed at %s for this %s (see nested causes).", at, docHint)
	case *kind.Group, *kind.Schema:
		return "meta.group", fmt.Sprintf("Value at %s failed meta-schema validation for this %s.", at, docHint)
	default:
		return "meta.catchall", fmt.Sprintf("Value at %s does not satisfy the published meta-schema for this %s. (keyword: %s)", at, docHint, metaKeywordHint(ek))
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

func metaMsgRequired(k *kind.Required, at, docHint string) string {
	if len(k.Missing) == 1 {
		return fmt.Sprintf("Missing required property %s at %s (%s).", strconv.Quote(k.Missing[0]), at, docHint)
	}
	return fmt.Sprintf("Missing required properties at %s (%s): %s.", at, docHint, joinQuotedMeta(k.Missing))
}

func metaMsgType(k *kind.Type, at, docHint string) string {
	want := strings.Join(k.Want, ", ")
	return fmt.Sprintf("Value at %s must be %s for this %s, but the document has type %q.", at, want, docHint, k.Got)
}

func metaMsgAdditionalProperties(k *kind.AdditionalProperties, at, docHint string) string {
	if len(k.Properties) == 1 {
		return fmt.Sprintf("Property %s at %s is not allowed for this %s.", strconv.Quote(k.Properties[0]), at, docHint)
	}
	return fmt.Sprintf("Properties %s at %s are not allowed for this %s.", joinQuotedMeta(k.Properties), at, docHint)
}

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
