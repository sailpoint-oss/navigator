package navigator

import (
	"strings"
	"unsafe"

	ts_yaml "github.com/sailpoint-oss/tree-sitter-openapi/bindings/go/openapi"
	ts_json "github.com/sailpoint-oss/tree-sitter-openapi/bindings/go/openapi_json"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// YAMLLanguage returns the tree-sitter Language for OpenAPI YAML.
func YAMLLanguage() *tree_sitter.Language {
	return tree_sitter.NewLanguage(unsafe.Pointer(ts_yaml.Language()))
}

// JSONLanguage returns the tree-sitter Language for OpenAPI JSON.
func JSONLanguage() *tree_sitter.Language {
	return tree_sitter.NewLanguage(unsafe.Pointer(ts_json.Language()))
}

// ParseTree parses using an existing tree-sitter CST for full source locations
// and break-glass access to tree-sitter nodes. Returns nil only when tree is
// nil; otherwise the returned Index always has Issues populated. Document may
// be nil when the root is not a mapping.
func ParseTree(tree *tree_sitter.Tree, content []byte, uri string, format FileFormat) *Index {
	if tree == nil {
		return nil
	}
	root := tree.RootNode()
	if root == nil {
		return nil
	}

	idx := newEmptyIndex()
	idx.uri = uri
	idx.content = content
	idx.tree = tree
	idx.Format = format

	tp := &treeParser{
		tree:    tree,
		content: content,
		format:  format,
	}

	semanticRoot := tp.buildSemanticTree(root)
	idx.semanticRoot = semanticRoot
	idx.Kind = DetectDocumentKind(semanticKeys(semanticRoot))
	if idx.Kind == DocumentKindArazzo {
		idx.Arazzo = parseArazzoDocumentFromSemantic(semanticRoot)
	} else {
		idx.Document = tp.parseDocument(semanticRoot)
		if idx.Document != nil {
			idx.Version = idx.Document.ParsedVersion
			if idx.Document.DocType == DocTypeFragment {
				idx.fragmentValue = semanticNodeOpenAPIValue(semanticRoot)
			}
		}
	}

	idx.indexPaths()
	idx.indexComponents()
	idx.indexTags()
	tp.collectRefs(root, idx, uri)

	idx.applyDefaultValidation()
	return idx
}

// ParseContent parses raw content using tree-sitter, creating the tree internally.
// Returns nil only for empty input; the returned Index always has Issues populated.
// Document may be nil when the root is not a mapping.
func ParseContent(content []byte, uri string) *Index {
	if len(content) == 0 {
		return nil
	}
	format := DetectFormat(uri, content)
	parser := tree_sitter.NewParser()
	defer parser.Close()

	var lang *tree_sitter.Language
	switch format {
	case FormatJSON:
		lang = JSONLanguage()
	default:
		lang = YAMLLanguage()
	}
	if err := parser.SetLanguage(lang); err != nil {
		return nil
	}
	tree := parser.Parse(content, nil)
	if tree == nil {
		return nil
	}
	return ParseTree(tree, content, uri, format)
}

type treeParser struct {
	tree    *tree_sitter.Tree
	content []byte
	format  FileFormat
}

func (tp *treeParser) nodeText(node *tree_sitter.Node) string {
	if node == nil {
		return ""
	}
	start := node.StartByte()
	end := node.EndByte()
	if int(start) >= len(tp.content) || int(end) > len(tp.content) {
		return ""
	}
	return string(tp.content[start:end])
}

func (tp *treeParser) buildSemanticTree(node *tree_sitter.Node) *SemanticNode {
	if node == nil {
		return nil
	}
	kind := node.Kind()

	// Skip document wrapper nodes
	if kind == "stream" || kind == "document" {
		for i := uint(0); i < node.ChildCount(); i++ {
			child := node.Child(i)
			if child != nil {
				k := child.Kind()
				if k != "---" && k != "..." && k != "comment" && k != "yaml_directive" && k != "tag_directive" {
					return tp.buildSemanticTree(child)
				}
			}
		}
		return nil
	}

	// Block/flow wrappers
	if kind == "block_node" || kind == "flow_node" {
		for i := uint(0); i < node.ChildCount(); i++ {
			child := node.Child(i)
			if child != nil {
				k := child.Kind()
				if k != "tag" && k != "anchor" {
					return tp.buildSemanticTree(child)
				}
			}
		}
		return nil
	}

	switch kind {
	case "block_mapping", "flow_mapping", "object":
		return tp.buildMapping(node)
	case "block_sequence", "flow_sequence", "array":
		return tp.buildSequence(node)
	case "block_mapping_pair", "flow_pair", "pair":
		return tp.buildMapping(node.Parent())
	default:
		return tp.buildScalar(node)
	}
}

func (tp *treeParser) buildMapping(node *tree_sitter.Node) *SemanticNode {
	if node == nil {
		return nil
	}
	sn := &SemanticNode{
		Kind:     NodeMapping,
		Range:    rangeFromTSNode(node),
		Children: make(map[string]*SemanticNode),
		CST:      node,
	}

	kind := node.Kind()
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		ck := child.Kind()

		var keyNode, valNode *tree_sitter.Node
		if ck == "block_mapping_pair" || ck == "flow_pair" || ck == "pair" {
			keyNode = child.ChildByFieldName("key")
			valNode = child.ChildByFieldName("value")
		} else if kind == "object" && ck == "pair" {
			keyNode = child.ChildByFieldName("key")
			valNode = child.ChildByFieldName("value")
		} else {
			continue
		}

		if keyNode == nil {
			continue
		}
		key := unquote(tp.nodeText(keyNode))
		var childSN *SemanticNode
		if valNode != nil {
			childSN = tp.buildSemanticTree(valNode)
		}
		if childSN == nil {
			childSN = &SemanticNode{Kind: NodeNull, Range: rangeFromTSNode(child), CST: valNode}
		}
		childSN.Key = key
		sn.Children[key] = childSN
	}
	return sn
}

