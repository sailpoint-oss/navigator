package navigator

import "fmt"

func (idx *Index) resolveSemanticRefParts(parts []string) (interface{}, error) {
	if idx == nil || idx.semanticRoot == nil {
		return nil, fmt.Errorf("resolve: semantic root unavailable")
	}
	target, err := semanticNodeAt(idx.semanticRoot, parts)
	if err != nil {
		return nil, err
	}
	return semanticValueForResolvedNode(parts, target, idx.Document), nil
}

func semanticNodeAt(root *SemanticNode, parts []string) (*SemanticNode, error) {
	current := root
	for _, part := range parts {
		if current == nil {
			return nil, fmt.Errorf("resolve: pointer segment %q not found", part)
		}
		switch current.Kind {
		case NodeMapping:
			next := current.Get(part)
			if next == nil {
				return nil, fmt.Errorf("resolve: pointer segment %q not found", part)
			}
			current = next
		case NodeSequence:
			i, err := parsePointerIndex(part)
			if err != nil {
				return nil, fmt.Errorf("resolve: invalid array index %q", part)
			}
			next := current.Index(i)
			if next == nil {
				return nil, fmt.Errorf("resolve: index %d out of range", i)
			}
			current = next
		default:
			return nil, fmt.Errorf("resolve: cannot traverse into scalar at %q", part)
		}
	}
	return current, nil
}

func semanticValueForResolvedNode(parts []string, target *SemanticNode, doc *Document) interface{} {
	if len(parts) == 0 {
		if doc != nil {
			return doc
		}
		return semanticNodeOpenAPIValue(target)
	}
	if target == nil {
		return nil
	}

	if target.Kind == NodeMapping {
		switch parts[0] {
		case "info":
			return (&treeParser{}).parseInfoSN(target)
		case "components":
			if len(parts) == 1 {
				if doc != nil && doc.Components != nil {
					return doc.Components
				}
				return parseComponentsFromSemantic(target)
			}
		case "tags":
			return parseTagFromSemantic(target)
		case "servers":
			return parseServerFromSemantic(target)
		}
	}

	if value := semanticNodeOpenAPIValueForParts(parts, target); value != nil {
		return value
	}

	switch target.Kind {
	case NodeScalar:
		return target.Value
	case NodeNull:
		return nil
	default:
		return target
	}
}

func semanticNodeOpenAPIValue(sn *SemanticNode) interface{} {
	return semanticNodeOpenAPIValueForParts(nil, sn)
}

func semanticNodeOpenAPIValueForParts(parts []string, sn *SemanticNode) interface{} {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	tp := &treeParser{}

	switch {
	case looksLikeInfoSemanticNode(sn):
		return tp.parseInfoSN(sn)
	case looksLikeTagSemanticNode(sn):
		return parseTagFromSemantic(sn)
	case looksLikeComponentsSemanticNode(sn):
		return parseComponentsFromSemantic(sn)
	}

	switch semanticFragmentType(sn) {
	case FragmentSchema:
		return tp.parseSchemaSN(sn)
	case FragmentPathItem:
		return parsePathItemFromSemantic(sn)
	case FragmentOperation:
		return parseOperationFromSemantic(sn, operationMethodFromParts(parts))
	case FragmentParameter:
		return tp.parseParameterSN(sn)
	case FragmentRequestBody:
		return tp.parseRequestBodySN(sn)
	case FragmentResponse:
		return parseResponseFromSemantic(sn)
	case FragmentHeader:
		return parseHeaderFromSemantic(sn)
	case FragmentSecurityScheme:
		return parseSecuritySchemeFromSemantic(sn)
	case FragmentComponents:
		return parseComponentsFromSemantic(sn)
	case FragmentServer:
		return parseServerFromSemantic(sn)
	default:
		return nil
	}
}

