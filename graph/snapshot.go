package graph

import (
	"sync"
	"sync/atomic"

	navigator "github.com/sailpoint-oss/navigator"
)

// Snapshot is an immutable point-in-time view of the workspace graph.
type Snapshot struct {
	ID             int64
	Nodes          map[string]SnapshotNode
	Roots          []string
	Diagnostics    map[string][]Diagnostic
	PointerIndices map[string]*navigator.PointerIndex
}

// SnapshotNode holds the state of a single document at snapshot time.
type SnapshotNode struct {
	URI          string
	Version      int64
	Raw          []byte
	StageResults map[StageName]*StageResult
}

// SnapshotManager produces immutable snapshots from a WorkspaceGraph.
type SnapshotManager struct {
	mu        sync.RWMutex
	graph     *WorkspaceGraph
	current   *Snapshot
	nextID    atomic.Int64
	callbacks []func(*Snapshot)
}

// NewSnapshotManager creates a manager that can optionally be bound to a graph.
// When called with no arguments, the graph must be passed to Build.
// When called with a graph argument, Build uses that graph.
func NewSnapshotManager(g ...*WorkspaceGraph) *SnapshotManager {
	sm := &SnapshotManager{}
	if len(g) > 0 {
		sm.graph = g[0]
	}
	sm.nextID.Store(1)
	return sm
}

// Build creates a new snapshot from the current (or given) graph state.
// Accepts an optional graph argument to match telescope's Build(graph) pattern.
func (sm *SnapshotManager) Build(g ...*WorkspaceGraph) *Snapshot {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	target := sm.graph
	if len(g) > 0 && g[0] != nil {
		target = g[0]
	}
	if target == nil {
		return nil
	}

	target.mu.RLock()
	defer target.mu.RUnlock()

	snap := &Snapshot{
		ID:             sm.nextID.Add(1) - 1,
		Nodes:          make(map[string]SnapshotNode, len(target.nodes)),
		Diagnostics:    make(map[string][]Diagnostic),
		PointerIndices: make(map[string]*navigator.PointerIndex),
	}

	for uri, node := range target.nodes {
		sn := SnapshotNode{
			URI:          uri,
			Version:      node.Version,
			Raw:          node.Raw,
			StageResults: make(map[StageName]*StageResult, len(node.StageResults)),
		}
		for stage, result := range node.StageResults {
			sn.StageResults[stage] = result
			if result.Diagnostics != nil {
				snap.Diagnostics[uri] = append(snap.Diagnostics[uri], result.Diagnostics...)
			}
			if po, ok := result.Data.(*ParseOutput); ok && po != nil && po.PointerIndex != nil {
				snap.PointerIndices[uri] = po.PointerIndex
			}
		}
		snap.Nodes[uri] = sn
	}

	for root := range target.roots {
		snap.Roots = append(snap.Roots, root)
	}

	sm.current = snap

	callbacks := make([]func(*Snapshot), len(sm.callbacks))
	copy(callbacks, sm.callbacks)

	go func() {
		for _, cb := range callbacks {
			cb(snap)
		}
	}()

	return snap
}

// Current returns the most recent snapshot.
func (sm *SnapshotManager) Current() *Snapshot {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current
}

// OnSnapshot registers a callback that fires whenever a new snapshot is built.
func (sm *SnapshotManager) OnSnapshot(fn func(*Snapshot)) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.callbacks = append(sm.callbacks, fn)
}
