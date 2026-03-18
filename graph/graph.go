package graph

import (
	"sync"
	"time"

	navigator "github.com/sailpoint-oss/navigator"
)

// ChangeAction describes what kind of graph mutation occurred.
type ChangeAction string

const (
	ChangeAddSource    ChangeAction = "add_source"
	ChangeRemoveSource ChangeAction = "remove_source"
	ChangeAddEdge      ChangeAction = "add_edge"
	ChangeRemoveEdges  ChangeAction = "remove_edges"
	ChangeInvalidate   ChangeAction = "invalidate"
	ChangeSetRoot      ChangeAction = "set_root"
)

// ChangeEntry records a single graph mutation for debugging and replay.
type ChangeEntry struct {
	Time     time.Time
	Action   ChangeAction
	URI      string
	Affected []string
}

// ReadOnlyGraph provides a read-only view of the workspace graph for SDK consumers.
type ReadOnlyGraph interface {
	Node(uri string) *GraphNode
	AllNodes() []string
	Roots() []string
	EdgesFrom(uri string) []Edge
	EdgesTo(uri string) []Edge
	Dependents(uri string) []string
	Dependencies(uri string) []string
	TransitiveDependents(uri string) []string
}

// WorkspaceGraph manages the full workspace topology for LSP-grade tools.
// It wraps navigator's workspace-wide IndexCache and FileGraph, adding
// per-node pipeline state, invalidation cascading, and snapshot management.
type WorkspaceGraph struct {
	mu        sync.RWMutex
	cache     *navigator.IndexCache
	graph     *navigator.FileGraph
	nodes     map[string]*GraphNode
	edges     map[string][]Edge
	revEdges  map[string][]Edge
	roots     map[string]bool
	ChangeLog []ChangeEntry
}

// New creates a WorkspaceGraph wrapping the given shared infrastructure.
func New(cache *navigator.IndexCache, graph *navigator.FileGraph) *WorkspaceGraph {
	return &WorkspaceGraph{
		cache:    cache,
		graph:    graph,
		nodes:    make(map[string]*GraphNode),
		edges:    make(map[string][]Edge),
		revEdges: make(map[string][]Edge),
		roots:    make(map[string]bool),
	}
}

// NewWorkspaceGraph creates a standalone WorkspaceGraph with its own
// IndexCache and FileGraph. Matches telescope's zero-arg constructor.
func NewWorkspaceGraph() *WorkspaceGraph {
	return New(navigator.NewIndexCache(), navigator.NewFileGraph())
}

// Cache returns the underlying IndexCache.
func (g *WorkspaceGraph) Cache() *navigator.IndexCache { return g.cache }

// FileGraph returns the underlying FileGraph.
func (g *WorkspaceGraph) FileGraph() *navigator.FileGraph { return g.graph }

// AddSource adds a document source to the graph.
func (g *WorkspaceGraph) AddSource(src navigator.DocumentSource) {
	g.mu.Lock()
	defer g.mu.Unlock()

	uri := src.URI()
	if node, ok := g.nodes[uri]; ok {
		node.Source = src
	} else {
		g.nodes[uri] = newGraphNode(src)
	}
	g.markAllDirty(uri)
	g.ChangeLog = append(g.ChangeLog, ChangeEntry{
		Time: time.Now(), Action: ChangeAddSource, URI: uri,
	})
}

// RemoveSource removes a document source from the graph.
// Returns the URIs of nodes that were affected by the removal.
func (g *WorkspaceGraph) RemoveSource(uri string) []string {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.nodes[uri]; !ok {
		return nil
	}

	affected := g.removeEdgesFromLocked(uri)
	g.removeIncomingEdgesLocked(uri)

	delete(g.nodes, uri)
	delete(g.roots, uri)
	g.cache.Delete(uri)
	g.graph.RemoveEdgesFrom(uri)
	g.ChangeLog = append(g.ChangeLog, ChangeEntry{
		Time: time.Now(), Action: ChangeRemoveSource, URI: uri, Affected: affected,
	})
	return affected
}