func (tp *treeParser) buildSequence(node *tree_sitter.Node) *SemanticNode {
	if node == nil {
		return nil
	}
	sn := &SemanticNode{
		Kind:  NodeSequence,
		Range: rangeFromTSNode(node),
		CST:   node,
	}

	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		ck := child.Kind()
		if ck == "block_sequence_item" {
			for j := uint(0); j < child.ChildCount(); j++ {
				inner := child.Child(j)
				if inner != nil && inner.Kind() != "-" {
					item := tp.buildSemanticTree(inner)
					if item != nil {
						sn.Items = append(sn.Items, item)
					}
				}
			}
		} else if ck == "[" || ck == "]" || ck == "," || ck == "comment" {
			continue
		} else {
			item := tp.buildSemanticTree(child)
			if item != nil {
				sn.Items = append(sn.Items, item)
			}
		}
	}
	return sn
}

func (tp *treeParser) buildScalar(node *tree_sitter.Node) *SemanticNode {
	if node == nil {
		return nil
	}
	text := unquote(tp.nodeText(node))
	return &SemanticNode{
		Kind:  NodeScalar,
		Value: text,
		Range: rangeFromTSNode(node),
		CST:   node,
	}
}

func (tp *treeParser) parseDocument(sn *SemanticNode) *Document {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	doc := &Document{
		Paths:      make(map[string]*PathItem),
		Extensions: make(map[string]*Node),
		Loc:        locFromSN(sn),
	}

	var rootKeys []string
	for k := range sn.Children {
		rootKeys = append(rootKeys, k)
	}
	doc.DocType = DetectDocType(rootKeys)

	if v := sn.Get("openapi"); v != nil {
		doc.Version = v.StringValue()
		doc.ParsedVersion = VersionFromString(doc.Version)
	} else if v := sn.Get("swagger"); v != nil {
		doc.Version = v.StringValue()
		doc.ParsedVersion = VersionFromString(doc.Version)
	}

	if infoNode := sn.Get("info"); infoNode != nil {
		doc.Info = tp.parseInfoSN(infoNode)
	}

	if serversNode := sn.Get("servers"); serversNode != nil {
		doc.Servers = tp.parseServersSN(serversNode)
	}

	if pathsNode := sn.Get("paths"); pathsNode != nil && pathsNode.Kind == NodeMapping {
		for pathKey, itemSN := range pathsNode.Children {
			pi := tp.parsePathItemSN(itemSN)
			if pi != nil {
				pi.PathLoc = locFromSN(itemSN)
				doc.Paths[pathKey] = pi
			}
		}
	}

	if compNode := sn.Get("components"); compNode != nil {
		doc.Components = tp.parseComponentsSN(compNode)
	}

	if tagsNode := sn.Get("tags"); tagsNode != nil && tagsNode.Kind == NodeSequence {
		for _, item := range tagsNode.Items {
			t := Tag{
				Name:        snStr(item, "name"),
				Description: snDescription(item),
				Loc:         locFromSN(item),
			}
			doc.Tags = append(doc.Tags, t)
		}
	}

	if secNode := sn.Get("security"); secNode != nil {
		doc.Security = tp.parseSecuritySN(secNode)
	}

	if edNode := sn.Get("externalDocs"); edNode != nil {
		doc.ExternalDocs = tp.parseExternalDocsSN(edNode)
	}
	doc.Extensions = snExtensions(sn)

	return doc
}

