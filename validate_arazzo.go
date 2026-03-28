package navigator

import (
	"fmt"
	"strconv"
	"strings"
)

func validateArazzoStructural(idx *Index, sink *issueSink) {
	if idx == nil || idx.Arazzo == nil {
		return
	}
	doc := idx.Arazzo
	root := idx.SemanticRoot()

	if root != nil {
		if n := root.Get("arazzo"); n != nil && n.Kind != NodeScalar {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"arazzo" must be a string`,
				Pointer:  "/arazzo",
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
		if n := root.Get("sourceDescriptions"); n != nil && n.Kind != NodeSequence {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"sourceDescriptions" must be an array`,
				Pointer:  "/sourceDescriptions",
				Range:    n.Range,
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		}
		if n := root.Get("workflows"); n != nil && n.Kind != NodeSequence {
			sink.add(Issue{
				Code:     "structural.ir-wrong-type",
				Message:  `"workflows" must be an array`,
				Pointer:  "/workflows",
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
	}

	if doc.Version != "" && !isSupportedArazzoVersion(doc.Version) {
		sink.add(Issue{
			Code:     "structural.unsupported-arazzo-version",
			Message:  fmt.Sprintf("unsupported Arazzo version %q", doc.Version),
			Pointer:  "/arazzo",
			Range:    doc.VersionLoc.Range,
			Severity: SeverityWarning,
			Category: CategoryStructural,
		})
	}

	if doc.Info == nil {
		sink.add(Issue{
			Code:     "structural.missing-info",
			Message:  "Arazzo documents require an `info` object",
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

	if len(doc.SourceDescriptions) == 0 {
		sink.add(Issue{
			Code:     "structural.missing-source-descriptions",
			Message:  "Arazzo documents require at least one `sourceDescriptions` entry",
			Pointer:  "/sourceDescriptions",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryStructural,
		})
	}

	if len(doc.Workflows) == 0 {
		sink.add(Issue{
			Code:     "structural.missing-workflows",
			Message:  "Arazzo documents require at least one workflow",
			Pointer:  "/workflows",
			Range:    doc.Loc.Range,
			Severity: SeverityError,
			Category: CategoryStructural,
		})
		return
	}

	seenWorkflows := make(map[string]int, len(doc.Workflows))
	for workflowIndex, workflow := range doc.Workflows {
		ptr := "/workflows/" + strconv.Itoa(workflowIndex)
		if workflow.WorkflowID == "" {
			sink.add(Issue{
				Code:     "structural.missing-workflow-id",
				Message:  "workflow must include `workflowId`",
				Pointer:  ptr + "/workflowId",
				Range:    firstNonEmptyRange(workflow.WorkflowIDLoc.Range, workflow.Loc.Range),
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		} else if other, ok := seenWorkflows[workflow.WorkflowID]; ok {
			sink.add(Issue{
				Code:     "structural.duplicate-workflow-id",
				Message:  fmt.Sprintf("workflowId %q is duplicated; first declared at /workflows/%d", workflow.WorkflowID, other),
				Pointer:  ptr + "/workflowId",
				Range:    firstNonEmptyRange(workflow.WorkflowIDLoc.Range, workflow.Loc.Range),
				Severity: SeverityError,
				Category: CategoryStructural,
			})
		} else {
			seenWorkflows[workflow.WorkflowID] = workflowIndex
		}

		seenSteps := make(map[string]int, len(workflow.Steps))
		for stepIndex, step := range workflow.Steps {
			stepPtr := ptr + "/steps/" + strconv.Itoa(stepIndex)
			if step.StepID == "" {
				sink.add(Issue{
					Code:     "structural.missing-step-id",
					Message:  "step must include `stepId`",
					Pointer:  stepPtr + "/stepId",
					Range:    firstNonEmptyRange(step.StepIDLoc.Range, step.Loc.Range),
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			} else if other, ok := seenSteps[step.StepID]; ok {
				sink.add(Issue{
					Code:     "structural.duplicate-step-id",
					Message:  fmt.Sprintf("stepId %q is duplicated within workflow %q; first declared at %s", step.StepID, workflow.WorkflowID, ptr+"/steps/"+strconv.Itoa(other)),
					Pointer:  stepPtr + "/stepId",
					Range:    firstNonEmptyRange(step.StepIDLoc.Range, step.Loc.Range),
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			} else {
				seenSteps[step.StepID] = stepIndex
			}

			targetCount := 0
			for _, value := range []string{step.OperationID, step.OperationPath, step.WorkflowID} {
				if strings.TrimSpace(value) != "" {
					targetCount++
				}
			}
			if targetCount != 1 {
				sink.add(Issue{
					Code:     "structural.step-target-count",
					Message:  "step must declare exactly one of `operationId`, `operationPath`, or `workflowId`",
					Pointer:  stepPtr,
					Range:    step.Loc.Range,
					Severity: SeverityError,
					Category: CategoryStructural,
				})
			}
		}
	}
}

func isSupportedArazzoVersion(version string) bool {
	return strings.HasPrefix(version, "1.0.")
}

func firstNonEmptyRange(primary, fallback Range) Range {
	if !IsEmpty(primary) {
		return primary
	}
	return fallback
}