func (g *WorkspaceGraph) removeEdgesFromLocked(uri string) []string {
	outgoing := g.edges[uri]
	affected := make(map[string]bool)
	for _, edge := range outgoing {
		affected[edge.TargetURI] = true
		revs := g.revEdges[edge.TargetURI]
		filtered := revs[:0]
		for _, r := range revs {
			if r.SourceURI != uri {
				filtered = append(filtered, r)
			}
		}
		if len(filtered) == 0 {
			delete(g.revEdges, edge.TargetURI)
		} else {
			g.revEdges[edge.TargetURI] = filtered
		}
	}
	delete(g.edges, uri)
	result := make([]string, 0, len(affected))
	for u := range affected {
		result = append(result, u)
	}
	return result
}

func (g *WorkspaceGraph) removeIncomingEdgesLocked(uri string) {
	incoming := g.revEdges[uri]
	for _, edge := range incoming {
		filtered := g.edges[edge.SourceURI][:0]
		for _, e := range g.edges[edge.SourceURI] {
			if e.TargetURI != uri {
				filtered = append(filtered, e)
			}
		}
		g.edges[edge.SourceURI] = filtered
	}
	delete(g.revEdges, uri)
}

// SetRoot marks or unmarks a URI as a root document. When isRoot is true the
// URI is added to the roots set; when false it is removed. This matches
// telescope's SetRoot(uri, isRoot) calling convention.
func (g *WorkspaceGraph) SetRoot(uri string, isRoot ...bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	root := true
	if len(isRoot) > 0 {
		root = isRoot[0]
	}
	if root {
		g.roots[uri] = true
	} else {
		delete(g.roots, uri)
	}
	g.ChangeLog = append(g.ChangeLog, ChangeEntry{
		Time: time.Now(), Action: ChangeSetRoot, URI: uri,
	})
}

// UnsetRoot removes a URI from the roots set.
func (g *WorkspaceGraph) UnsetRoot(uri string) {
	g.SetRoot(uri, false)
}

// Roots returns all current root URIs.
func (g *WorkspaceGraph) Roots() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]string, 0, len(g.roots))
	for r := range g.roots {
		result = append(result, r)
	}
	return result
}

// Invalidate marks a file and all transitive dependents as dirty.
// Returns the full set of affected URIs.
func (g *WorkspaceGraph) Invalidate(uri string) []string {
	g.mu.Lock()
	defer g.mu.Unlock()

	visited := make(map[string]bool)
	queue := []string{uri}
	visited[uri] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		g.markAllDirty(current)
		for _, edge := range g.revEdges[current] {
			if !visited[edge.SourceURI] {
				visited[edge.SourceURI] = true
				queue = append(queue, edge.SourceURI)
			}
		}
	}

	// Also cascade via the FileGraph reverse edges
	for _, dep := range g.graph.TransitiveDependentsOf(uri) {
		if !visited[dep] {
			visited[dep] = true
			g.markAllDirty(dep)
		}
	}

	affected := make([]string, 0, len(visited))
	for u := range visited {
		affected = append(affected, u)
	}
	g.ChangeLog = append(g.ChangeLog, ChangeEntry{
		Time: time.Now(), Action: ChangeInvalidate, URI: uri, Affected: affected,
	})
	return affected
}

// Node returns the GraphNode for a URI, or nil.
func (g *WorkspaceGraph) Node(uri string) *GraphNode {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.nodes[uri]
}

// AllNodes returns all URIs with nodes.
func (g *WorkspaceGraph) AllNodes() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]string, 0, len(g.nodes))
	for uri := range g.nodes {
		result = append(result, uri)
	}
	return result
}

// AddEdge adds a graph-level edge (separate from FileGraph edges).
func (g *WorkspaceGraph) AddEdge(e Edge) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges[e.SourceURI] = append(g.edges[e.SourceURI], e)
	g.revEdges[e.TargetURI] = append(g.revEdges[e.TargetURI], e)
	g.ChangeLog = append(g.ChangeLog, ChangeEntry{
		Time: time.Now(), Action: ChangeAddEdge, URI: e.SourceURI,
		Affected: []string{e.TargetURI},
	})
}

// RemoveEdgesFrom removes all outgoing graph-level edges from a URI.
// Returns the URIs of previously referenced targets.
func (g *WorkspaceGraph) RemoveEdgesFrom(uri string) []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	affected := g.removeEdgesFromLocked(uri)
	g.ChangeLog = append(g.ChangeLog, ChangeEntry{
		Time: time.Now(), Action: ChangeRemoveEdges, URI: uri, Affected: affected,
	})
	return affected
}