func (tp *treeParser) parseInfoSN(sn *SemanticNode) *Info {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	info := &Info{
		Title:          snStr(sn, "title"),
		Description:    snDescription(sn),
		TermsOfService: snStr(sn, "termsOfService"),
		Version:        snStr(sn, "version"),
		Extensions:     snExtensions(sn),
		Loc:            locFromSN(sn),
		TitleLoc:       locFromSN(sn.Get("title")),
		VersionLoc:     locFromSN(sn.Get("version")),
	}
	if contactSN := sn.Get("contact"); contactSN != nil && contactSN.Kind == NodeMapping {
		info.Contact = &Contact{
			Name:  snStr(contactSN, "name"),
			URL:   snStr(contactSN, "url"),
			Email: snStr(contactSN, "email"),
			Loc:   locFromSN(contactSN),
		}
	}
	if licenseSN := sn.Get("license"); licenseSN != nil && licenseSN.Kind == NodeMapping {
		info.License = &License{
			Name:       snStr(licenseSN, "name"),
			Identifier: snStr(licenseSN, "identifier"),
			URL:        snStr(licenseSN, "url"),
			Loc:        locFromSN(licenseSN),
		}
	}
	return info
}

func (tp *treeParser) parseServersSN(sn *SemanticNode) []Server {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	var servers []Server
	for _, item := range sn.Items {
		if item.Kind != NodeMapping {
			continue
		}
		s := Server{
			URL:         snStr(item, "url"),
			Description: snDescription(item),
			Extensions:  snExtensions(item),
			Loc:         locFromSN(item),
			URLLoc:      locFromSN(item.Get("url")),
		}
		servers = append(servers, s)
	}
	return servers
}

func (tp *treeParser) parsePathItemSN(sn *SemanticNode) *PathItem {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	item := &PathItem{
		Summary:     snStr(sn, "summary"),
		Description: snDescription(sn),
		Ref:         snStr(sn, "$ref"),
		Extensions:  snExtensions(sn),
		Loc:         locFromSN(sn),
	}

	methods := []string{"get", "put", "post", "delete", "options", "head", "patch", "trace"}
	for _, m := range methods {
		if opSN := sn.Get(m); opSN != nil {
			op := tp.parseOperationSN(opSN, m)
			switch m {
			case "get":
				item.Get = op
			case "put":
				item.Put = op
			case "post":
				item.Post = op
			case "delete":
				item.Delete = op
			case "options":
				item.Options = op
			case "head":
				item.Head = op
			case "patch":
				item.Patch = op
			case "trace":
				item.Trace = op
			}
		}
	}

	if paramsSN := sn.Get("parameters"); paramsSN != nil {
		item.Parameters = tp.parseParametersSN(paramsSN)
	}
	return item
}

func (tp *treeParser) parseOperationSN(sn *SemanticNode, method string) *Operation {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	op := &Operation{
		OperationID:    snStr(sn, "operationId"),
		Summary:        snStr(sn, "summary"),
		Description:    snDescription(sn),
		Deprecated:     snStr(sn, "deprecated") == "true",
		Extensions:     snExtensions(sn),
		Loc:            locFromSN(sn),
		MethodLoc:      locFromSN(sn),
		OperationIDLoc: locFromSN(sn.Get("operationId")),
		TagsLoc:        locFromSN(sn.Get("tags")),
		ResponsesLoc:   locFromSN(sn.Get("responses")),
		ParametersLoc:  locFromSN(sn.Get("parameters")),
	}

	if tagsSN := sn.Get("tags"); tagsSN != nil && tagsSN.Kind == NodeSequence {
		for _, t := range tagsSN.Items {
			op.Tags = append(op.Tags, TagUsage{Name: t.StringValue(), Loc: locFromSN(t)})
		}
	}

	if paramsSN := sn.Get("parameters"); paramsSN != nil {
		op.Parameters = tp.parseParametersSN(paramsSN)
	}

	if rbSN := sn.Get("requestBody"); rbSN != nil {
		op.RequestBody = tp.parseRequestBodySN(rbSN)
	}

	if respSN := sn.Get("responses"); respSN != nil && respSN.Kind == NodeMapping {
		op.Responses = make(map[string]*Response)
		for code, rSN := range respSN.Children {
			r := tp.parseResponseSN(rSN)
			if r != nil {
				op.Responses[code] = r
			}
		}
	}

	if secSN := sn.Get("security"); secSN != nil {
		op.Security = tp.parseSecuritySN(secSN)
	}

	return op
}

