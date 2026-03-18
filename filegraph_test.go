package navigator

import (
	"sort"
	"testing"
)

func TestFileGraph_AddEdge(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "b"})
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "c"})
	g.AddEdge(RefEdge{FromURI: "b", ToURI: "c"})

	deps := g.DependenciesOf("a")
	sort.Strings(deps)
	if len(deps) != 2 || deps[0] != "b" || deps[1] != "c" {
		t.Errorf("DependenciesOf(a) = %v", deps)
	}

	rdeps := g.DependentsOf("c")
	sort.Strings(rdeps)
	if len(rdeps) != 2 || rdeps[0] != "a" || rdeps[1] != "b" {
		t.Errorf("DependentsOf(c) = %v", rdeps)
	}
}

func TestFileGraph_Diamond(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "A", ToURI: "B"})
	g.AddEdge(RefEdge{FromURI: "A", ToURI: "C"})
	g.AddEdge(RefEdge{FromURI: "B", ToURI: "D"})
	g.AddEdge(RefEdge{FromURI: "C", ToURI: "D"})

	dependents := g.TransitiveDependentsOf("D")
	sort.Strings(dependents)
	if len(dependents) != 3 {
		t.Fatalf("expected 3 transitive dependents of D, got %d: %v", len(dependents), dependents)
	}
	if dependents[0] != "A" || dependents[1] != "B" || dependents[2] != "C" {
		t.Errorf("unexpected dependents: %v", dependents)
	}
}

func TestFileGraph_Cycle(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "b"})
	g.AddEdge(RefEdge{FromURI: "b", ToURI: "a"})

	deps := g.TransitiveDependenciesOf("a")
	if len(deps) != 1 || deps[0] != "b" {
		t.Errorf("TransitiveDependenciesOf(a) = %v, want [b]", deps)
	}

	if !g.HasCycle() {
		t.Error("should detect cycle")
	}
}

func TestFileGraph_NoCycle(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "b"})
	g.AddEdge(RefEdge{FromURI: "b", ToURI: "c"})

	if g.HasCycle() {
		t.Error("should not detect cycle in DAG")
	}
}

func TestFileGraph_RemoveEdgesFrom(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "b"})
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "c"})
	g.AddEdge(RefEdge{FromURI: "b", ToURI: "c"})

	g.RemoveEdgesFrom("a")

	if deps := g.DependenciesOf("a"); len(deps) != 0 {
		t.Errorf("a should have no deps after removal, got %v", deps)
	}
	rdeps := g.DependentsOf("c")
	if len(rdeps) != 1 || rdeps[0] != "b" {
		t.Errorf("c should only have b as dependent, got %v", rdeps)
	}
}

func TestFileGraph_AllURIs(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "b"})
	g.AddEdge(RefEdge{FromURI: "b", ToURI: "c"})

	uris := g.AllURIs()
	sort.Strings(uris)
	if len(uris) != 3 || uris[0] != "a" || uris[1] != "b" || uris[2] != "c" {
		t.Errorf("AllURIs = %v", uris)
	}
}

func TestFileGraph_CyclesInvolving(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "b"})
	g.AddEdge(RefEdge{FromURI: "b", ToURI: "a"})

	cycles := g.CyclesInvolving("a")
	if len(cycles) == 0 {
		t.Error("should find at least one cycle involving a")
	}
}

func TestFileGraph_TransitiveDependencies(t *testing.T) {
	g := NewFileGraph()
	g.AddEdge(RefEdge{FromURI: "a", ToURI: "b"})
	g.AddEdge(RefEdge{FromURI: "b", ToURI: "c"})
	g.AddEdge(RefEdge{FromURI: "c", ToURI: "d"})

	deps := g.TransitiveDependenciesOf("a")
	sort.Strings(deps)
	if len(deps) != 3 || deps[0] != "b" || deps[1] != "c" || deps[2] != "d" {
		t.Errorf("TransitiveDependenciesOf(a) = %v", deps)
	}
}
