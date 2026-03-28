package navigator

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

// Parse parses raw YAML/JSON content into a fully indexed OpenAPI document.
// Uses the standalone yaml.v3 parser (no CGO or tree-sitter required).
// Returns nil when yaml.Unmarshal fails; otherwise Issues is always populated.
// Document may be nil when the root is not a mapping.
func Parse(content []byte) *Index {
	return ParseURI("", content)
}

// ParseAndIndex is an alias for Parse, matching telescope's standalone API.
func ParseAndIndex(content []byte) *Index {
	return Parse(content)
}

// ParseURI is like Parse but associates a URI with the result for ref tracking.
// Returns nil when yaml.Unmarshal fails; otherwise Issues is always populated.
// Document may be nil when the root is not a mapping.
func ParseURI(uri string, content []byte) *Index {
	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil
	}
	if root.Kind == 0 {
		return nil
	}

	idx := newEmptyIndex()
	idx.uri = uri
	idx.content = content
	idx.Format = DetectFormat(uri, content)

	sp := &standaloneParser{root: &root, uri: uri}
	idx.semanticRoot = semanticRootFromYAML(sp.docNode())
	idx.Kind = DetectDocumentKind(yamlMapKeys(sp.docNode()))
	if idx.Kind == DocumentKindArazzo {
		idx.Arazzo = parseArazzoDocumentFromSemantic(idx.semanticRoot)
	} else {
		idx.Document = sp.parseDocument()
		if idx.Document != nil {
			idx.Version = idx.Document.ParsedVersion
			if idx.Document.DocType == DocTypeFragment {
				idx.fragmentValue = semanticNodeOpenAPIValue(idx.semanticRoot)
			}
		}
	}

	idx.indexPaths()
	idx.indexComponents()
	idx.indexTags()
	sp.collectRefs(idx)

	idx.applyDefaultValidation()
	return idx
}

const maxSchemaDepth = 64

type standaloneParser struct {
	root *yaml.Node
	uri  string
}

func (sp *standaloneParser) parseDocument() *Document {
	node := sp.docNode()
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	doc := &Document{
		Paths:      make(map[string]*PathItem),
		Extensions: make(map[string]*Node),
		Loc:        yamlLoc(node),
	}

	rootKeys := yamlMapKeys(node)
	doc.DocType = DetectDocType(rootKeys)

	if v := yamlStr(node, "openapi"); v != "" {
		doc.Version = v
		doc.ParsedVersion = VersionFromString(v)
	} else if v := yamlStr(node, "swagger"); v != "" {
		doc.Version = v
		doc.ParsedVersion = VersionFromString(v)
	}

	if infoNode := yamlMapGet(node, "info"); infoNode != nil {
		doc.Info = sp.parseInfo(infoNode)
	}

	if serversNode := yamlMapGet(node, "servers"); serversNode != nil {
		doc.Servers = sp.parseServers(serversNode)
	}

	if pathsNode := yamlMapGet(node, "paths"); pathsNode != nil {
		doc.Paths = sp.parsePaths(pathsNode)
	}

	if compNode := yamlMapGet(node, "components"); compNode != nil {
		doc.Components = sp.parseComponents(compNode)
	}

	if tagsNode := yamlMapGet(node, "tags"); tagsNode != nil {
		doc.Tags = sp.parseTags(tagsNode)
	}

	if secNode := yamlMapGet(node, "security"); secNode != nil {
		doc.Security = sp.parseSecurity(secNode)
	}

	if edNode := yamlMapGet(node, "externalDocs"); edNode != nil {
		doc.ExternalDocs = sp.parseExternalDocs(edNode)
	}

	doc.Schemes = yamlStringSlice(node, "schemes")
	doc.Extensions = yamlExtensions(node)

	return doc
}

func (sp *standaloneParser) parseInfo(node *yaml.Node) *Info {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	info := &Info{
		Loc:        yamlLoc(node),
		Extensions: make(map[string]*Node),
	}
	info.Title = yamlStr(node, "title")
	info.TitleLoc = yamlKeyLoc(node, "title")
	info.Description = yamlDescription(node)
	info.TermsOfService = yamlStr(node, "termsOfService")
	info.Version = yamlStr(node, "version")
	info.VersionLoc = yamlKeyLoc(node, "version")

	if contactNode := yamlMapGet(node, "contact"); contactNode != nil {
		info.Contact = &Contact{
			Name:  yamlStr(contactNode, "name"),
			URL:   yamlStr(contactNode, "url"),
			Email: yamlStr(contactNode, "email"),
			Loc:   yamlLoc(contactNode),
		}
	}
	if licenseNode := yamlMapGet(node, "license"); licenseNode != nil {
		info.License = &License{
			Name:       yamlStr(licenseNode, "name"),
			Identifier: yamlStr(licenseNode, "identifier"),
			URL:        yamlStr(licenseNode, "url"),
			Loc:        yamlLoc(licenseNode),
		}
	}
	info.Extensions = yamlExtensions(node)
	return info
}

