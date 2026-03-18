package navigator

import (
	"fmt"
	"strings"
)

// CrossFileResolver resolves $refs across multiple files using a workspace-wide
// IndexCache and FileGraph. It is cycle-safe and fragment-aware.
type CrossFileResolver struct {
	cache *IndexCache
	graph *FileGraph
}

// NewCrossFileResolver creates a resolver wrapping the given cache and graph.
func NewCrossFileResolver(cache *IndexCache, graph *FileGraph) *CrossFileResolver {
	return &CrossFileResolver{cache: cache, graph: graph}
}

// ResolveResult holds the outcome of a cross-file $ref resolution.
type ResolveResult struct {
	TargetURI   string
	TargetIndex *Index
	Value       interface{}
	Pointer     string
}

// Resolve resolves a $ref from one file to its target across the workspace.
// Handles local refs (#/...), external refs (./file.yaml#/...), and whole-file
// refs (./file.yaml with no fragment).
func (r *CrossFileResolver) Resolve(fromURI, ref string) (*ResolveResult, error) {
	return r.resolveWithVisited(fromURI, ref, make(map[string]bool))
}

func (r *CrossFileResolver) resolveWithVisited(fromURI, ref string, visited map[string]bool) (*ResolveResult, error) {
	if ref == "" {
		return nil, fmt.Errorf("crossfile: empty ref")
	}

	// Local ref: resolve within the source file's index.
	if strings.HasPrefix(ref, "#") {
		idx := r.cache.Get(fromURI)
		if idx == nil {
			return nil, fmt.Errorf("crossfile: source index not found for %q", fromURI)
		}
		val, err := idx.ResolveRef(ref)
		if err != nil {
			return nil, err
		}
		return &ResolveResult{
			TargetURI:   fromURI,
			TargetIndex: idx,
			Value:       val,
			Pointer:     ref,
		}, nil
	}

	filePart, fragment := SplitRefURI(ref)
	targetURI := ResolveRelativeURI(fromURI, filePart)

	visitKey := targetURI + fragment
	if visited[visitKey] {
		return nil, fmt.Errorf("crossfile: cycle detected resolving %q from %q", ref, fromURI)
	}
	visited[visitKey] = true

	targetIdx := r.cache.Get(targetURI)
	if targetIdx == nil {
		return nil, fmt.Errorf("crossfile: target file not found: %q", targetURI)
	}

	// Whole-file ref (no fragment): return the document or its primary content.
	if fragment == "" || fragment == "#" {
		return &ResolveResult{
			TargetURI:   targetURI,
			TargetIndex: targetIdx,
			Value:       targetIdx.Document,
			Pointer:     "#",
		}, nil
	}

	val, err := targetIdx.ResolveRef(fragment)
	if err != nil {
		return nil, fmt.Errorf("crossfile: resolve %q in %q: %w", fragment, targetURI, err)
	}
	return &ResolveResult{
		TargetURI:   targetURI,
		TargetIndex: targetIdx,
		Value:       val,
		Pointer:     fragment,
	}, nil
}

// CanResolve returns true if the ref can be resolved without error.
func (r *CrossFileResolver) CanResolve(fromURI, ref string) bool {
	_, err := r.Resolve(fromURI, ref)
	return err == nil
}
