package navigator

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"gopkg.in/yaml.v3"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func validateSyntaxTreeSitter(root *tree_sitter.Node, sink *issueSink) {
	if root == nil {
		return
	}
	if root.HasError() {
		var found bool
		var walk func(*tree_sitter.Node)
		walk = func(n *tree_sitter.Node) {
			if n == nil || !sink.canAdd() {
				return
			}
			if n.Kind() == "ERROR" {
				found = true
				sink.add(Issue{
					Code:     "syntax.parse-error",
					Message:  "syntax error in document",
					Pointer:  "",
					Range:    rangeFromTSNode(n),
					Severity: SeverityError,
					Category: CategorySyntax,
				})
			}
			for i := uint(0); i < n.ChildCount(); i++ {
				walk(n.Child(i))
			}
		}
		walk(root)
		if !found {
			sink.add(Issue{
				Code:     "syntax.parse-error",
				Message:  "syntax error in document",
				Pointer:  "",
				Range:    rangeFromTSNode(root),
				Severity: SeverityError,
				Category: CategorySyntax,
			})
		}
	}
}

func validateDuplicateKeysTreeSitter(node *tree_sitter.Node, content []byte, path string, sink *issueSink) {
	if node == nil || !sink.canAdd() {
		return
	}
	switch node.Kind() {
	case "block_mapping", "flow_mapping", "object":
		seen := make(map[string]*tree_sitter.Node)
		for i := uint(0); i < node.ChildCount(); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			ck := child.Kind()
			if ck != "block_mapping_pair" && ck != "flow_pair" && ck != "pair" {
				continue
			}
			keyNode := child.ChildByFieldName("key")
			if keyNode == nil {
				continue
			}
			keyText := unquoteTS(content, keyNode)
			if _, dup := seen[keyText]; dup {
				ptr := path + "/" + EscapeJSONPointer(keyText)
				sink.add(Issue{
					Code:     "syntax.duplicate-key",
					Message:  fmt.Sprintf("duplicate key %q", keyText),
					Pointer:  ptr,
					Range:    rangeFromTSNode(keyNode),
					Severity: SeverityError,
					Category: CategorySyntax,
				})
			} else {
				seen[keyText] = keyNode
			}
			valNode := child.ChildByFieldName("value")
			if valNode != nil {
				childPath := path + "/" + EscapeJSONPointer(keyText)
				validateDuplicateKeysTreeSitter(valNode, content, childPath, sink)
			}
		}
		return
	case "block_sequence", "flow_sequence", "array":
		idx := 0
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
						childPath := path + "/" + strconv.Itoa(idx)
						validateDuplicateKeysTreeSitter(inner, content, childPath, sink)
						idx++
					}
				}
			} else if ck == "[" || ck == "]" || ck == "," || ck == "comment" {
				continue
			} else {
				childPath := path + "/" + strconv.Itoa(idx)
				validateDuplicateKeysTreeSitter(child, content, childPath, sink)
				idx++
			}
		}
		return
	case "stream", "document":
		for i := uint(0); i < node.ChildCount(); i++ {
			c := node.Child(i)
			if c == nil {
				continue
			}
			if k := c.Kind(); k != "---" && k != "..." && k != "comment" && k != "yaml_directive" && k != "tag_directive" {
				validateDuplicateKeysTreeSitter(c, content, path, sink)
				return
			}
		}
		return
	case "block_node", "flow_node":
		for i := uint(0); i < node.ChildCount(); i++ {
			c := node.Child(i)
			if c == nil {
				continue
			}
			if k := c.Kind(); k != "tag" && k != "anchor" {
				validateDuplicateKeysTreeSitter(c, content, path, sink)
				return
			}
		}
		return
	default:
		for i := uint(0); i < node.ChildCount(); i++ {
			validateDuplicateKeysTreeSitter(node.Child(i), content, path, sink)
		}
	}
}

func unquoteTS(content []byte, n *tree_sitter.Node) string {
	if n == nil {
		return ""
	}
	start := int(n.StartByte())
	end := int(n.EndByte())
	if start < 0 || end > len(content) || start >= end {
		return ""
	}
	return unquote(string(content[start:end]))
}

func validateDuplicateKeysYAML(idx *Index, sink *issueSink) {
	if idx == nil || len(idx.Content()) == 0 {
		return
	}
	var root yaml.Node
	if err := yaml.Unmarshal(idx.Content(), &root); err != nil {
		return
	}
	doc := yamlDocNodeForValidation(&root)
	if doc == nil {
		return
	}
	walkYAMLDupKeys(doc, "", sink)
}

func yamlDocNodeForValidation(root *yaml.Node) *yaml.Node {
	if root == nil {
		return nil
	}
	if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		return root.Content[0]
	}
	return root
}

func walkYAMLDupKeys(node *yaml.Node, path string, sink *issueSink) {
	if node == nil || !sink.canAdd() {
		return
	}
	switch node.Kind {
	case yaml.MappingNode:
		seen := make(map[string]bool)
		for i := 0; i+1 < len(node.Content); i += 2 {
			kn := node.Content[i]
			vn := node.Content[i+1]
			key := kn.Value
			ptr := path + "/" + EscapeJSONPointer(key)
			if seen[key] {
				sink.add(Issue{
					Code:     "syntax.duplicate-key",
					Message:  fmt.Sprintf("duplicate key %q", key),
					Pointer:  ptr,
					Range:    rangeFromYAMLScalarNode(kn),
					Severity: SeverityError,
					Category: CategorySyntax,
				})
			}
			seen[key] = true
			walkYAMLDupKeys(vn, ptr, sink)
		}
	case yaml.SequenceNode:
		for i, item := range node.Content {
			walkYAMLDupKeys(item, path+"/"+strconv.Itoa(i), sink)
		}
	}
}

func rangeFromYAMLScalarNode(n *yaml.Node) Range {
	if n == nil {
		return Range{}
	}
	startLine := n.Line - 1
	startCol := n.Column - 1
	if startLine < 0 {
		startLine = 0
	}
	if startCol < 0 {
		startCol = 0
	}
	endLine := startLine
	endCol := startCol + utf16LenYAMLValidation(n.Value)
	return Range{
		Start: Position{Line: uint32(startLine), Character: uint32(startCol)},
		End:   Position{Line: uint32(endLine), Character: uint32(endCol)},
	}
}

func utf16LenYAMLValidation(s string) int {
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