func (sp *standaloneParser) parseServers(node *yaml.Node) []Server {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}
	var servers []Server
	for _, item := range node.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}
		s := Server{
			URL:         yamlStr(item, "url"),
			Description: yamlDescription(item),
			Loc:         yamlLoc(item),
			URLLoc:      yamlKeyLoc(item, "url"),
			Extensions:  yamlExtensions(item),
		}
		if varsNode := yamlMapGet(item, "variables"); varsNode != nil && varsNode.Kind == yaml.MappingNode {
			s.Variables = make(map[string]*ServerVariable)
			for i := 0; i+1 < len(varsNode.Content); i += 2 {
				name := varsNode.Content[i].Value
				vNode := varsNode.Content[i+1]
				sv := &ServerVariable{
					Default:     yamlStr(vNode, "default"),
					Description: yamlDescription(vNode),
					Enum:        yamlStringSlice(vNode, "enum"),
					Loc:         yamlLoc(vNode),
				}
				s.Variables[name] = sv
			}
		}
		servers = append(servers, s)
	}
	return servers
}

func (sp *standaloneParser) parsePaths(node *yaml.Node) map[string]*PathItem {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	paths := make(map[string]*PathItem)
	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		pathTemplate := keyNode.Value
		item := sp.parsePathItem(valNode)
		if item != nil {
			item.PathLoc = yamlLoc(keyNode)
			paths[pathTemplate] = item
		}
	}
	return paths
}

