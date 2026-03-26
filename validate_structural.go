package navigator

import (
	"fmt"
	"strings"
)

func validateStructural(idx *Index, sink *issueSink) {
	if idx == nil {
		return
	}
	doc := idx.Document
	root := idx.SemanticRoot()

	if doc == nil {
		if root != nil {
			sink.add(Issue{
				Code:     "structural.root-not-mapping",
				Message:  "document root must be a YAML/JSON object",
				Pointer:  "",
				Range:    root.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		return
	}

	if doc.DocType == DocTypeFragment {
		validateFragmentStructural(idx, sink)
		return
	}

	if doc.DocType != DocTypeRoot {
		return
	}

	if root != nil {
		if n := root.Get("openapi"); n != nil && n.Kind != NodeScalar {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"openapi" must be a string`,
				Pointer:  "/openapi",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("swagger"); n != nil && n.Kind != NodeScalar {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"swagger" must be a string`,
				Pointer:  "/swagger",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("info"); n != nil && n.Kind != NodeMapping {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"info" must be an object`,
				Pointer:  "/info",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("paths"); n != nil && n.Kind != NodeMapping {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"paths" must be an object`,
				Pointer:  "/paths",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("components"); n != nil && n.Kind != NodeMapping {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"components" must be an object`,
				Pointer:  "/components",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("servers"); n != nil && n.Kind != NodeSequence {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"servers" must be an array`,
				Pointer:  "/servers",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("tags"); n != nil && n.Kind != NodeSequence {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"tags" must be an array`,
				Pointer:  "/tags",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("security"); n != nil && n.Kind != NodeSequence {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"security" must be an array`,
				Pointer:  "/security",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
	}

	if doc.ParsedVersion == VersionUnknown && doc.Version != "" {
		sink.add(Issue{
			Code:     "structural.unsupported-version",
			Message:  fmt.Sprintf("unsupported OpenAPI version %q", doc.Version),
			Pointer:  versionPointer(root, doc),
			Range:    versionRange(doc, root),
			Severity: SeverityWarning,
			Category: CategoryStructural,
		})
	}

	if isOpenAPI3(doc.ParsedVersion) || looksLikeOAS3(root, doc) {
		if doc.Info == nil {
			sink.add(Issue{
				Code:     "structural.missing-info",
				Message:  "OpenAPI 3.x documents require an `info` object",
				Pointer:  "/info",
				Range:    doc.Loc.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		} else {
			if doc.Info.Title == "" {
				sink.add(Issue{
					Code:     "structural.missing-info-title",
					Message:  "`info.title` is required",
					Pointer:  "/info/title",
					Range:    firstNonEmptyLoc(doc.Info.TitleLoc, doc.Info.Loc).Range,
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			}
			if doc.Info.Version == "" {
				sink.add(Issue{
					Code:     "structural.missing-info-version",
					Message:  "`info.version` is required",
					Pointer:  "/info/version",
					Range:    firstNonEmptyLoc(doc.Info.VersionLoc, doc.Info.Loc).Range,
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			}
		}
	}

	if doc.ParsedVersion == Version20 {
		if doc.Info == nil {
			sink.add(Issue{
				Code:     "structural.missing-info",
				Message:  "Swagger 2.0 documents require an `info` object",
				Pointer:  "/info",
				Range:    doc.Loc.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		} else {
			if doc.Info.Title == "" {
				sink.add(Issue{
					Code:     "structural.missing-info-title",
					Message:  "`info.title` is required",
					Pointer:  "/info/title",
					Range:    firstNonEmptyLoc(doc.Info.TitleLoc, doc.Info.Loc).Range,
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			}
			if doc.Info.Version == "" {
				sink.add(Issue{
					Code:     "structural.missing-info-version",
					Message:  "`info.version` is required",
					Pointer:  "/info/version",
					Range:    firstNonEmptyLoc(doc.Info.VersionLoc, doc.Info.Loc).Range,
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			}
		}
	}

	if doc.Paths == nil {
		return
	}
	for pth, item := range doc.Paths {
		if item == nil {
			continue
		}
		if item.Ref != "" && !pathItemHasOperations(item) {
			continue
		}
		for _, mo := range item.Operations() {
			op := mo.Operation
			if op == nil {
				continue
			}
			if len(op.Responses) == 0 {
				rng := op.ResponsesLoc.Range
				if IsEmpty(rng) {
					rng = op.Loc.Range
				}
				sink.add(Issue{
					Code:     "structural.missing-responses",
					Message:  fmt.Sprintf("operation %s %s must declare non-empty `responses`", strings.ToUpper(mo.Method), pth),
					Pointer:  "/paths/" + EscapeJSONPointer(pth) + "/" + mo.Method + "/responses",
					Range:    rng,
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			}
		}
	}
}

func validateFragmentStructural(idx *Index, sink *issueSink) {
	if idx == nil || sink == nil || idx.Document == nil {
		return
	}
	version := sink.effectiveVersion()
	if version == VersionUnknown {
		return
	}
	if version != Version30 {
		return
	}

	root := idx.SemanticRoot()
	if root == nil || root.Kind != NodeMapping {
		return
	}
	if fragmentTypeForIndex(idx) != FragmentOperation {
		return
	}
	if root.Get("responses") != nil {
		return
	}

	sink.add(Issue{
		Code:     "structural.missing-responses",
		Message:  "OpenAPI 3.0 operation fragments require a `responses` object",
		Pointer:  "/responses",
		Range:    idx.Document.Loc.Range,
		Severity: SeverityError,
		Category: CategoryStructural,
	})
}

func fragmentTypeForIndex(idx *Index) FragmentType {
	if idx == nil {
		return FragmentUnknown
	}
	root := idx.SemanticRoot()
	if root == nil || root.Kind != NodeMapping {
		return FragmentUnknown
	}
	keys := make([]string, 0, len(root.Children))
	for key := range root.Children {
		keys = append(keys, key)
	}
	return DetectFragmentType(keys)
}

func isOpenAPI3(v Version) bool {
	switch v {
	case Version30, Version31, Version32:
		return true
	default:
		return false
	}
}

func looksLikeOAS3(root *SemanticNode, doc *Document) bool {
	if doc == nil || doc.DocType != DocTypeRoot {
		return false
	}
	if isOpenAPI3(doc.ParsedVersion) {
		return false
	}
	if doc.ParsedVersion != VersionUnknown {
		return false
	}
	if strings.HasPrefix(doc.Version, "3.") {
		return true
	}
	if root == nil {
		return false
	}
	if o := root.Get("openapi"); o != nil && o.Kind == NodeScalar {
		return strings.HasPrefix(o.StringValue(), "3.")
	}
	return false
}

func versionPointer(root *SemanticNode, doc *Document) string {
	if root != nil && root.Get("openapi") != nil {
		return "/openapi"
	}
	if root != nil && root.Get("swagger") != nil {
		return "/swagger"
	}
	if doc != nil {
		return "/openapi"
	}
	return ""
}

func versionRange(doc *Document, root *SemanticNode) Range {
	if root != nil {
		if o := root.Get("openapi"); o != nil {
			return o.Range
		}
		if s := root.Get("swagger"); s != nil {
			return s.Range
		}
	}
	if doc != nil {
		return doc.Loc.Range
	}
	return FileStartRange
}

func firstNonEmptyLoc(a, b Loc) Loc {
	if !IsEmpty(a.Range) {
		return a
	}
	return b
}

func pathItemHasOperations(p *PathItem) bool {
	if p == nil {
		return false
	}
	return p.Get != nil || p.Put != nil || p.Post != nil || p.Delete != nil ||
		p.Options != nil || p.Head != nil || p.Patch != nil || p.Trace != nil
}
