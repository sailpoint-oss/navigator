package navigator

import "strconv"

// ArazzoDocument is the typed root model for an Arazzo 1.0.x description.
type ArazzoDocument struct {
	Version            string
	Info               *ArazzoInfo
	SourceDescriptions []ArazzoSourceDescription
	Workflows          []ArazzoWorkflow
	Components         *ArazzoComponents
	Extensions         map[string]*Node
	Loc                Loc
	VersionLoc         Loc
}

// ArazzoInfo provides metadata about the Arazzo description.
type ArazzoInfo struct {
	Title       string
	Summary     string
	Description DescriptionValue
	Version     string
	Extensions  map[string]*Node
	Loc         Loc
	TitleLoc    Loc
	VersionLoc  Loc
}

// ArazzoSourceDescription references a source API or workflow document.
type ArazzoSourceDescription struct {
	Name       string
	URL        string
	Type       string
	Extensions map[string]*Node
	Loc        Loc
	NameLoc    Loc
	URLLoc     Loc
}

// ArazzoWorkflow describes a single reusable workflow.
type ArazzoWorkflow struct {
	WorkflowID     string
	Summary        string
	Description    DescriptionValue
	Inputs         *SemanticNode
	DependsOn      []string
	Steps          []ArazzoStep
	SuccessActions []ArazzoAction
	FailureActions []ArazzoAction
	Outputs        map[string]string
	Parameters     []ArazzoParameter
	Extensions     map[string]*Node
	Loc            Loc
	WorkflowIDLoc  Loc
	StepsLoc       Loc
}

// ArazzoStep describes one workflow step.
type ArazzoStep struct {
	StepID          string
	Description     DescriptionValue
	OperationID     string
	OperationPath   string
	WorkflowID      string
	Parameters      []ArazzoParameter
	RequestBody     *ArazzoRequestBody
	SuccessCriteria []ArazzoCriterion
	OnSuccess       []ArazzoAction
	OnFailure       []ArazzoAction
	Outputs         map[string]string
	Extensions      map[string]*Node
	Loc             Loc
	StepIDLoc       Loc
}

// ArazzoParameter describes a direct parameter object or a reusable reference.
type ArazzoParameter struct {
	Name         string
	In           string
	Reference    string
	Value        *SemanticNode
	Extensions   map[string]*Node
	Loc          Loc
	NameLoc      Loc
	ReferenceLoc Loc
}

// ArazzoRequestBody describes the request body passed to an operation step.
type ArazzoRequestBody struct {
	ContentType    string
	Payload        *SemanticNode
	Replacements   []ArazzoPayloadReplacement
	Extensions     map[string]*Node
	Loc            Loc
	ContentTypeLoc Loc
}

// ArazzoPayloadReplacement rewrites one payload location before execution.
type ArazzoPayloadReplacement struct {
	Target     string
	Value      string
	Extensions map[string]*Node
	Loc        Loc
	TargetLoc  Loc
	ValueLoc   Loc
}

// ArazzoCriterion describes a runtime assertion used by Arazzo.
type ArazzoCriterion struct {
	Context    string
	Condition  string
	Type       string
	Version    string
	Extensions map[string]*Node
	Loc        Loc
}

// ArazzoAction describes success/failure actions or reusable action references.
type ArazzoAction struct {
	Name         string
	Type         string
	WorkflowID   string
	StepID       string
	Reference    string
	Value        *SemanticNode
	RetryAfter   *float64
	RetryLimit   *int
	Criteria     []ArazzoCriterion
	Extensions   map[string]*Node
	Loc          Loc
	NameLoc      Loc
	ReferenceLoc Loc
}

// ArazzoComponents holds reusable workflow components.
type ArazzoComponents struct {
	Inputs         map[string]*SemanticNode
	Parameters     map[string]*ArazzoParameter
	SuccessActions map[string]*ArazzoAction
	FailureActions map[string]*ArazzoAction
	Extensions     map[string]*Node
	Loc            Loc
}