func (sp *standaloneParser) parsePathItem(node *yaml.Node) *PathItem {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	item := &PathItem{
		Summary:     yamlStr(node, "summary"),
		Description: yamlDescription(node),
		Ref:         yamlStr(node, "$ref"),
		Loc:         yamlLoc(node),
		Extensions:  yamlExtensions(node),
	}

	methods := []string{"get", "put", "post", "delete", "options", "head", "patch", "trace"}
	for _, m := range methods {
		if opNode := yamlMapGet(node, m); opNode != nil {
			op := sp.parseOperation(opNode, m)
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

	if paramsNode := yamlMapGet(node, "parameters"); paramsNode != nil {
		item.Parameters = sp.parseParameters(paramsNode)
	}
	if serversNode := yamlMapGet(node, "servers"); serversNode != nil {
		item.Servers = sp.parseServers(serversNode)
	}
	return item
}

func (sp *standaloneParser) parseOperation(node *yaml.Node, method string) *Operation {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	op := &Operation{
		OperationID:    yamlStr(node, "operationId"),
		Summary:        yamlStr(node, "summary"),
		Description:    yamlDescription(node),
		Deprecated:     yamlStr(node, "deprecated") == "true",
		Loc:            yamlLoc(node),
		MethodLoc:      yamlKeyLoc(node, method),
		OperationIDLoc: yamlValLoc(node, "operationId"),
		TagsLoc:        yamlKeyLoc(node, "tags"),
		ResponsesLoc:   yamlKeyLoc(node, "responses"),
		ParametersLoc:  yamlKeyLoc(node, "parameters"),
		Extensions:     yamlExtensions(node),
	}

	if tagsNode := yamlMapGet(node, "tags"); tagsNode != nil && tagsNode.Kind == yaml.SequenceNode {
		for _, t := range tagsNode.Content {
			op.Tags = append(op.Tags, TagUsage{Name: t.Value, Loc: yamlLoc(t)})
		}
	}

	if paramsNode := yamlMapGet(node, "parameters"); paramsNode != nil {
		op.Parameters = sp.parseParameters(paramsNode)
	}

	if rbNode := yamlMapGet(node, "requestBody"); rbNode != nil {
		op.RequestBody = sp.parseRequestBody(rbNode)
	}

	if respNode := yamlMapGet(node, "responses"); respNode != nil {
		op.Responses = sp.parseResponses(respNode)
	}

	if secNode := yamlMapGet(node, "security"); secNode != nil {
		op.Security = sp.parseSecurity(secNode)
	}

	if serversNode := yamlMapGet(node, "servers"); serversNode != nil {
		op.Servers = sp.parseServers(serversNode)
	}

	if edNode := yamlMapGet(node, "externalDocs"); edNode != nil {
		op.ExternalDocs = sp.parseExternalDocs(edNode)
	}

	return op
}

func (sp *standaloneParser) parseParameters(node *yaml.Node) []*Parameter {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}
	var params []*Parameter
	for _, item := range node.Content {
		p := sp.parseParameter(item)
		if p != nil {
			params = append(params, p)
		}
	}
	return params
}

func (sp *standaloneParser) parseParameter(node *yaml.Node) *Parameter {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	p := &Parameter{
		Name:            yamlStr(node, "name"),
		In:              yamlStr(node, "in"),
		Description:     yamlDescription(node),
		Required:        yamlStr(node, "required") == "true",
		Deprecated:      yamlStr(node, "deprecated") == "true",
		AllowEmptyValue: yamlStr(node, "allowEmptyValue") == "true",
		Ref:             yamlStr(node, "$ref"),
		Loc:             yamlLoc(node),
		NameLoc:         yamlValLoc(node, "name"),
		Extensions:      yamlExtensions(node),
	}
	if schemaNode := yamlMapGet(node, "schema"); schemaNode != nil {
		p.Schema = sp.parseSchema(schemaNode)
	}
	if exNode := yamlMapGet(node, "example"); exNode != nil {
		p.Example = yamlNode(exNode)
	}
	if exsNode := yamlMapGet(node, "examples"); exsNode != nil {
		p.Examples = sp.parseExamples(exsNode)
	}
	return p
}

func (sp *standaloneParser) parseRequestBody(node *yaml.Node) *RequestBody {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	rb := &RequestBody{
		Description: yamlDescription(node),
		Required:    yamlStr(node, "required") == "true",
		Ref:         yamlStr(node, "$ref"),
		Loc:         yamlLoc(node),
	}
	if contentNode := yamlMapGet(node, "content"); contentNode != nil {
		rb.Content = sp.parseMediaTypes(contentNode)
	}
	return rb
}

func (sp *standaloneParser) parseResponses(node *yaml.Node) map[string]*Response {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	responses := make(map[string]*Response)
	for i := 0; i+1 < len(node.Content); i += 2 {
		code := node.Content[i].Value
		rNode := node.Content[i+1]
		r := sp.parseResponse(rNode)
		if r != nil {
			r.CodeLoc = yamlLoc(node.Content[i])
			responses[code] = r
		}
	}
	return responses
}

func (sp *standaloneParser) parseResponse(node *yaml.Node) *Response {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	r := &Response{
		Description: yamlDescription(node),
		Ref:         yamlStr(node, "$ref"),
		Loc:         yamlLoc(node),
		HeadersLoc:  yamlKeyLoc(node, "headers"),
		Extensions:  yamlExtensions(node),
	}
	if contentNode := yamlMapGet(node, "content"); contentNode != nil {
		r.Content = sp.parseMediaTypes(contentNode)
	}
	if headersNode := yamlMapGet(node, "headers"); headersNode != nil && headersNode.Kind == yaml.MappingNode {
		r.Headers = make(map[string]*Header)
		for i := 0; i+1 < len(headersNode.Content); i += 2 {
			name := headersNode.Content[i].Value
			hNode := headersNode.Content[i+1]
			h := sp.parseHeader(hNode)
			if h != nil {
				h.NameLoc = yamlLoc(headersNode.Content[i])
				r.Headers[name] = h
			}
		}
	}
	if linksNode := yamlMapGet(node, "links"); linksNode != nil && linksNode.Kind == yaml.MappingNode {
		r.Links = make(map[string]*Link)
		for i := 0; i+1 < len(linksNode.Content); i += 2 {
			name := linksNode.Content[i].Value
			lNode := linksNode.Content[i+1]
			l := sp.parseLink(lNode)
			if l != nil {
				l.NameLoc = yamlLoc(linksNode.Content[i])
				r.Links[name] = l
			}
		}
	}
	return r
}

func (sp *standaloneParser) parseHeader(node *yaml.Node) *Header {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	h := &Header{
		Description: yamlDescription(node),
		Required:    yamlStr(node, "required") == "true",
		Deprecated:  yamlStr(node, "deprecated") == "true",
		Ref:         yamlStr(node, "$ref"),
		Loc:         yamlLoc(node),
	}
	if schemaNode := yamlMapGet(node, "schema"); schemaNode != nil {
		h.Schema = sp.parseSchema(schemaNode)
	}
	return h
}

func (sp *standaloneParser) parseLink(node *yaml.Node) *Link {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	return &Link{
		OperationRef:   yamlStr(node, "operationRef"),
		OperationID:    yamlStr(node, "operationId"),
		Description:    yamlDescription(node),
		Ref:            yamlStr(node, "$ref"),
		Loc:            yamlLoc(node),
		OperationIDLoc: yamlValLoc(node, "operationId"),
	}
}

func (sp *standaloneParser) parseMediaTypes(node *yaml.Node) map[string]*MediaType {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	mts := make(map[string]*MediaType)
	for i := 0; i+1 < len(node.Content); i += 2 {
		mediaType := node.Content[i].Value
		mtNode := node.Content[i+1]
		mt := &MediaType{Loc: yamlLoc(mtNode)}
		if schemaNode := yamlMapGet(mtNode, "schema"); schemaNode != nil {
			mt.Schema = sp.parseSchema(schemaNode)
		}
		if exNode := yamlMapGet(mtNode, "example"); exNode != nil {
			mt.Example = yamlNode(exNode)
		}
		if exsNode := yamlMapGet(mtNode, "examples"); exsNode != nil {
			mt.Examples = sp.parseExamples(exsNode)
		}
		mts[mediaType] = mt
	}
	return mts
}

func (sp *standaloneParser) parseSchema(node *yaml.Node) *Schema {
	return sp.parseSchemaDepth(node, 0)
}

func (sp *standaloneParser) parseSchemaDepth(node *yaml.Node, depth int) *Schema {
	if node == nil || depth > maxSchemaDepth {
		return nil
	}
	if node.Kind != yaml.MappingNode {
		return nil
	}
	s := &Schema{
		Type:        yamlStr(node, "type"),
		Format:      yamlStr(node, "format"),
		Title:       yamlStr(node, "title"),
		Description: yamlDescription(node),
		Pattern:     yamlStr(node, "pattern"),
		Nullable:    yamlStr(node, "nullable") == "true",
		ReadOnly:    yamlStr(node, "readOnly") == "true",
		WriteOnly:   yamlStr(node, "writeOnly") == "true",
		Deprecated:  yamlStr(node, "deprecated") == "true",
		Ref:         yamlStr(node, "$ref"),
		Loc:         yamlLoc(node),
		TypeLoc:     yamlValLoc(node, "type"),
		Extensions:  yamlExtensions(node),
	}

	s.Enum = yamlStringSlice(node, "enum")
	s.Required = yamlStringSlice(node, "required")

	if defNode := yamlMapGet(node, "default"); defNode != nil {
		s.Default = yamlNode(defNode)
	}
	if exNode := yamlMapGet(node, "example"); exNode != nil {
		s.Example = yamlNode(exNode)
	}

	if yamlMapGet(node, "const") != nil {
		s.HasConst = true
	}

	s.MinLength = yamlIntPtr(node, "minLength")
	s.MaxLength = yamlIntPtr(node, "maxLength")
	s.MinItems = yamlIntPtr(node, "minItems")
	s.MaxItems = yamlIntPtr(node, "maxItems")
	s.MaxProperties = yamlIntPtr(node, "maxProperties")
	s.Minimum = yamlFloatPtr(node, "minimum")
	s.Maximum = yamlFloatPtr(node, "maximum")
	s.ExclusiveMinimum = yamlFloatPtr(node, "exclusiveMinimum")
	s.ExclusiveMaximum = yamlFloatPtr(node, "exclusiveMaximum")

	if propsNode := yamlMapGet(node, "properties"); propsNode != nil && propsNode.Kind == yaml.MappingNode {
		s.Properties = make(map[string]*Schema)
		for i := 0; i+1 < len(propsNode.Content); i += 2 {
			name := propsNode.Content[i].Value
			propSchema := sp.parseSchemaDepth(propsNode.Content[i+1], depth+1)
			if propSchema != nil {
				propSchema.NameLoc = yamlLoc(propsNode.Content[i])
				s.Properties[name] = propSchema
			}
		}
	}

	if apNode := yamlMapGet(node, "additionalProperties"); apNode != nil {
		if apNode.Kind == yaml.ScalarNode && apNode.Value == "false" {
			s.AdditionalPropertiesFalse = true
		} else if apNode.Kind == yaml.MappingNode {
			s.AdditionalProperties = sp.parseSchemaDepth(apNode, depth+1)
		}
	}

	if upNode := yamlMapGet(node, "unevaluatedProperties"); upNode != nil {
		if upNode.Kind == yaml.ScalarNode && upNode.Value == "false" {
			s.UnevaluatedPropertiesFalse = true
		} else if upNode.Kind == yaml.MappingNode {
			s.UnevaluatedProperties = sp.parseSchemaDepth(upNode, depth+1)
		}
	}

	if itemsNode := yamlMapGet(node, "items"); itemsNode != nil {
		s.Items = sp.parseSchemaDepth(itemsNode, depth+1)
	}
	if notNode := yamlMapGet(node, "not"); notNode != nil {
		s.Not = sp.parseSchemaDepth(notNode, depth+1)
	}
	s.AllOf = sp.parseSchemaListDepth(yamlMapGet(node, "allOf"), depth+1)
	s.AnyOf = sp.parseSchemaListDepth(yamlMapGet(node, "anyOf"), depth+1)
	s.OneOf = sp.parseSchemaListDepth(yamlMapGet(node, "oneOf"), depth+1)

	if discNode := yamlMapGet(node, "discriminator"); discNode != nil && discNode.Kind == yaml.MappingNode {
		s.Discriminator = &Discriminator{
			PropertyName: yamlStr(discNode, "propertyName"),
			Loc:          yamlLoc(discNode),
		}
		if mappingNode := yamlMapGet(discNode, "mapping"); mappingNode != nil && mappingNode.Kind == yaml.MappingNode {
			s.Discriminator.Mapping = make(map[string]string)
			for i := 0; i+1 < len(mappingNode.Content); i += 2 {
				s.Discriminator.Mapping[mappingNode.Content[i].Value] = mappingNode.Content[i+1].Value
			}
		}
	}

	if edNode := yamlMapGet(node, "externalDocs"); edNode != nil {
		s.ExternalDocs = sp.parseExternalDocs(edNode)
	}

	return s
}

func (sp *standaloneParser) parseSchemaList(node *yaml.Node) []*Schema {
	return sp.parseSchemaListDepth(node, 0)
}

func (sp *standaloneParser) parseSchemaListDepth(node *yaml.Node, depth int) []*Schema {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}
	var schemas []*Schema
	for _, item := range node.Content {
		if s := sp.parseSchemaDepth(item, depth); s != nil {
			schemas = append(schemas, s)
		}
	}
	return schemas
}