func semanticFragmentType(sn *SemanticNode) FragmentType {
	keys := semanticKeys(sn)
	if len(keys) == 0 {
		return FragmentUnknown
	}
	keySet := make(map[string]bool, len(keys))
	for _, key := range keys {
		keySet[key] = true
	}
	if keySet["openapi"] || keySet["swagger"] {
		return FragmentUnknown
	}
	if keySet["components"] && !keySet["type"] && !keySet["properties"] {
		return FragmentComponents
	}
	if keySet["type"] && (keySet["scheme"] || keySet["flows"] || keySet["openIdConnectUrl"]) {
		return FragmentSecurityScheme
	}
	if keySet["in"] && keySet["name"] {
		return FragmentParameter
	}
	if keySet["schema"] && !keySet["in"] && !keySet["name"] {
		return FragmentHeader
	}
	return DetectFragmentType(keys)
}

func semanticKeys(sn *SemanticNode) []string {
	if sn == nil || sn.Kind != NodeMapping || len(sn.Children) == 0 {
		return nil
	}
	keys := make([]string, 0, len(sn.Children))
	for key := range sn.Children {
		keys = append(keys, key)
	}
	return keys
}

func looksLikeInfoSemanticNode(sn *SemanticNode) bool {
	return sn != nil &&
		sn.Kind == NodeMapping &&
		sn.Get("title") != nil &&
		sn.Get("version") != nil &&
		sn.Get("type") == nil &&
		sn.Get("properties") == nil
}

func looksLikeTagSemanticNode(sn *SemanticNode) bool {
	return sn != nil &&
		sn.Kind == NodeMapping &&
		sn.Get("name") != nil &&
		sn.Get("in") == nil &&
		sn.Get("schema") == nil &&
		sn.Get("url") == nil &&
		sn.Get("operationId") == nil &&
		sn.Get("responses") == nil
}

func looksLikeComponentsSemanticNode(sn *SemanticNode) bool {
	if sn == nil || sn.Kind != NodeMapping {
		return false
	}
	if sn.Get("openapi") != nil || sn.Get("swagger") != nil {
		return false
	}
	componentKeys := []string{
		"schemas",
		"responses",
		"parameters",
		"examples",
		"requestBodies",
		"headers",
		"securitySchemes",
		"links",
		"pathItems",
		"callbacks",
	}
	for _, key := range componentKeys {
		if sn.Get(key) != nil {
			return true
		}
	}
	return false
}

func parsePathItemFromSemantic(sn *SemanticNode) *PathItem {
	item := (&treeParser{}).parsePathItemSN(sn)
	if item == nil {
		return nil
	}
	if servers := sn.Get("servers"); servers != nil {
		item.Servers = parseServersFromSemantic(servers)
	}
	return item
}

func parseOperationFromSemantic(sn *SemanticNode, method string) *Operation {
	op := (&treeParser{}).parseOperationSN(sn, method)
	if op == nil {
		return nil
	}
	if servers := sn.Get("servers"); servers != nil {
		op.Servers = parseServersFromSemantic(servers)
	}
	if externalDocs := sn.Get("externalDocs"); externalDocs != nil {
		op.ExternalDocs = (&treeParser{}).parseExternalDocsSN(externalDocs)
	}
	return op
}

func parseResponseFromSemantic(sn *SemanticNode) *Response {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	tp := &treeParser{}
	r := &Response{
		Description: snDescription(sn),
		Ref:         snStr(sn, "$ref"),
		Loc:         locFromSN(sn),
	}
	r.Content = tp.parseContentMediaTypesSN(sn.Get("content"))
	if headers := sn.Get("headers"); headers != nil && headers.Kind == NodeMapping {
		r.Headers = make(map[string]*Header)
		for name, child := range headers.Children {
			if h := parseHeaderFromSemantic(child); h != nil {
				r.Headers[name] = h
			}
		}
	}
	if links := sn.Get("links"); links != nil && links.Kind == NodeMapping {
		r.Links = make(map[string]*Link)
		for name, child := range links.Children {
			if link := parseLinkFromSemantic(child); link != nil {
				r.Links[name] = link
			}
		}
	}
	return r
}

func parseHeaderFromSemantic(sn *SemanticNode) *Header {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	h := &Header{
		Description: snDescription(sn),
		Required:    semanticBool(sn, "required"),
		Deprecated:  semanticBool(sn, "deprecated"),
		Ref:         snStr(sn, "$ref"),
		Loc:         locFromSN(sn),
	}
	if schema := sn.Get("schema"); schema != nil {
		h.Schema = (&treeParser{}).parseSchemaSN(schema)
	}
	return h
}

