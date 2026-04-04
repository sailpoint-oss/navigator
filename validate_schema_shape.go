package navigator

import (
	"fmt"
	"strconv"
)

var wellKnownSchemaTypes = map[string]struct{}{
	"array": {}, "boolean": {}, "integer": {}, "number": {},
	"object": {}, "string": {}, "null": {},
}

func validateSchemaShape(idx *Index, sink *issueSink) {
	if idx == nil || idx.Document == nil {
		return
	}
	if idx.Document.DocType == DocTypeNonOpenAPI {
		return
	}
	doc := idx.Document
	treeBacked := idx.Tree() != nil

	for path, item := range doc.Paths {
		if item == nil {
			continue
		}
		for i, p := range item.Parameters {
			ptr := "/paths/" + EscapeJSONPointer(path) + "/parameters/" + strconv.Itoa(i)
			validateParameterShape(p, ptr, sink)
		}
		for _, mo := range item.Operations() {
			if mo.Operation == nil {
				continue
			}
			op := mo.Operation
			for i, p := range op.Parameters {
				ptr := "/paths/" + EscapeJSONPointer(path) + "/" + mo.Method + "/parameters/" + strconv.Itoa(i)
				validateParameterShape(p, ptr, sink)
			}
		}
	}

	if doc.Components == nil {
		return
	}
	c := doc.Components

	for name, p := range c.Parameters {
		ptr := "/components/parameters/" + EscapeJSONPointer(name)
		validateParameterShape(p, ptr, sink)
	}

	for name, s := range c.Schemas {
		ptr := "/components/schemas/" + EscapeJSONPointer(name)
		validateSchemaTypeShape(s, ptr, sink, 0)
	}

	for name, ss := range c.SecuritySchemes {
		ptr := "/components/securitySchemes/" + EscapeJSONPointer(name)
		validateSecuritySchemeShape(ss, ptr, sink, treeBacked)
	}
}

func validateParameterShape(p *Parameter, ptr string, sink *issueSink) {
	if p == nil || p.Ref != "" {
		return
	}
	if p.Name == "" {
		rng := p.NameLoc.Range
		if IsEmpty(rng) {
			rng = p.Loc.Range
		}
		sink.add(Issue{
			Code:     "schema.parameter-missing-name",
			Message:  "Parameter Object must include 'name' when not a $ref",
			Pointer:  ptr + "/name",
			Range:    rng,
			Severity: SeverityError,
			Category: CategorySchema,
		})
	}
	if p.In == "" {
		sink.add(Issue{
			Code:     "schema.parameter-missing-in",
			Message:  "Parameter Object must include 'in' (query, path, header, or cookie) when not a $ref",
			Pointer:  ptr + "/in",
			Range:    p.Loc.Range,
			Severity: SeverityError,
			Category: CategorySchema,
		})
	}
}

func validateSchemaTypeShape(s *Schema, ptr string, sink *issueSink, depth int) {
	if s == nil || depth > 64 {
		return
	}
	if s.Ref != "" {
		return
	}
	if s.Type != "" {
		if _, ok := wellKnownSchemaTypes[s.Type]; !ok {
			sink.add(Issue{
				Code:     "schema.unknown-type",
				Message:  fmt.Sprintf("unknown schema type %q", s.Type),
				Pointer:  ptr + "/type",
				Range:    s.Loc.Range,
				Severity: SeverityWarning,
				Category: CategorySchema,
			})
		}
	}
	if s.Items != nil {
		validateSchemaTypeShape(s.Items, ptr+"/items", sink, depth+1)
	}
	for name, ps := range s.Properties {
		validateSchemaTypeShape(ps, ptr+"/properties/"+EscapeJSONPointer(name), sink, depth+1)
	}
	if s.AdditionalProperties != nil {
		validateSchemaTypeShape(s.AdditionalProperties, ptr+"/additionalProperties", sink, depth+1)
	}
	if s.UnevaluatedProperties != nil {
		validateSchemaTypeShape(s.UnevaluatedProperties, ptr+"/unevaluatedProperties", sink, depth+1)
	}
	for i, sub := range s.AllOf {
		validateSchemaTypeShape(sub, ptr+"/allOf/"+strconv.Itoa(i), sink, depth+1)
	}
	for i, sub := range s.AnyOf {
		validateSchemaTypeShape(sub, ptr+"/anyOf/"+strconv.Itoa(i), sink, depth+1)
	}
	for i, sub := range s.OneOf {
		validateSchemaTypeShape(sub, ptr+"/oneOf/"+strconv.Itoa(i), sink, depth+1)
	}
	if s.Not != nil {
		validateSchemaTypeShape(s.Not, ptr+"/not", sink, depth+1)
	}
}

func validateSecuritySchemeShape(ss *SecurityScheme, ptr string, sink *issueSink, treeBacked bool) {
	if ss == nil || ss.Ref != "" {
		return
	}
	if ss.Type == "" {
		sink.add(Issue{
			Code:     "schema.security-scheme-missing-type",
			Message:  "Security Scheme must declare 'type'",
			Pointer:  ptr + "/type",
			Range:    ss.Loc.Range,
			Severity: SeverityError,
			Category: CategorySchema,
		})
		return
	}
	switch ss.Type {
	case "apiKey":
		if ss.Name == "" {
			sink.add(Issue{
				Code:     "schema.security-scheme-apikey-fields",
				Message:  "API Key Security Scheme requires 'name'",
				Pointer:  ptr + "/name",
				Range:    ss.Loc.Range,
				Severity: SeverityError,
				Category: CategorySchema,
			})
		}
		if ss.In == "" {
			sink.add(Issue{
				Code:     "schema.security-scheme-apikey-fields",
				Message:  "API Key Security Scheme requires 'in'",
				Pointer:  ptr + "/in",
				Range:    ss.Loc.Range,
				Severity: SeverityError,
				Category: CategorySchema,
			})
		}
	case "http":
		if ss.Scheme == "" {
			sink.add(Issue{
				Code:     "schema.security-scheme-http-fields",
				Message:  "HTTP Security Scheme requires 'scheme'",
				Pointer:  ptr + "/scheme",
				Range:    ss.Loc.Range,
				Severity: SeverityError,
				Category: CategorySchema,
			})
		}
	case "oauth2":
		if treeBacked {
			return
		}
		if ss.Flows == nil {
			sink.add(Issue{
				Code:     "schema.security-scheme-oauth2-fields",
				Message:  "OAuth2 Security Scheme requires 'flows'",
				Pointer:  ptr + "/flows",
				Range:    ss.Loc.Range,
				Severity: SeverityError,
				Category: CategorySchema,
			})
		}
	case "openIdConnect":
		if treeBacked {
			return
		}
		if ss.OpenIDConnectURL == "" {
			sink.add(Issue{
				Code:     "schema.security-scheme-openid-fields",
				Message:  "OpenID Connect Security Scheme requires 'openIdConnectUrl'",
				Pointer:  ptr + "/openIdConnectUrl",
				Range:    ss.Loc.Range,
				Severity: SeverityError,
				Category: CategorySchema,
			})
		}
	case "mutualTLS":
		// no extra required fields
	default:
		sink.add(Issue{
			Code:     "schema.security-scheme-unknown-type",
			Message:  fmt.Sprintf("unknown security scheme type %q", ss.Type),
			Pointer:  ptr + "/type",
			Range:    ss.Loc.Range,
			Severity: SeverityWarning,
			Category: CategorySchema,
		})
	}
}