func (sp *standaloneParser) parseComponents(node *yaml.Node) *Components {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	c := &Components{Loc: yamlLoc(node)}

	if n := yamlMapGet(node, "schemas"); n != nil && n.Kind == yaml.MappingNode {
		c.Schemas = make(map[string]*Schema)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			s := sp.parseSchema(n.Content[i+1])
			if s != nil {
				s.NameLoc = yamlLoc(n.Content[i])
				c.Schemas[name] = s
			}
		}
	}

	if n := yamlMapGet(node, "responses"); n != nil && n.Kind == yaml.MappingNode {
		c.Responses = make(map[string]*Response)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			r := sp.parseResponse(n.Content[i+1])
			if r != nil {
				r.NameLoc = yamlLoc(n.Content[i])
				c.Responses[name] = r
			}
		}
	}

	if n := yamlMapGet(node, "parameters"); n != nil && n.Kind == yaml.MappingNode {
		c.Parameters = make(map[string]*Parameter)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			p := sp.parseParameter(n.Content[i+1])
			if p != nil {
				p.NameLoc = yamlLoc(n.Content[i])
				c.Parameters[name] = p
			}
		}
	}

	if n := yamlMapGet(node, "examples"); n != nil && n.Kind == yaml.MappingNode {
		c.Examples = sp.parseExamples(n)
	}

	if n := yamlMapGet(node, "requestBodies"); n != nil && n.Kind == yaml.MappingNode {
		c.RequestBodies = make(map[string]*RequestBody)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			rb := sp.parseRequestBody(n.Content[i+1])
			if rb != nil {
				rb.NameLoc = yamlLoc(n.Content[i])
				c.RequestBodies[name] = rb
			}
		}
	}

	if n := yamlMapGet(node, "headers"); n != nil && n.Kind == yaml.MappingNode {
		c.Headers = make(map[string]*Header)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			h := sp.parseHeader(n.Content[i+1])
			if h != nil {
				h.NameLoc = yamlLoc(n.Content[i])
				c.Headers[name] = h
			}
		}
	}

	if n := yamlMapGet(node, "securitySchemes"); n != nil && n.Kind == yaml.MappingNode {
		c.SecuritySchemes = make(map[string]*SecurityScheme)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			ss := sp.parseSecurityScheme(n.Content[i+1])
			if ss != nil {
				ss.NameLoc = yamlLoc(n.Content[i])
				c.SecuritySchemes[name] = ss
			}
		}
	}

	if n := yamlMapGet(node, "links"); n != nil && n.Kind == yaml.MappingNode {
		c.Links = make(map[string]*Link)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			l := sp.parseLink(n.Content[i+1])
			if l != nil {
				l.NameLoc = yamlLoc(n.Content[i])
				c.Links[name] = l
			}
		}
	}

	if n := yamlMapGet(node, "pathItems"); n != nil && n.Kind == yaml.MappingNode {
		c.PathItems = make(map[string]*PathItem)
		for i := 0; i+1 < len(n.Content); i += 2 {
			name := n.Content[i].Value
			pi := sp.parsePathItem(n.Content[i+1])
			if pi != nil {
				c.PathItems[name] = pi
			}
		}
	}

	return c
}

