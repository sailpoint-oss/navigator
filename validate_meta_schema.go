package navigator

import (
	"encoding/json"
	"fmt"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

func validateMetaSchema(idx *Index, sink *issueSink) {
	if idx == nil || sink == nil {
		return
	}
	if len(idx.Content()) == 0 {
		// Meta-schema validates decoded source bytes; programmatic indexes
		// (e.g. NewIndexFromDocument) have no serialized document.
		return
	}
	if idx.IsArazzo() {
		validateArazzoMetaSchema(idx, sink)
		return
	}
	doc := idx.Document
	if doc == nil || doc.DocType == DocTypeUnknown || doc.DocType == DocTypeNonOpenAPI {
		return
	}

	rootSch, fragSch, cerr := compiledMetaSchemas()
	if cerr != nil {
		sink.add(Issue{
			Code:     "meta.compiler",
			Message:  fmt.Sprintf("OpenAPI meta-schemas failed to load or compile: %v", cerr),
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
		return
	}
	rootByVersion, fragByVersion, cerr := compiledVersionedMetaSchemas()
	if cerr != nil {
		sink.add(Issue{
			Code:     "meta.compiler",
			Message:  fmt.Sprintf("OpenAPI meta-schemas failed to load or compile: %v", cerr),
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
		return
	}

	inst, err := decodeDocumentInstanceForMeta(idx)
	if err != nil {
		sink.add(Issue{
			Code:     "meta.decode",
			Message:  fmt.Sprintf("Document could not be decoded for meta-schema validation: %v", err),
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
		return
	}

	fragment := doc.DocType == DocTypeFragment
	var sch *jsonschema.Schema
	if fragment {
		fragType := fragmentTypeForIndex(idx)
		if fragType == FragmentUnknown {
			return
		}
		if byType := fragByVersion[sink.effectiveVersion()]; byType != nil {
			sch = byType[fragType]
		}
		if sch == nil {
			sch = fragSch
		}
	} else {
		if version := sink.effectiveVersion(); version != VersionUnknown {
			sch = rootByVersion[version]
		}
		if sch == nil {
			sch = rootSch
		}
	}

	err = sch.Validate(inst)
	if err == nil {
		return
	}

	verr, ok := err.(*jsonschema.ValidationError)
	if !ok {
		sink.add(Issue{
			Code:     "meta.catchall",
			Message:  fmt.Sprintf("Meta-schema validation failed: %v", err),
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
		return
	}

	// Walk the ValidationError tree directly to extract leaf-level errors.
	// The BasicOutput() approach loses detail because OpenAPI JSON Schemas
	// use $ref extensively, producing kind.Reference wrappers that obscure
	// the actual constraint violations.
	leafErrors := collectLeafValidationErrors(verr)

	truncated := false
	for _, le := range leafErrors {
		if !sink.canAdd() {
			truncated = true
			break
		}
		sink.add(issueFromLeafError(le, fragment, idx))
	}
	if truncated && sink.canAdd() {
		sink.add(Issue{
			Code:     "meta.truncated",
			Message:  "Further meta-schema issues were omitted (issue limit reached).",
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
	}
}

func validateArazzoMetaSchema(idx *Index, sink *issueSink) {
	if idx == nil || idx.Arazzo == nil || sink == nil {
		return
	}
	doc := idx.Arazzo
	sch, cerr := compiledArazzoMetaSchema()
	if cerr != nil {
		sink.add(Issue{
			Code:     "meta.compiler",
			Message:  fmt.Sprintf("Arazzo meta-schema failed to load or compile: %v", cerr),
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
		return
	}

	inst, err := decodeDocumentInstanceForMeta(idx)
	if err != nil {
		sink.add(Issue{
			Code:     "meta.decode",
			Message:  fmt.Sprintf("Document could not be decoded for meta-schema validation: %v", err),
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
		return
	}

	err = sch.Validate(inst)
	if err == nil {
		return
	}

	verr, ok := err.(*jsonschema.ValidationError)
	if !ok {
		sink.add(Issue{
			Code:     "meta.catchall",
			Message:  fmt.Sprintf("Meta-schema validation failed: %v", err),
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
		return
	}

	leafErrors := collectLeafValidationErrors(verr)

	truncated := false
	for _, le := range leafErrors {
		if !sink.canAdd() {
			truncated = true
			break
		}
		sink.add(issueFromLeafError(le, false, idx))
	}
	if truncated && sink.canAdd() {
		sink.add(Issue{
			Code:     "meta.truncated",
			Message:  "Further meta-schema issues were omitted (issue limit reached).",
			Pointer:  "",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryMeta,
		})
	}
}

func decodeDocumentInstanceForMeta(idx *Index) (any, error) {
	if idx == nil || len(idx.Content()) == 0 {
		return map[string]any{}, nil
	}
	switch idx.Format {
	case FormatJSON:
		var v any
		if err := json.Unmarshal(idx.Content(), &v); err != nil {
			return nil, err
		}
		return v, nil
	default:
		var v any
		if err := yaml.Unmarshal(idx.Content(), &v); err != nil {
			return nil, err
		}
		return v, nil
	}
}
