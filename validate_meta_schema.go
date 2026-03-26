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
	doc := idx.Document
	if doc == nil || doc.DocType == DocTypeUnknown {
		return
	}
	if doc.DocType != DocTypeRoot && doc.DocType != DocTypeFragment {
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

	out := verr.BasicOutput()
	var units []*jsonschema.OutputUnit
	collectMetaOutputUnits(out, &units)

	truncated := false
	for _, u := range units {
		if !sink.canAdd() {
			truncated = true
			break
		}
		if u == nil || u.Error == nil {
			continue
		}
		sink.add(issueFromMetaOutputUnit(u, fragment, idx))
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