func (tp *treeParser) parseParametersSN(sn *SemanticNode) []*Parameter {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	var params []*Parameter
	for _, item := range sn.Items {
		p := tp.parseParameterSN(item)
		if p != nil {
			params = append(params, p)
		}
	}
	return params
}

func (tp *treeParser) parseParameterSN(sn *SemanticNode) *Parameter {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	p := &Parameter{
		Name:            snStr(sn, "name"),
		In:              snStr(sn, "in"),
		Description:     snDescription(sn),
		Required:        snStr(sn, "required") == "true",
		Deprecated:      snStr(sn, "deprecated") == "true",
		AllowEmptyValue: snStr(sn, "allowEmptyValue") == "true",
		Ref:             snStr(sn, "$ref"),
		Extensions:      snExtensions(sn),
		Loc:             locFromSN(sn),
		NameLoc:         locFromSN(sn.Get("name")),
	}
	if schemaSN := sn.Get("schema"); schemaSN != nil {
		p.Schema = tp.parseSchemaSN(schemaSN)
	}
	if exSN := sn.Get("example"); exSN != nil {
		p.Example = snToNode(exSN)
	}
	if exsSN := sn.Get("examples"); exsSN != nil {
		p.Examples = tp.parseExamplesSN(exsSN)
	}
	return p
}

func (tp *treeParser) parseRequestBodySN(sn *SemanticNode) *RequestBody {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	rb := &RequestBody{
		Description: snDescription(sn),
		Required:    snStr(sn, "required") == "true",
		Ref:         snStr(sn, "$ref"),
		NameLoc:     locFromSN(sn),
		Loc:         locFromSN(sn),
	}
	rb.Content = tp.parseContentMediaTypesSN(sn.Get("content"))
	return rb
}

func (tp *treeParser) parseResponseSN(sn *SemanticNode) *Response {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	r := &Response{
		Description: snDescription(sn),
		Ref:         snStr(sn, "$ref"),
		Loc:         locFromSN(sn),
	}
	if headersSN := sn.Get("headers"); headersSN != nil && headersSN.Kind == NodeMapping {
		r.HeadersLoc = locFromSN(headersSN)
		r.Headers = tp.parseHeadersSN(headersSN)
	}
	r.Content = tp.parseContentMediaTypesSN(sn.Get("content"))
	if linksSN := sn.Get("links"); linksSN != nil && linksSN.Kind == NodeMapping {
		r.Links = tp.parseLinksSN(linksSN)
	}
	return r
}

func (tp *treeParser) parseContentMediaTypesSN(contentSN *SemanticNode) map[string]*MediaType {
	if contentSN == nil || contentSN.Kind != NodeMapping {
		return nil
	}
	result := make(map[string]*MediaType)
	for mt, mtSN := range contentSN.Children {
		mediaType := &MediaType{Loc: locFromSN(mtSN)}
		if schemaSN := mtSN.Get("schema"); schemaSN != nil {
			mediaType.Schema = tp.parseSchemaSN(schemaSN)
		}
		if exSN := mtSN.Get("example"); exSN != nil {
			mediaType.Example = snToNode(exSN)
		}
		if exsSN := mtSN.Get("examples"); exsSN != nil {
			mediaType.Examples = tp.parseExamplesSN(exsSN)
		}
		result[mt] = mediaType
	}
	return result
}

func (tp *treeParser) parseExamplesSN(sn *SemanticNode) map[string]*Example {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	examples := make(map[string]*Example)
	for name, eSN := range sn.Children {
		if eSN.Kind != NodeMapping {
			continue
		}
		ex := &Example{
			Summary:       snStr(eSN, "summary"),
			Description:   snDescription(eSN),
			ExternalValue: snStr(eSN, "externalValue"),
			Ref:           snStr(eSN, "$ref"),
			Loc:           locFromSN(eSN),
		}
		if valSN := eSN.Get("value"); valSN != nil {
			ex.Value = snToNode(valSN)
		}
		examples[name] = ex
	}
	return examples
}