// EdgesFrom returns outgoing edges from a URI.
func (g *WorkspaceGraph) EdgesFrom(uri string) []Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]Edge, len(g.edges[uri]))
	copy(result, g.edges[uri])
	return result
}

// EdgesTo returns incoming edges to a URI.
func (g *WorkspaceGraph) EdgesTo(uri string) []Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]Edge, len(g.revEdges[uri]))
	copy(result, g.revEdges[uri])
	return result
}

// Resolve provides workspace-wide cross-file $ref resolution.
func (g *WorkspaceGraph) Resolve(fromURI, ref string) (*navigator.ResolveResult, error) {
	resolver := navigator.NewCrossFileResolver(g.cache, g.graph)
	return resolver.Resolve(fromURI, ref)
}

// SetStageResult stores a pipeline stage result for a node.
func (g *WorkspaceGraph) SetStageResult(uri string, result *StageResult) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if node, ok := g.nodes[uri]; ok {
		node.StageResults[result.Stage] = result
		delete(node.DirtyStages, result.Stage)
	}
}

// StageResult returns the cached stage result for a node, or nil.
func (g *WorkspaceGraph) StageResult(uri string, stage StageName) *StageResult {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if node, ok := g.nodes[uri]; ok {
		return node.StageResults[stage]
	}
	return nil
}

// IsStageDirty returns true if the stage needs to be re-run.
func (g *WorkspaceGraph) IsStageDirty(uri string, stage StageName) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if node, ok := g.nodes[uri]; ok {
		return node.DirtyStages[stage]
	}
	return true
}

// Dependents returns URIs that reference the given URI via any edge kind.
func (g *WorkspaceGraph) Dependents(uri string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	seen := make(map[string]bool)
	for _, e := range g.revEdges[uri] {
		seen[e.SourceURI] = true
	}
	return mapToSlice(seen)
}

// Dependencies returns URIs that the given URI references via any edge kind.
func (g *WorkspaceGraph) Dependencies(uri string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	seen := make(map[string]bool)
	for _, e := range g.edges[uri] {
		seen[e.TargetURI] = true
	}
	return mapToSlice(seen)
}

// TransitiveDependents returns all URIs reachable by walking reverse edges
// from the given URI (breadth-first). The starting URI is not included.
func (g *WorkspaceGraph) TransitiveDependents(uri string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := map[string]bool{uri: true}
	queue := []string{uri}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, e := range g.revEdges[current] {
			if !visited[e.SourceURI] {
				visited[e.SourceURI] = true
				queue = append(queue, e.SourceURI)
			}
		}
	}

	delete(visited, uri)
	return mapToSlice(visited)
}

// DetectCycles returns all URIs that participate in reference cycles.
func (g *WorkspaceGraph) DetectCycles() [][]string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	const (
		white = iota
		gray
		black
	)
	color := make(map[string]int)
	var cycles [][]string

	var dfs func(uri string, path []string)
	dfs = func(uri string, path []string) {
		color[uri] = gray
		path = append(path, uri)
		for _, e := range g.edges[uri] {
			target := e.TargetURI
			switch color[target] {
			case white:
				dfs(target, path)
			case gray:
				start := -1
				for i, u := range path {
					if u == target {
						start = i
						break
					}
				}
				if start >= 0 {
					cycle := make([]string, len(path)-start)
					copy(cycle, path[start:])
					cycles = append(cycles, cycle)
				}
			}
		}
		color[uri] = black
	}

	for uri := range g.nodes {
		if color[uri] == white {
			dfs(uri, nil)
		}
	}
	return cycles
}

func mapToSlice(m map[string]bool) []string {
	if len(m) == 0 {
		return nil
	}
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func (g *WorkspaceGraph) markAllDirty(uri string) {
	if node, ok := g.nodes[uri]; ok {
		for stage := range node.StageResults {
			node.DirtyStages[stage] = true
		}
		node.DirtyStages[StageRaw] = true
		node.DirtyStages[StageParse] = true
		node.DirtyStages[StageLint] = true
		node.DirtyStages[StageBind] = true
		node.DirtyStages[StageValidate] = true
		node.DirtyStages[StageAnalyze] = true
	}
}
