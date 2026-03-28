package navigator

import "context"

// Issue code registry — stable identifiers consumed by downstream (barrelman
// LSP diagnostics). Adding new codes is non-breaking; renaming or removing
// existing codes is a breaking change.
//
// Syntax (validate_syntax.go):
//
//	syntax.parse-error                    tree-sitter ERROR node or root HasError
//	syntax.duplicate-key                  duplicate mapping key in YAML/JSON
//
// Structural (validate_structural.go):
//
//	structural.root-not-mapping           document root is not a YAML/JSON object
//	structural.ir-wrong-type              top-level key has unexpected YAML kind
//	structural.unsupported-version        openapi/swagger version not recognized
//	structural.unsupported-arazzo-version arazzo version not recognized
//	structural.missing-info               info object absent (OAS 3.x / Swagger 2.0)
//	structural.missing-info-title         info.title required but missing
//	structural.missing-info-version       info.version required but missing
//	structural.missing-source-descriptions arazzo sourceDescriptions absent
//	structural.missing-workflows          arazzo workflows absent
//	structural.missing-workflow-id        arazzo workflowId required but missing
//	structural.missing-step-id            arazzo stepId required but missing
//	structural.duplicate-workflow-id      arazzo workflowId must be unique
//	structural.duplicate-step-id          arazzo stepId must be unique per workflow
//	structural.step-target-count          arazzo step must target exactly one next action
//	structural.missing-responses          operation has no responses object
//
// Schema-shape (validate_schema_shape.go):
//
//	schema.parameter-missing-name         non-$ref parameter lacks name
//	schema.parameter-missing-in           non-$ref parameter lacks in
//	schema.unknown-type                   schema type not in OAS type set
//	schema.security-scheme-missing-type   securityScheme has no type
//	schema.security-scheme-unknown-type   securityScheme type not recognized
//	schema.security-scheme-apikey-fields  apiKey scheme missing name or in
//	schema.security-scheme-http-fields    http scheme missing scheme field
//	schema.security-scheme-oauth2-fields  oauth2 scheme missing flows
//	schema.security-scheme-openid-fields  openIdConnect scheme missing openIdConnectUrl
//
// Meta-schema (validate_meta_schema.go, Draft 2020-12 via jsonschema/v6):
//
//	meta.compiler                       embedded meta-schemas failed to compile
//	meta.decode                         YAML/JSON decode failed before meta validation
//	meta.required                       required property missing
//	meta.type                           JSON type mismatch
//	meta.enum                           enum mismatch
//	meta.one-of                         oneOf mismatch
//	meta.any-of                         anyOf mismatch
//	meta.all-of                         allOf mismatch
//	meta.additional-property            additionalProperties violation
//	meta.format                         format assertion failed
//	meta.const                          const mismatch
//	meta.not                            not keyword failed
//	meta.property-names                 propertyNames failed
//	meta.property-count                 minProperties / maxProperties
//	meta.items                          array item constraints
//	meta.unique-items                   uniqueItems violation
//	meta.contains                       contains keyword
//	meta.string                         string constraints
//	meta.number                         numeric constraints
//	meta.dependency                     dependencies / dependentRequired
//	meta.instance-type                  invalid JSON value type
//	meta.false-schema                   false schema
//	meta.ref                            $ref validation wrapper
//	meta.group                          grouped / schema-level wrapper
//	meta.catchall                       unclassified keyword
//	meta.truncated                      MaxIssues cap hit during meta pass

// knownIssueCodes is the canonical set of Issue.Code values. Tests assert that
// no validation pass emits a code outside this set.
var knownIssueCodes = map[string]struct{}{
	"syntax.parse-error":                     {},
	"syntax.duplicate-key":                   {},
	"structural.root-not-mapping":            {},
	"structural.ir-wrong-type":               {},
	"structural.unsupported-version":         {},
	"structural.unsupported-arazzo-version":  {},
	"structural.missing-info":                {},
	"structural.missing-info-title":          {},
	"structural.missing-info-version":        {},
	"structural.missing-source-descriptions": {},
	"structural.missing-workflows":           {},
	"structural.missing-workflow-id":         {},
	"structural.missing-step-id":             {},
	"structural.duplicate-workflow-id":       {},
	"structural.duplicate-step-id":           {},
	"structural.step-target-count":           {},
	"structural.missing-responses":           {},
	"schema.parameter-missing-name":          {},
	"schema.parameter-missing-in":            {},
	"schema.unknown-type":                    {},
	"schema.security-scheme-missing-type":    {},
	"schema.security-scheme-unknown-type":    {},
	"schema.security-scheme-apikey-fields":   {},
	"schema.security-scheme-http-fields":     {},
	"schema.security-scheme-oauth2-fields":   {},
	"schema.security-scheme-openid-fields":   {},
	"meta.compiler":                          {},
	"meta.decode":                            {},
	"meta.required":                          {},
	"meta.type":                              {},
	"meta.enum":                              {},
	"meta.one-of":                            {},
	"meta.any-of":                            {},
	"meta.all-of":                            {},
	"meta.additional-property":               {},
	"meta.format":                            {},
	"meta.const":                             {},
	"meta.not":                               {},
	"meta.property-names":                    {},
	"meta.property-count":                    {},
	"meta.items":                             {},
	"meta.unique-items":                      {},
	"meta.contains":                          {},
	"meta.string":                            {},
	"meta.number":                            {},
	"meta.dependency":                        {},
	"meta.instance-type":                     {},
	"meta.false-schema":                      {},
	"meta.ref":                               {},
	"meta.group":                             {},
	"meta.catchall":                          {},
	"meta.truncated":                         {},
}