func (tp *treeParser) parseHeadersSN(sn *SemanticNode) map[string]*Header {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	headers := make(map[string]*Header)
	for name, hSN := range sn.Children {
		if h := tp.parseHeaderSN(hSN); h != nil {
			h.NameLoc = locFromSN(hSN)
			headers[name] = h
		}
	}
	return headers
}

func (tp *treeParser) parseHeaderSN(sn *SemanticNode) *Header {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	h := &Header{
		Description: snDescription(sn),
		Required:    snStr(sn, "required") == "true",
		Deprecated:  snStr(sn, "deprecated") == "true",
		Ref:         snStr(sn, "$ref"),
		Loc:         locFromSN(sn),
	}
	if schemaSN := sn.Get("schema"); schemaSN != nil {
		h.Schema = tp.parseSchemaSN(schemaSN)
	}
	return h
}

func (tp *treeParser) parseLinksSN(sn *SemanticNode) map[string]*Link {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	links := make(map[string]*Link)
	for name, lSN := range sn.Children {
		if l := tp.parseLinkSN(lSN); l != nil {
			l.NameLoc = locFromSN(lSN)
			links[name] = l
		}
	}
	return links
}

func (tp *treeParser) parseLinkSN(sn *SemanticNode) *Link {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	return &Link{
		OperationRef:   snStr(sn, "operationRef"),
		OperationID:    snStr(sn, "operationId"),
		Description:    snDescription(sn),
		Ref:            snStr(sn, "$ref"),
		Loc:            locFromSN(sn),
		OperationIDLoc: locFromSN(sn.Get("operationId")),
	}
}

func (tp *treeParser) parseSchemaSN(sn *SemanticNode) *Schema {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	s := &Schema{
		Type:        snStr(sn, "type"),
		Format:      snStr(sn, "format"),
		Title:       snStr(sn, "title"),
		Description: snDescription(sn),
		Pattern:     snStr(sn, "pattern"),
		Nullable:    snStr(sn, "nullable") == "true",
		ReadOnly:    snStr(sn, "readOnly") == "true",
		WriteOnly:   snStr(sn, "writeOnly") == "true",
		Deprecated:  snStr(sn, "deprecated") == "true",
		Default:     snToNode(sn.Get("default")),
		Minimum:     snFloatValue(sn.Get("minimum")),
		Maximum:     snFloatValue(sn.Get("maximum")),
		Ref:         snStr(sn, "$ref"),
		Extensions:  snExtensions(sn),
		Loc:         locFromSN(sn),
		TypeLoc:     locFromSN(sn.Get("type")),
	}

	if itemsSN := sn.Get("items"); itemsSN != nil {
		s.Items = tp.parseSchemaSN(itemsSN)
	}
	if notSN := sn.Get("not"); notSN != nil {
		s.Not = tp.parseSchemaSN(notSN)
	}

	if propsSN := sn.Get("properties"); propsSN != nil && propsSN.Kind == NodeMapping {
		s.Properties = make(map[string]*Schema)
		for name, pSN := range propsSN.Children {
			ps := tp.parseSchemaSN(pSN)
			if ps != nil {
				ps.NameLoc = locFromSN(pSN)
				s.Properties[name] = ps
			}
		}
	}

	s.AllOf = tp.parseSchemaListSN(sn.Get("allOf"))
	s.AnyOf = tp.parseSchemaListSN(sn.Get("anyOf"))
	s.OneOf = tp.parseSchemaListSN(sn.Get("oneOf"))

	if reqSN := sn.Get("required"); reqSN != nil && reqSN.Kind == NodeSequence {
		s.HasRequired = true
		for _, item := range reqSN.Items {
			s.Required = append(s.Required, item.StringValue())
		}
	}

	if enumSN := sn.Get("enum"); enumSN != nil && enumSN.Kind == NodeSequence {
		for _, item := range enumSN.Items {
			s.Enum = append(s.Enum, item.StringValue())
		}
	}

	if exSN := sn.Get("example"); exSN != nil {
		s.Example = snToNode(exSN)
	}

	return s
}