func parseLinkFromSemantic(sn *SemanticNode) *Link {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	return &Link{
		OperationRef: snStr(sn, "operationRef"),
		OperationID:  snStr(sn, "operationId"),
		Description:  snDescription(sn),
		Ref:          snStr(sn, "$ref"),
		Loc:          locFromSN(sn),
	}
}

func parseExampleFromSemantic(sn *SemanticNode) *Example {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	example := &Example{
		Summary:       snStr(sn, "summary"),
		Description:   snDescription(sn),
		ExternalValue: snStr(sn, "externalValue"),
		Ref:           snStr(sn, "$ref"),
		Loc:           locFromSN(sn),
	}
	if value := sn.Get("value"); value != nil {
		example.Value = snToNode(value)
	}
	return example
}

func parseTagFromSemantic(sn *SemanticNode) *Tag {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	tag := &Tag{
		Name:        snStr(sn, "name"),
		Description: snDescription(sn),
		Loc:         locFromSN(sn),
	}
	if externalDocs := sn.Get("externalDocs"); externalDocs != nil {
		tag.ExternalDocs = (&treeParser{}).parseExternalDocsSN(externalDocs)
	}
	return tag
}

func parseServersFromSemantic(sn *SemanticNode) []Server {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	servers := make([]Server, 0, len(sn.Items))
	for _, item := range sn.Items {
		if server := parseServerFromSemantic(item); server != nil {
			servers = append(servers, *server)
		}
	}
	return servers
}

func parseServerFromSemantic(sn *SemanticNode) *Server {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	server := &Server{
		URL:         snStr(sn, "url"),
		Description: snDescription(sn),
		Loc:         locFromSN(sn),
	}
	if variables := sn.Get("variables"); variables != nil && variables.Kind == NodeMapping {
		server.Variables = make(map[string]*ServerVariable)
		for name, child := range variables.Children {
			if variable := parseServerVariableFromSemantic(child); variable != nil {
				server.Variables[name] = variable
			}
		}
	}
	return server
}

func parseServerVariableFromSemantic(sn *SemanticNode) *ServerVariable {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	variable := &ServerVariable{
		Default:     snStr(sn, "default"),
		Description: snDescription(sn),
		Loc:         locFromSN(sn),
	}
	if enum := sn.Get("enum"); enum != nil && enum.Kind == NodeSequence {
		for _, item := range enum.Items {
			variable.Enum = append(variable.Enum, item.StringValue())
		}
	}
	return variable
}

func parseSecuritySchemeFromSemantic(sn *SemanticNode) *SecurityScheme {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	scheme := &SecurityScheme{
		Type:             snStr(sn, "type"),
		Description:      snDescription(sn),
		Name:             snStr(sn, "name"),
		In:               snStr(sn, "in"),
		Scheme:           snStr(sn, "scheme"),
		BearerFormat:     snStr(sn, "bearerFormat"),
		OpenIDConnectURL: snStr(sn, "openIdConnectUrl"),
		Ref:              snStr(sn, "$ref"),
		Loc:              locFromSN(sn),
	}
	if flows := sn.Get("flows"); flows != nil {
		scheme.Flows = parseOAuthFlowsFromSemantic(flows)
	}
	return scheme
}

func parseOAuthFlowsFromSemantic(sn *SemanticNode) *OAuthFlows {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	flows := &OAuthFlows{Loc: locFromSN(sn)}
	if implicit := sn.Get("implicit"); implicit != nil {
		flows.Implicit = parseOAuthFlowFromSemantic(implicit)
	}
	if password := sn.Get("password"); password != nil {
		flows.Password = parseOAuthFlowFromSemantic(password)
	}
	if clientCredentials := sn.Get("clientCredentials"); clientCredentials != nil {
		flows.ClientCredentials = parseOAuthFlowFromSemantic(clientCredentials)
	}
	if authorizationCode := sn.Get("authorizationCode"); authorizationCode != nil {
		flows.AuthorizationCode = parseOAuthFlowFromSemantic(authorizationCode)
	}
	return flows
}