// ValidationOptions configures validation passes run after indexing.
type ValidationOptions struct {
	// MaxIssues stops collecting after this many issues (0 = unlimited).
	MaxIssues int
	// Context, when non-nil, aborts further issues when Done.
	Context context.Context
	// TargetVersion overrides the effective OpenAPI version used for
	// version-sensitive validation when the document itself does not declare one
	// (for example standalone fragments opened in an editor).
	TargetVersion Version
	// SkipSyntax skips CST / YAML duplicate-key and ERROR-node checks.
	SkipSyntax bool
	// SkipStructural skips root/info/paths/operation structural rules.
	SkipStructural bool
	// SkipSchemaShape skips parameter and security-scheme shape checks.
	SkipSchemaShape bool
	// SkipMetaSchema skips JSON Schema meta-validation against embedded
	// Root/Fragment OpenAPI contracts (Draft 2020-12).
	SkipMetaSchema bool
}

// DefaultValidationOptions returns options used by ParseTree, ParseContent, ParseURI, and NewIndexFromDocument.
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{}
}

// Revalidate recomputes Issues and stores them on the index. Nil-safe.
func (idx *Index) Revalidate(opts ValidationOptions) []Issue {
	if idx == nil {
		return nil
	}
	idx.Issues = idx.computeIssues(opts)
	return idx.Issues
}

func (idx *Index) computeIssues(opts ValidationOptions) []Issue {
	if idx == nil {
		return nil
	}
	var issues []Issue
	sink := &issueSink{opts: opts, issues: &issues, idx: idx}

	if !opts.SkipSyntax {
		if t := idx.Tree(); t != nil {
			if root := t.RootNode(); root != nil {
				validateSyntaxTreeSitter(root, sink)
				validateDuplicateKeysTreeSitter(root, idx.Content(), "", sink)
			}
		} else if idx.Document != nil {
			validateDuplicateKeysYAML(idx, sink)
		}
	}

	if !opts.SkipStructural {
		validateStructural(idx, sink)
	}

	if !opts.SkipSchemaShape {
		validateSchemaShape(idx, sink)
	}

	if !opts.SkipMetaSchema {
		validateMetaSchema(idx, sink)
	}

	return issues
}

type issueSink struct {
	opts   ValidationOptions
	issues *[]Issue
	idx    *Index
}

func (s *issueSink) effectiveVersion() Version {
	if s == nil {
		return VersionUnknown
	}
	if s.opts.TargetVersion != VersionUnknown {
		return s.opts.TargetVersion
	}
	if s.idx != nil {
		return s.idx.Version
	}
	return VersionUnknown
}

func (s *issueSink) canAdd() bool {
	if s.opts.MaxIssues > 0 && len(*s.issues) >= s.opts.MaxIssues {
		return false
	}
	if s.opts.Context != nil {
		select {
		case <-s.opts.Context.Done():
			return false
		default:
		}
	}
	return true
}

func (s *issueSink) add(iss Issue) {
	if !s.canAdd() {
		return
	}
	if iss.DocumentKind == DocumentKindUnknown && s.idx != nil {
		iss.DocumentKind = s.idx.Kind
	}
	if iss.SpecVersion == VersionUnknown {
		iss.SpecVersion = s.effectiveVersion()
	}
	*s.issues = append(*s.issues, iss)
}

func (idx *Index) applyDefaultValidation() {
	if idx == nil {
		return
	}
	idx.Issues = idx.computeIssues(DefaultValidationOptions())
}