func (tp *treeParser) parseSchemaListSN(sn *SemanticNode) []*Schema {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	var schemas []*Schema
	for _, item := range sn.Items {
		if s := tp.parseSchemaSN(item); s != nil {
			schemas = append(schemas, s)
		}
	}
	return schemas
}

func (tp *treeParser) parseComponentsSN(sn *SemanticNode) *Components {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	c := &Components{Loc: locFromSN(sn)}

	if n := sn.Get("schemas"); n != nil && n.Kind == NodeMapping {
		c.Schemas = make(map[string]*Schema)
		for name, sSN := range n.Children {
			if s := tp.parseSchemaSN(sSN); s != nil {
				c.Schemas[name] = s
			}
		}
	}

	if n := sn.Get("parameters"); n != nil && n.Kind == NodeMapping {
		c.Parameters = make(map[string]*Parameter)
		for name, pSN := range n.Children {
			if p := tp.parseParameterSN(pSN); p != nil {
				p.NameLoc = locFromSN(pSN)
				c.Parameters[name] = p
			}
		}
	}

	if n := sn.Get("examples"); n != nil && n.Kind == NodeMapping {
		c.Examples = tp.parseExamplesSN(n)
	}

	if n := sn.Get("requestBodies"); n != nil && n.Kind == NodeMapping {
		c.RequestBodies = make(map[string]*RequestBody)
		for name, rbSN := range n.Children {
			if rb := tp.parseRequestBodySN(rbSN); rb != nil {
				rb.NameLoc = locFromSN(rbSN)
				c.RequestBodies[name] = rb
			}
		}
	}

	if n := sn.Get("responses"); n != nil && n.Kind == NodeMapping {
		c.Responses = make(map[string]*Response)
		for name, rSN := range n.Children {
			if r := tp.parseResponseSN(rSN); r != nil {
				r.NameLoc = locFromSN(rSN)
				c.Responses[name] = r
			}
		}
	}

	if n := sn.Get("headers"); n != nil && n.Kind == NodeMapping {
		c.Headers = tp.parseHeadersSN(n)
	}

	if n := sn.Get("securitySchemes"); n != nil && n.Kind == NodeMapping {
		c.SecuritySchemes = make(map[string]*SecurityScheme)
		for name, ssSN := range n.Children {
			if ss := tp.parseSecuritySchemeSN(ssSN); ss != nil {
				ss.NameLoc = locFromSN(ssSN)
				c.SecuritySchemes[name] = ss
			}
		}
	}

	if n := sn.Get("links"); n != nil && n.Kind == NodeMapping {
		c.Links = tp.parseLinksSN(n)
	}

	if n := sn.Get("pathItems"); n != nil && n.Kind == NodeMapping {
		c.PathItems = make(map[string]*PathItem)
		for name, piSN := range n.Children {
			if pi := tp.parsePathItemSN(piSN); pi != nil {
				c.PathItems[name] = pi
			}
		}
	}

	return c
}

