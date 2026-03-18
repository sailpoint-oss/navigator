package navigator

import "strings"

// ExtractExternalRefs returns all RefEdges for external (non-local) $refs in the index.
func ExtractExternalRefs(sourceURI string, idx *Index) []RefEdge {
	if idx == nil {
		return nil
	}
	var edges []RefEdge
	for _, ref := range idx.AllRefs {
		if strings.HasPrefix(ref.Target, "#") {
			continue
		}
		filePart, pointer := SplitRefURI(ref.Target)
		targetURI := ResolveRelativeURI(sourceURI, filePart)
		edges = append(edges, RefEdge{
			FromURI:     sourceURI,
			FromPointer: ref.From,
			ToURI:       targetURI,
			ToPointer:   pointer,
			RefValue:    ref.Target,
		})
	}
	return edges
}

// CollectExternalRefTargets returns the set of target URIs referenced from the index.
func CollectExternalRefTargets(sourceURI string, idx *Index) []string {
	if idx == nil {
		return nil
	}
	seen := make(map[string]bool)
	var targets []string
	for _, ref := range idx.AllRefs {
		if strings.HasPrefix(ref.Target, "#") {
			continue
		}
		filePart, _ := SplitRefURI(ref.Target)
		targetURI := ResolveRelativeURI(sourceURI, filePart)
		if !seen[targetURI] {
			seen[targetURI] = true
			targets = append(targets, targetURI)
		}
	}
	return targets
}

// UpdateGraphFromIndex removes old edges from sourceURI and adds current ones.
func UpdateGraphFromIndex(g *FileGraph, sourceURI string, idx *Index) {
	g.RemoveEdgesFrom(sourceURI)
	for _, edge := range ExtractExternalRefs(sourceURI, idx) {
		g.AddEdge(edge)
	}
}