func (sp *standaloneParser) parseSecurityScheme(node *yaml.Node) *SecurityScheme {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	ss := &SecurityScheme{
		Type:             yamlStr(node, "type"),
		Description:      yamlDescription(node),
		Name:             yamlStr(node, "name"),
		In:               yamlStr(node, "in"),
		Scheme:           yamlStr(node, "scheme"),
		BearerFormat:     yamlStr(node, "bearerFormat"),
		OpenIDConnectURL: yamlStr(node, "openIdConnectUrl"),
		Ref:              yamlStr(node, "$ref"),
		Loc:              yamlLoc(node),
		Extensions:       yamlExtensions(node),
	}
	if flowsNode := yamlMapGet(node, "flows"); flowsNode != nil && flowsNode.Kind == yaml.MappingNode {
		ss.Flows = sp.parseOAuthFlows(flowsNode)
	}
	return ss
}

func (sp *standaloneParser) parseOAuthFlows(node *yaml.Node) *OAuthFlows {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	flows := &OAuthFlows{Loc: yamlLoc(node)}
	if n := yamlMapGet(node, "implicit"); n != nil {
		flows.Implicit = sp.parseOAuthFlow(n)
	}
	if n := yamlMapGet(node, "password"); n != nil {
		flows.Password = sp.parseOAuthFlow(n)
	}
	if n := yamlMapGet(node, "clientCredentials"); n != nil {
		flows.ClientCredentials = sp.parseOAuthFlow(n)
	}
	if n := yamlMapGet(node, "authorizationCode"); n != nil {
		flows.AuthorizationCode = sp.parseOAuthFlow(n)
	}
	return flows
}