func (tp *treeParser) parseSecuritySchemeSN(sn *SemanticNode) *SecurityScheme {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	ss := &SecurityScheme{
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
	if flowsSN := sn.Get("flows"); flowsSN != nil && flowsSN.Kind == NodeMapping {
		ss.Flows = tp.parseOAuthFlowsSN(flowsSN)
	}
	return ss
}

func (tp *treeParser) parseOAuthFlowsSN(sn *SemanticNode) *OAuthFlows {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	flows := &OAuthFlows{Loc: locFromSN(sn)}
	flows.Implicit = tp.parseOAuthFlowSN(sn.Get("implicit"))
	flows.Password = tp.parseOAuthFlowSN(sn.Get("password"))
	flows.ClientCredentials = tp.parseOAuthFlowSN(sn.Get("clientCredentials"))
	flows.AuthorizationCode = tp.parseOAuthFlowSN(sn.Get("authorizationCode"))
	return flows
}

func (tp *treeParser) parseOAuthFlowSN(sn *SemanticNode) *OAuthFlow {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	flow := &OAuthFlow{
		AuthorizationURL:    snStr(sn, "authorizationUrl"),
		TokenURL:            snStr(sn, "tokenUrl"),
		RefreshURL:          snStr(sn, "refreshUrl"),
		Loc:                 locFromSN(sn),
		AuthorizationURLLoc: locFromSN(sn.Get("authorizationUrl")),
		TokenURLLoc:         locFromSN(sn.Get("tokenUrl")),
	}
	if scopesSN := sn.Get("scopes"); scopesSN != nil && scopesSN.Kind == NodeMapping {
		flow.Scopes = make(map[string]string)
		for name, scopeSN := range scopesSN.Children {
			flow.Scopes[name] = scopeSN.StringValue()
		}
	}
	return flow
}

func (tp *treeParser) parseSecuritySN(sn *SemanticNode) []SecurityRequirement {
	if sn == nil || sn.Kind != NodeSequence {
		return nil
	}
	var reqs []SecurityRequirement
	for _, item := range sn.Items {
		if item.Kind != NodeMapping {
			continue
		}
		req := SecurityRequirement{Loc: locFromSN(item)}
		for name, child := range item.Children {
			entry := SecurityRequirementEntry{Name: name, NameLoc: locFromSN(child)}
			if child.Kind == NodeSequence {
				for _, s := range child.Items {
					entry.Scopes = append(entry.Scopes, s.StringValue())
					entry.ScopeLocs = append(entry.ScopeLocs, locFromSN(s))
				}
			}
			req.Entries = append(req.Entries, entry)
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func (tp *treeParser) parseExternalDocsSN(sn *SemanticNode) *ExternalDocs {
	if sn == nil || sn.Kind != NodeMapping {
		return nil
	}
	return &ExternalDocs{
		Description: snDescription(sn),
		URL:         snStr(sn, "url"),
		Loc:         locFromSN(sn),
	}
}

func (tp *treeParser) collectRefs(root *tree_sitter.Node, idx *Index, uri string) {
	if root == nil {
		return
	}
	tp.walkForRefs(root, idx, nil)
	for i := range idx.AllRefs {
		idx.AllRefs[i].URI = uri
	}
	for target := range idx.Refs {
		for i := range idx.Refs[target] {
			idx.Refs[target][i].URI = uri
		}
	}
}

func (tp *treeParser) walkForRefs(node *tree_sitter.Node, idx *Index, path []string) {
	if node == nil {
		return
	}
	kind := node.Kind()

	isKVPair := (tp.format == FormatYAML && (kind == "block_mapping_pair" || kind == "flow_pair")) ||
		(tp.format == FormatJSON && kind == "pair")

	if isKVPair {
		keyNode := node.ChildByFieldName("key")
		valNode := node.ChildByFieldName("value")
		if keyNode != nil && valNode != nil {
			keyText := unquote(tp.nodeText(keyNode))
			if keyText == "$ref" {
				refTarget := unquote(tp.nodeText(valNode))
				usage := RefUsage{
					Loc:    LocFromNode(valNode),
					Target: refTarget,
					From:   "/" + strings.Join(path, "/"),
				}
				idx.Refs[refTarget] = append(idx.Refs[refTarget], usage)
				idx.AllRefs = append(idx.AllRefs, usage)
				return
			}
			childPath := append(path, EscapeJSONPointer(keyText))
			tp.walkForRefs(valNode, idx, childPath)
			return
		}
	}

	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child != nil {
			tp.walkForRefs(child, idx, path)
		}
	}
}

// --- helpers ---

func rangeFromTSNode(node *tree_sitter.Node) Range {
	if node == nil {
		return Range{}
	}
	start := node.StartPosition()
	end := node.EndPosition()
	return Range{
		Start: Position{Line: uint32(start.Row), Character: uint32(start.Column)},
		End:   Position{Line: uint32(end.Row), Character: uint32(end.Column)},
	}
}

func locFromSN(sn *SemanticNode) Loc {
	if sn == nil {
		return Loc{}
	}
	l := Loc{Range: sn.Range}
	if sn.CST != nil {
		l.Node = sn.CST
	}
	return l
}

func snToNode(sn *SemanticNode) *Node {
	if sn == nil {
		return nil
	}
	return &Node{
		Value: sn.StringValue(),
		Loc:   locFromSN(sn),
	}
}

func snStr(sn *SemanticNode, key string) string {
	if sn == nil {
		return ""
	}
	child := sn.Get(key)
	if child == nil {
		return ""
	}
	return child.StringValue()
}

func snDescription(sn *SemanticNode) DescriptionValue {
	child := sn.Get("description")
	if child == nil {
		return DescriptionValue{}
	}
	return DescriptionValue{
		Text: child.StringValue(),
		Loc:  locFromSN(child),
	}
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