// WorkflowIDs returns the workflow IDs in document order.
func (d *ArazzoDocument) WorkflowIDs() []string {
	if d == nil {
		return nil
	}
	ids := make([]string, 0, len(d.Workflows))
	for _, workflow := range d.Workflows {
		if workflow.WorkflowID != "" {
			ids = append(ids, workflow.WorkflowID)
		}
	}
	return ids
}

func parseArazzoDocumentFromSemantic(sn *SemanticNode) *ArazzoDocument {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	if DetectDocumentKind(semanticKeys(sn)) != DocumentKindArazzo {
		return nil
	}
	doc := &ArazzoDocument{
		Version:    snStr(sn, "arazzo"),
		Info:       parseArazzoInfoSN(sn.Get("info")),
		Extensions: snExtensions(sn),
		Loc:        locFromSN(sn),
		VersionLoc: locFromSN(sn.Get("arazzo")),
	}
	if sources := sn.Get("sourceDescriptions"); sources != nil {
		doc.SourceDescriptions = parseArazzoSourceDescriptionsSN(sources)
	}
	if workflows := sn.Get("workflows"); workflows != nil {
		doc.Workflows = parseArazzoWorkflowsSN(workflows)
	}
	if components := sn.Get("components"); components != nil {
		doc.Components = parseArazzoComponentsSN(components)
	}
	return doc
}

func parseArazzoInfoSN(sn *SemanticNode) *ArazzoInfo {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	return &ArazzoInfo{
		Title:       snStr(sn, "title"),
		Summary:     snStr(sn, "summary"),
		Description: snDescription(sn),
		Version:     snStr(sn, "version"),
		Extensions:  snExtensions(sn),
		Loc:         locFromSN(sn),
		TitleLoc:    locFromSN(sn.Get("title")),
		VersionLoc:  locFromSN(sn.Get("version")),
	}
}

func parseArazzoSourceDescriptionsSN(sn *SemanticNode) []ArazzoSourceDescription {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	items := make([]ArazzoSourceDescription, 0, len(sn.Items))
	for _, item := range sn.Items {
		if item == nil || item.Kind != NodeMapping {
			continue
		}
		items = append(items, ArazzoSourceDescription{
			Name:       snStr(item, "name"),
			URL:        snStr(item, "url"),
			Type:       snStr(item, "type"),
			Extensions: snExtensions(item),
			Loc:        locFromSN(item),
			NameLoc:    locFromSN(item.Get("name")),
			URLLoc:     locFromSN(item.Get("url")),
		})
	}
	return items
}

func parseArazzoWorkflowsSN(sn *SemanticNode) []ArazzoWorkflow {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	workflows := make([]ArazzoWorkflow, 0, len(sn.Items))
	for _, item := range sn.Items {
		if workflow := parseArazzoWorkflowSN(item); workflow != nil {
			workflows = append(workflows, *workflow)
		}
	}
	return workflows
}

func parseArazzoWorkflowSN(sn *SemanticNode) *ArazzoWorkflow {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	workflow := &ArazzoWorkflow{
		WorkflowID:     snStr(sn, "workflowId"),
		Summary:        snStr(sn, "summary"),
		Description:    snDescription(sn),
		Inputs:         sn.Get("inputs"),
		DependsOn:      snStringSlice(sn.Get("dependsOn")),
		SuccessActions: parseArazzoActionsSN(sn.Get("successActions")),
		FailureActions: parseArazzoActionsSN(sn.Get("failureActions")),
		Outputs:        snStringMap(sn.Get("outputs")),
		Parameters:     parseArazzoParametersSN(sn.Get("parameters")),
		Extensions:     snExtensions(sn),
		Loc:            locFromSN(sn),
		WorkflowIDLoc:  locFromSN(sn.Get("workflowId")),
		StepsLoc:       locFromSN(sn.Get("steps")),
	}
	if steps := sn.Get("steps"); steps != nil {
		workflow.Steps = parseArazzoStepsSN(steps)
	}
	return workflow
}