func (sp *standaloneParser) parseOAuthFlow(node *yaml.Node) *OAuthFlow {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	flow := &OAuthFlow{
		AuthorizationURL:    yamlStr(node, "authorizationUrl"),
		TokenURL:            yamlStr(node, "tokenUrl"),
		RefreshURL:          yamlStr(node, "refreshUrl"),
		Loc:                 yamlLoc(node),
		AuthorizationURLLoc: yamlValLoc(node, "authorizationUrl"),
		TokenURLLoc:         yamlValLoc(node, "tokenUrl"),
	}
	if scopesNode := yamlMapGet(node, "scopes"); scopesNode != nil && scopesNode.Kind == yaml.MappingNode {
		flow.Scopes = make(map[string]string)
		for i := 0; i+1 < len(scopesNode.Content); i += 2 {
			flow.Scopes[scopesNode.Content[i].Value] = scopesNode.Content[i+1].Value
		}
	}
	return flow
}

func (sp *standaloneParser) parseTags(node *yaml.Node) []Tag {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}
	var tags []Tag
	for _, item := range node.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}
		t := Tag{
			Name:        yamlStr(item, "name"),
			Description: yamlDescription(item),
			Loc:         yamlLoc(item),
			NameLoc:     yamlValLoc(item, "name"),
			Extensions:  yamlExtensions(item),
		}
		if edNode := yamlMapGet(item, "externalDocs"); edNode != nil {
			t.ExternalDocs = sp.parseExternalDocs(edNode)
		}
		tags = append(tags, t)
	}
	return tags
}