func parseOAuthFlowFromSemantic(sn *SemanticNode) *OAuthFlow {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	flow := &OAuthFlow{
		AuthorizationURL: snStr(sn, "authorizationUrl"),
		TokenURL:         snStr(sn, "tokenUrl"),
		RefreshURL:       snStr(sn, "refreshUrl"),
		Loc:              locFromSN(sn),
	}
	if scopes := sn.Get("scopes"); scopes != nil && scopes.Kind == NodeMapping {
		flow.Scopes = make(map[string]string)
		for name, child := range scopes.Children {
			flow.Scopes[name] = child.StringValue()
		}
	}
	return flow
}

func parseComponentsFromSemantic(sn *SemanticNode) *Components {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	tp := &treeParser{}
	components := &Components{Loc: locFromSN(sn)}

	if schemas := sn.Get("schemas"); schemas != nil && schemas.Kind == NodeMapping {
		components.Schemas = make(map[string]*Schema)
		for name, child := range schemas.Children {
			if schema := tp.parseSchemaSN(child); schema != nil {
				components.Schemas[name] = schema
			}
		}
	}
	if responses := sn.Get("responses"); responses != nil && responses.Kind == NodeMapping {
		components.Responses = make(map[string]*Response)
		for name, child := range responses.Children {
			if response := parseResponseFromSemantic(child); response != nil {
				components.Responses[name] = response
			}
		}
	}
	if parameters := sn.Get("parameters"); parameters != nil && parameters.Kind == NodeMapping {
		components.Parameters = make(map[string]*Parameter)
		for name, child := range parameters.Children {
			if parameter := tp.parseParameterSN(child); parameter != nil {
				components.Parameters[name] = parameter
			}
		}
	}
	if examples := sn.Get("examples"); examples != nil && examples.Kind == NodeMapping {
		components.Examples = make(map[string]*Example)
		for name, child := range examples.Children {
			if example := parseExampleFromSemantic(child); example != nil {
				components.Examples[name] = example
			}
		}
	}
	if requestBodies := sn.Get("requestBodies"); requestBodies != nil && requestBodies.Kind == NodeMapping {
		components.RequestBodies = make(map[string]*RequestBody)
		for name, child := range requestBodies.Children {
			if requestBody := tp.parseRequestBodySN(child); requestBody != nil {
				components.RequestBodies[name] = requestBody
			}
		}
	}
	if headers := sn.Get("headers"); headers != nil && headers.Kind == NodeMapping {
		components.Headers = make(map[string]*Header)
		for name, child := range headers.Children {
			if header := parseHeaderFromSemantic(child); header != nil {
				components.Headers[name] = header
			}
		}
	}
	if securitySchemes := sn.Get("securitySchemes"); securitySchemes != nil && securitySchemes.Kind == NodeMapping {
		components.SecuritySchemes = make(map[string]*SecurityScheme)
		for name, child := range securitySchemes.Children {
			if scheme := parseSecuritySchemeFromSemantic(child); scheme != nil {
				components.SecuritySchemes[name] = scheme
			}
		}
	}
	if links := sn.Get("links"); links != nil && links.Kind == NodeMapping {
		components.Links = make(map[string]*Link)
		for name, child := range links.Children {
			if link := parseLinkFromSemantic(child); link != nil {
				components.Links[name] = link
			}
		}
	}
	if pathItems := sn.Get("pathItems"); pathItems != nil && pathItems.Kind == NodeMapping {
		components.PathItems = make(map[string]*PathItem)
		for name, child := range pathItems.Children {
			if item := parsePathItemFromSemantic(child); item != nil {
				components.PathItems[name] = item
			}
		}
	}

	return components
}

func operationMethodFromParts(parts []string) string {
	for i := len(parts) - 1; i >= 0; i-- {
		if isOperationMethod(parts[i]) {
			return parts[i]
		}
	}
	return ""
}

func isOperationMethod(part string) bool {
	switch part {
	case "get", "put", "post", "delete", "options", "head", "patch", "trace":
		return true
	default:
		return false
	}
}

func semanticBool(sn *SemanticNode, key string) bool {
	return snStr(sn, key) == "true"
}

func parsePointerIndex(part string) (int, error) {
	if part == "" {
		return 0, fmt.Errorf("empty index")
	}
	i := 0
	for _, ch := range part {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("non-numeric index")
		}
		i = i*10 + int(ch-'0')
	}
	return i, nil
}