func parseArazzoStepsSN(sn *SemanticNode) []ArazzoStep {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	steps := make([]ArazzoStep, 0, len(sn.Items))
	for _, item := range sn.Items {
		if step := parseArazzoStepSN(item); step != nil {
			steps = append(steps, *step)
		}
	}
	return steps
}

func parseArazzoStepSN(sn *SemanticNode) *ArazzoStep {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	step := &ArazzoStep{
		StepID:          snStr(sn, "stepId"),
		Description:     snDescription(sn),
		OperationID:     snStr(sn, "operationId"),
		OperationPath:   snStr(sn, "operationPath"),
		WorkflowID:      snStr(sn, "workflowId"),
		Parameters:      parseArazzoParametersSN(sn.Get("parameters")),
		SuccessCriteria: parseArazzoCriteriaSN(sn.Get("successCriteria")),
		OnSuccess:       parseArazzoActionsSN(sn.Get("onSuccess")),
		OnFailure:       parseArazzoActionsSN(sn.Get("onFailure")),
		Outputs:         snStringMap(sn.Get("outputs")),
		Extensions:      snExtensions(sn),
		Loc:             locFromSN(sn),
		StepIDLoc:       locFromSN(sn.Get("stepId")),
	}
	if requestBody := sn.Get("requestBody"); requestBody != nil {
		step.RequestBody = parseArazzoRequestBodySN(requestBody)
	}
	return step
}

func parseArazzoParametersSN(sn *SemanticNode) []ArazzoParameter {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	params := make([]ArazzoParameter, 0, len(sn.Items))
	for _, item := range sn.Items {
		if param := parseArazzoParameterSN(item); param != nil {
			params = append(params, *param)
		}
	}
	return params
}

func parseArazzoParameterSN(sn *SemanticNode) *ArazzoParameter {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	return &ArazzoParameter{
		Name:         snStr(sn, "name"),
		In:           snStr(sn, "in"),
		Reference:    snStr(sn, "reference"),
		Value:        sn.Get("value"),
		Extensions:   snExtensions(sn),
		Loc:          locFromSN(sn),
		NameLoc:      locFromSN(sn.Get("name")),
		ReferenceLoc: locFromSN(sn.Get("reference")),
	}
}

func parseArazzoRequestBodySN(sn *SemanticNode) *ArazzoRequestBody {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	body := &ArazzoRequestBody{
		ContentType:    snStr(sn, "contentType"),
		Payload:        sn.Get("payload"),
		Extensions:     snExtensions(sn),
		Loc:            locFromSN(sn),
		ContentTypeLoc: locFromSN(sn.Get("contentType")),
	}
	if replacements := sn.Get("replacements"); replacements != nil && replacements.Kind == NodeSequence {
		body.Replacements = make([]ArazzoPayloadReplacement, 0, len(replacements.Items))
		for _, item := range replacements.Items {
			if item == nil || item.Kind != NodeMapping {
				continue
			}
			body.Replacements = append(body.Replacements, ArazzoPayloadReplacement{
				Target:     snStr(item, "target"),
				Value:      snStr(item, "value"),
				Extensions: snExtensions(item),
				Loc:        locFromSN(item),
				TargetLoc:  locFromSN(item.Get("target")),
				ValueLoc:   locFromSN(item.Get("value")),
			})
		}
	}
	return body
}

func parseArazzoCriteriaSN(sn *SemanticNode) []ArazzoCriterion {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	criteria := make([]ArazzoCriterion, 0, len(sn.Items))
	for _, item := range sn.Items {
		if item == nil || item.Kind != NodeMapping {
			continue
		}
		criteria = append(criteria, ArazzoCriterion{
			Context:    snStr(item, "context"),
			Condition:  snStr(item, "condition"),
			Type:       snStr(item, "type"),
			Version:    snStr(item, "version"),
			Extensions: snExtensions(item),
			Loc:        locFromSN(item),
		})
	}
	return criteria
}

