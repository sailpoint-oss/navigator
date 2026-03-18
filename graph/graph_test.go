package graph

import (
	"sort"
	"testing"

	navigator "github.com/sailpoint-oss/navigator"
)

func TestWorkspaceGraph_AddRemoveSource(t *testing.T) {
	cache := navigator.NewIndexCache()
	fg := navigator.NewFileGraph()
	g := New(cache, fg)

	src := navigator.NewMemorySource("file:///a.yaml", []byte("openapi: '3.0.0'"))
	g.AddSource(src)

	nodes := g.AllNodes()
	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}

	g.RemoveSource("file:///a.yaml")
	nodes = g.AllNodes()
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes after remove, got %d", len(nodes))
	}
}

func TestWorkspaceGraph_Roots(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	g.SetRoot("file:///a.yaml")
	g.SetRoot("file:///b.yaml")

	roots := g.Roots()
	sort.Strings(roots)
	if len(roots) != 2 {
		t.Errorf("expected 2 roots, got %d", len(roots))
	}

	g.UnsetRoot("file:///a.yaml")
	roots = g.Roots()
	if len(roots) != 1 {
		t.Errorf("expected 1 root after unset, got %d", len(roots))
	}
}

func TestWorkspaceGraph_Invalidate(t *testing.T) {
	cache := navigator.NewIndexCache()
	fg := navigator.NewFileGraph()
	g := New(cache, fg)

	fg.AddEdge(navigator.RefEdge{FromURI: "a", ToURI: "b"})
	fg.AddEdge(navigator.RefEdge{FromURI: "a", ToURI: "c"})
	fg.AddEdge(navigator.RefEdge{FromURI: "b", ToURI: "d"})
	fg.AddEdge(navigator.RefEdge{FromURI: "c", ToURI: "d"})

	for _, uri := range []string{"a", "b", "c", "d"} {
		g.AddSource(navigator.NewMemorySource(uri, nil))
	}

	affected := g.Invalidate("d")
	sort.Strings(affected)
	// d itself + a, b, c should all be affected
	if len(affected) < 4 {
		t.Errorf("expected at least 4 affected URIs, got %d: %v", len(affected), affected)
	}
}

func TestWorkspaceGraph_Edges(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	g.AddEdge(Edge{SourceURI: "a", TargetURI: "b", Kind: EdgeRef})
	g.AddEdge(Edge{SourceURI: "a", TargetURI: "c", Kind: EdgeExternal})

	from := g.EdgesFrom("a")
	if len(from) != 2 {
		t.Errorf("expected 2 edges from a, got %d", len(from))
	}

	to := g.EdgesTo("b")
	if len(to) != 1 {
		t.Errorf("expected 1 edge to b, got %d", len(to))
	}

	g.RemoveEdgesFrom("a")
	from = g.EdgesFrom("a")
	if len(from) != 0 {
		t.Errorf("expected 0 edges after removal, got %d", len(from))
	}
	to = g.EdgesTo("b")
	if len(to) != 0 {
		t.Errorf("expected 0 reverse edges after removal, got %d", len(to))
	}
}

func TestWorkspaceGraph_StageResults(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	src := navigator.NewMemorySource("file:///test.yaml", nil)
	g.AddSource(src)

	g.SetStageResult("file:///test.yaml", &StageResult{Stage: StageParse, Version: 1})

	result := g.StageResult("file:///test.yaml", StageParse)
	if result == nil {
		t.Fatal("stage result should exist")
	}
	if result.Version != 1 {
		t.Errorf("version = %d", result.Version)
	}

	if !g.IsStageDirty("file:///test.yaml", StageRaw) {
		t.Error("raw stage should be dirty initially")
	}
	if g.IsStageDirty("file:///test.yaml", StageParse) {
		t.Error("parse stage should not be dirty after SetStageResult")
	}
}

func TestWorkspaceGraph_Resolve(t *testing.T) {
	cache := navigator.NewIndexCache()
	idx := navigator.ParseURI("file:///api.yaml", []byte(`openapi: "3.0.3"
info:
  title: Test
  version: "1.0"
paths: {}
components:
  schemas:
    Pet:
      type: object`))
	cache.Set("file:///api.yaml", idx)

	g := New(cache, navigator.NewFileGraph())
	result, err := g.Resolve("file:///api.yaml", "#/components/schemas/Pet")
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("resolve should succeed")
	}
}
