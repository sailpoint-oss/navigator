package navigator

import "sync"

// FileGraph is a workspace-wide directed graph tracking $ref dependencies
// between files. It maintains both forward edges (what a file references)
// and reverse edges (what references a file) for efficient invalidation.
type FileGraph struct {
	mu    sync.RWMutex
	deps  map[string]map[string]bool // forward: uri -> set of uris it refs
	rdeps map[string]map[string]bool // reverse: uri -> set of uris that ref it
	edges map[string][]RefEdge       // per-uri edge details
}

// NewFileGraph creates an empty file graph.
func NewFileGraph() *FileGraph {
	return &FileGraph{
		deps:  make(map[string]map[string]bool),
		rdeps: make(map[string]map[string]bool),
		edges: make(map[string][]RefEdge),
	}
}

// AddEdge adds a directed $ref edge from one file to another.
func (g *FileGraph) AddEdge(edge RefEdge) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.deps[edge.FromURI] == nil {
		g.deps[edge.FromURI] = make(map[string]bool)
	}
	g.deps[edge.FromURI][edge.ToURI] = true

	if g.rdeps[edge.ToURI] == nil {
		g.rdeps[edge.ToURI] = make(map[string]bool)
	}
	g.rdeps[edge.ToURI][edge.FromURI] = true

	g.edges[edge.FromURI] = append(g.edges[edge.FromURI], edge)
}

// RemoveEdgesFrom removes all outgoing edges from a URI.
func (g *FileGraph) RemoveEdgesFrom(uri string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if targets, ok := g.deps[uri]; ok {
		for target := range targets {
			if g.rdeps[target] != nil {
				delete(g.rdeps[target], uri)
				if len(g.rdeps[target]) == 0 {
					delete(g.rdeps, target)
				}
			}
		}
		delete(g.deps, uri)
	}
	delete(g.edges, uri)
}

// DependenciesOf returns the direct forward dependencies of a URI.
func (g *FileGraph) DependenciesOf(uri string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return setToSlice(g.deps[uri])
}

// DependentsOf returns the direct reverse dependents of a URI.
func (g *FileGraph) DependentsOf(uri string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return setToSlice(g.rdeps[uri])
}

// TransitiveDependentsOf returns all URIs that transitively depend on the given URI.
// Uses a visited set to handle cycles safely.
func (g *FileGraph) TransitiveDependentsOf(uri string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	var result []string
	g.walkReverse(uri, visited, &result)
	return result
}

func (g *FileGraph) walkReverse(uri string, visited map[string]bool, result *[]string) {
	for dep := range g.rdeps[uri] {
		if !visited[dep] {
			visited[dep] = true
			*result = append(*result, dep)
			g.walkReverse(dep, visited, result)
		}
	}
}

// TransitiveDependenciesOf returns all URIs that the given URI transitively depends on.
// Uses a visited set to handle cycles safely.
func (g *FileGraph) TransitiveDependenciesOf(uri string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	visited[uri] = true // exclude the starting node from results
	var result []string
	g.walkForward(uri, visited, &result)
	return result
}

func (g *FileGraph) walkForward(uri string, visited map[string]bool, result *[]string) {
	for dep := range g.deps[uri] {
		if !visited[dep] {
			visited[dep] = true
			*result = append(*result, dep)
			g.walkForward(dep, visited, result)
		}
	}
}

// EdgesFrom returns the detailed edges from a URI.
func (g *FileGraph) EdgesFrom(uri string) []RefEdge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	edges := g.edges[uri]
	result := make([]RefEdge, len(edges))
	copy(result, edges)
	return result
}

// AllURIs returns all URIs known to the graph (as source or target).
func (g *FileGraph) AllURIs() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	uris := make(map[string]bool)
	for u := range g.deps {
		uris[u] = true
	}
	for u := range g.rdeps {
		uris[u] = true
	}
	return setToSlice(uris)
}

// HasCycle returns true if the graph contains at least one cycle.
func (g *FileGraph) HasCycle() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	for uri := range g.deps {
		if g.hasCycleDFS(uri, visited, inStack) {
			return true
		}
	}
	return false
}

func (g *FileGraph) hasCycleDFS(uri string, visited, inStack map[string]bool) bool {
	if inStack[uri] {
		return true
	}
	if visited[uri] {
		return false
	}
	visited[uri] = true
	inStack[uri] = true

	for dep := range g.deps[uri] {
		if g.hasCycleDFS(dep, visited, inStack) {
			return true
		}
	}

	inStack[uri] = false
	return false
}

// CyclesInvolving returns all cycles that include the given URI.
func (g *FileGraph) CyclesInvolving(uri string) [][]string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var cycles [][]string
	g.findCycles(uri, uri, []string{uri}, make(map[string]bool), &cycles)
	return cycles
}

func (g *FileGraph) findCycles(start, current string, path []string, visited map[string]bool, cycles *[][]string) {
	for dep := range g.deps[current] {
		if dep == start && len(path) > 1 {
			cycle := make([]string, len(path))
			copy(cycle, path)
			*cycles = append(*cycles, cycle)
			continue
		}
		if !visited[dep] {
			visited[dep] = true
			g.findCycles(start, dep, append(path, dep), visited, cycles)
			visited[dep] = false
		}
	}
}

func setToSlice(m map[string]bool) []string {
	if len(m) == 0 {
		return nil
	}
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}