func parseArazzoActionsSN(sn *SemanticNode) []ArazzoAction {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	actions := make([]ArazzoAction, 0, len(sn.Items))
	for _, item := range sn.Items {
		if action := parseArazzoActionSN(item); action != nil {
			actions = append(actions, *action)
		}
	}
	return actions
}

func parseArazzoActionSN(sn *SemanticNode) *ArazzoAction {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	return &ArazzoAction{
		Name:         snStr(sn, "name"),
		Type:         snStr(sn, "type"),
		WorkflowID:   snStr(sn, "workflowId"),
		StepID:       snStr(sn, "stepId"),
		Reference:    snStr(sn, "reference"),
		Value:        sn.Get("value"),
		RetryAfter:   snFloatValue(sn.Get("retryAfter")),
		RetryLimit:   snIntValue(sn.Get("retryLimit")),
		Criteria:     parseArazzoCriteriaSN(sn.Get("criteria")),
		Extensions:   snExtensions(sn),
		Loc:          locFromSN(sn),
		NameLoc:      locFromSN(sn.Get("name")),
		ReferenceLoc: locFromSN(sn.Get("reference")),
	}
}

func parseArazzoComponentsSN(sn *SemanticNode) *ArazzoComponents {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	components := &ArazzoComponents{
		Inputs:         parseSemanticMapSN(sn.Get("inputs")),
		Parameters:     parseArazzoParameterMapSN(sn.Get("parameters")),
		SuccessActions: parseArazzoActionMapSN(sn.Get("successActions")),
		FailureActions: parseArazzoActionMapSN(sn.Get("failureActions")),
		Extensions:     snExtensions(sn),
		Loc:            locFromSN(sn),
	}
	return components
}

func parseArazzoParameterMapSN(sn *SemanticNode) map[string]*ArazzoParameter {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	items := make(map[string]*ArazzoParameter, len(sn.Children))
	for key, child := range sn.Children {
		items[key] = parseArazzoParameterSN(child)
	}
	return items
}

func parseArazzoActionMapSN(sn *SemanticNode) map[string]*ArazzoAction {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	items := make(map[string]*ArazzoAction, len(sn.Children))
	for key, child := range sn.Children {
		items[key] = parseArazzoActionSN(child)
	}
	return items
}

func parseSemanticMapSN(sn *SemanticNode) map[string]*SemanticNode {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	items := make(map[string]*SemanticNode, len(sn.Children))
	for key, child := range sn.Children {
		items[key] = child
	}
	return items
}

func snStringMap(sn *SemanticNode) map[string]string {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	items := make(map[string]string, len(sn.Children))
	for key, child := range sn.Children {
		items[key] = child.StringValue()
	}
	return items
}

func snStringSlice(sn *SemanticNode) []string {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	items := make([]string, 0, len(sn.Items))
	for _, item := range sn.Items {
		items = append(items, item.StringValue())
	}
	return items
}

func snExtensions(sn *SemanticNode) map[string]*Node {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	var extensions map[string]*Node
	for key, child := range sn.Children {
		if len(key) >= 2 && key[0] == 'x' && key[1] == '-' {
			if extensions == nil {
				extensions = make(map[string]*Node)
			}
			extensions[key] = snToNode(child)
		}
	}
	return extensions
}

func snFloatValue(sn *SemanticNode) *float64 {
	if sn == nil {
		return nil
	}
	value, err := strconv.ParseFloat(sn.StringValue(), 64)
	if err != nil {
		return nil
	}
	return &value
}

func snIntValue(sn *SemanticNode) *int {
	if sn == nil {
		return nil
	}
	value, err := strconv.Atoi(sn.StringValue())
	if err != nil {
		return nil
	}
	return &value
}