func (sp *standaloneParser) parseSecurity(node *yaml.Node) []SecurityRequirement {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil
	}
	var reqs []SecurityRequirement
	for _, item := range node.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}
		req := SecurityRequirement{Loc: yamlLoc(item)}
		for i := 0; i+1 < len(item.Content); i += 2 {
			entry := SecurityRequirementEntry{
				Name:    item.Content[i].Value,
				NameLoc: yamlLoc(item.Content[i]),
			}
			if item.Content[i+1].Kind == yaml.SequenceNode {
				for _, s := range item.Content[i+1].Content {
					entry.Scopes = append(entry.Scopes, s.Value)
				}
			}
			req.Entries = append(req.Entries, entry)
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func (sp *standaloneParser) parseExternalDocs(node *yaml.Node) *ExternalDocs {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	return &ExternalDocs{
		Description: yamlDescription(node),
		URL:         yamlStr(node, "url"),
		Loc:         yamlLoc(node),
		URLLoc:      yamlValLoc(node, "url"),
	}
}

func (sp *standaloneParser) parseExamples(node *yaml.Node) map[string]*Example {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	examples := make(map[string]*Example)
	for i := 0; i+1 < len(node.Content); i += 2 {
		name := node.Content[i].Value
		eNode := node.Content[i+1]
		if eNode.Kind != yaml.MappingNode {
			continue
		}
		ex := &Example{
			Summary:       yamlStr(eNode, "summary"),
			Description:   yamlDescription(eNode),
			ExternalValue: yamlStr(eNode, "externalValue"),
			Ref:           yamlStr(eNode, "$ref"),
			Loc:           yamlLoc(eNode),
			NameLoc:       yamlLoc(node.Content[i]),
		}
		if valNode := yamlMapGet(eNode, "value"); valNode != nil {
			ex.Value = yamlNode(valNode)
		}
		examples[name] = ex
	}
	return examples
}

func (sp *standaloneParser) collectRefs(idx *Index) {
	node := sp.docNode()
	if node == nil {
		return
	}
	sp.walkRefsYAML(node, idx, nil)
	for i := range idx.AllRefs {
		idx.AllRefs[i].URI = sp.uri
	}
	for target := range idx.Refs {
		for i := range idx.Refs[target] {
			idx.Refs[target][i].URI = sp.uri
		}
	}
}

func (sp *standaloneParser) walkRefsYAML(node *yaml.Node, idx *Index, path []string) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valNode := node.Content[i+1]
			if keyNode.Value == "$ref" && valNode.Kind == yaml.ScalarNode {
				usage := RefUsage{
					Loc:    yamlLoc(valNode),
					Target: valNode.Value,
					From:   "/" + strings.Join(path, "/"),
				}
				idx.Refs[valNode.Value] = append(idx.Refs[valNode.Value], usage)
				idx.AllRefs = append(idx.AllRefs, usage)
			} else {
				childPath := append(path, EscapeJSONPointer(keyNode.Value))
				sp.walkRefsYAML(valNode, idx, childPath)
			}
		}
	case yaml.SequenceNode:
		for i, item := range node.Content {
			childPath := append(path, strconv.Itoa(i))
			sp.walkRefsYAML(item, idx, childPath)
		}
	case yaml.DocumentNode:
		for _, child := range node.Content {
			sp.walkRefsYAML(child, idx, path)
		}
	}
}

func (sp *standaloneParser) docNode() *yaml.Node {
	if sp.root == nil {
		return nil
	}
	if sp.root.Kind == yaml.DocumentNode && len(sp.root.Content) > 0 {
		return sp.root.Content[0]
	}
	return sp.root
}

// semanticRootFromYAML builds a full SemanticNode tree from a yaml.Node so
// standalone parsing can participate in pointer-based resolution just like the
// tree-sitter path.
func semanticRootFromYAML(node *yaml.Node) *SemanticNode {
	if node == nil {
		return nil
	}
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil
		}
		return semanticRootFromYAML(node.Content[0])
	case yaml.MappingNode:
		sn := &SemanticNode{
			Kind:     NodeMapping,
			Range:    yamlLoc(node).Range,
			Children: make(map[string]*SemanticNode, len(node.Content)/2),
			Tag:      node.Tag,
			Anchor:   node.Anchor,
		}
		for i := 0; i+1 < len(node.Content); i += 2 {
			key, val := node.Content[i], node.Content[i+1]
			child := semanticRootFromYAML(val)
			if child == nil {
				child = &SemanticNode{
					Kind:   NodeNull,
					Range:  yamlLoc(val).Range,
					Tag:    val.Tag,
					Anchor: val.Anchor,
				}
			}
			child.Key = key.Value
			sn.Children[key.Value] = child
		}
		return sn
	case yaml.SequenceNode:
		sn := &SemanticNode{
			Kind:   NodeSequence,
			Range:  yamlLoc(node).Range,
			Items:  make([]*SemanticNode, 0, len(node.Content)),
			Tag:    node.Tag,
			Anchor: node.Anchor,
		}
		for _, item := range node.Content {
			child := semanticRootFromYAML(item)
			if child == nil {
				child = &SemanticNode{
					Kind:   NodeNull,
					Range:  yamlLoc(item).Range,
					Tag:    item.Tag,
					Anchor: item.Anchor,
				}
			}
			sn.Items = append(sn.Items, child)
		}
		return sn
	case yaml.ScalarNode:
		kind := NodeScalar
		var value any = node.Value
		if node.Tag == "!!null" {
			kind = NodeNull
			value = nil
		}
		return &SemanticNode{
			Kind:   kind,
			Value:  value,
			Range:  yamlLoc(node).Range,
			Tag:    node.Tag,
			Anchor: node.Anchor,
		}
	case yaml.AliasNode:
		alias := ""
		if node.Alias != nil {
			alias = node.Alias.Value
		}
		return &SemanticNode{
			Kind:   NodeScalar,
			Value:  alias,
			Range:  yamlLoc(node).Range,
			Tag:    node.Tag,
			Anchor: node.Anchor,
			Alias:  alias,
		}
	default:
		return &SemanticNode{
			Kind:   NodeScalar,
			Value:  node.Value,
			Range:  yamlLoc(node).Range,
			Tag:    node.Tag,
			Anchor: node.Anchor,
		}
	}
}

// --- yaml.Node helpers ---

func yamlStr(node *yaml.Node, key string) string {
	v := yamlMapGet(node, key)
	if v == nil {
		return ""
	}
	return v.Value
}

func yamlMapGet(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func yamlMapKeys(node *yaml.Node) []string {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	keys := make([]string, 0, len(node.Content)/2)
	for i := 0; i+1 < len(node.Content); i += 2 {
		keys = append(keys, node.Content[i].Value)
	}
	return keys
}

func yamlKeyLoc(node *yaml.Node, key string) Loc {
	if node == nil || node.Kind != yaml.MappingNode {
		return Loc{}
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return yamlLoc(node.Content[i])
		}
	}
	return Loc{}
}

func yamlValLoc(node *yaml.Node, key string) Loc {
	v := yamlMapGet(node, key)
	if v == nil {
		return Loc{}
	}
	return yamlLoc(v)
}

func yamlLoc(node *yaml.Node) Loc {
	if node == nil {
		return Loc{}
	}
	startLine := node.Line - 1
	startCol := node.Column - 1
	if startLine < 0 {
		startLine = 0
	}
	if startCol < 0 {
		startCol = 0
	}
	endLine := startLine
	endCol := startCol + utf16Len(node.Value)
	return Loc{
		Range: Range{
			Start: Position{Line: uint32(startLine), Character: uint32(startCol)},
			End:   Position{Line: uint32(endLine), Character: uint32(endCol)},
		},
	}
}

func yamlDescription(node *yaml.Node) DescriptionValue {
	v := yamlMapGet(node, "description")
	if v == nil {
		return DescriptionValue{}
	}
	dv := DescriptionValue{
		Text: v.Value,
		Loc:  yamlLoc(v),
	}
	if v.Style == yaml.LiteralStyle || v.Style == yaml.FoldedStyle {
		dv.LineOffset = 1
	}
	return dv
}

func yamlStringSlice(node *yaml.Node, key string) []string {
	v := yamlMapGet(node, key)
	if v == nil || v.Kind != yaml.SequenceNode {
		return nil
	}
	var ss []string
	for _, item := range v.Content {
		ss = append(ss, item.Value)
	}
	return ss
}

func yamlExtensions(node *yaml.Node) map[string]*Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	var exts map[string]*Node
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if strings.HasPrefix(key, "x-") {
			if exts == nil {
				exts = make(map[string]*Node)
			}
			exts[key] = yamlNode(node.Content[i+1])
		}
	}
	return exts
}

func yamlNode(node *yaml.Node) *Node {
	if node == nil {
		return nil
	}
	return &Node{
		Value: node.Value,
		Loc:   yamlLoc(node),
	}
}

func yamlIntPtr(node *yaml.Node, key string) *int {
	v := yamlStr(node, key)
	if v == "" {
		return nil
	}
	n := 0
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			if ch == '-' && n == 0 {
				continue
			}
			return nil
		}
		n = n*10 + int(ch-'0')
	}
	if len(v) > 0 && v[0] == '-' {
		n = -n
	}
	return &n
}

func yamlFloatPtr(node *yaml.Node, key string) *float64 {
	v := yamlStr(node, key)
	if v == "" {
		return nil
	}
	var f float64
	var neg bool
	var decimalFound bool
	var decimalDiv float64 = 1
	for _, ch := range v {
		if ch == '-' && f == 0 && !neg {
			neg = true
			continue
		}
		if ch == '.' && !decimalFound {
			decimalFound = true
			continue
		}
		if ch < '0' || ch > '9' {
			return nil
		}
		if decimalFound {
			decimalDiv *= 10
			f += float64(ch-'0') / decimalDiv
		} else {
			f = f*10 + float64(ch-'0')
		}
	}
	if neg {
		f = -f
	}
	return &f
}

func utf16Len(s string) int {
	n := 0
	for _, r := range s {
		if r >= 0x10000 {
			n += 2
		} else {
			n++
		}
	}
	if n == 0 {
		n = utf8.RuneCountInString(s)
	}
	return n
}
