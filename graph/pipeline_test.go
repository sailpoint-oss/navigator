package graph

import (
	"context"
	"testing"

	navigator "github.com/sailpoint-oss/navigator"
)

func TestPipelineRunner_DefaultStages(t *testing.T) {
	cache := navigator.NewIndexCache()
	fg := navigator.NewFileGraph()
	g := New(cache, fg)

	content := []byte(`openapi: "3.0.3"
info:
  title: Pipeline Test
  version: "1.0"
paths:
  /test:
    get:
      operationId: getTest
      responses:
        "200":
          description: OK`)

	src := navigator.NewMemorySource("file:///pipeline.yaml", content)
	g.AddSource(src)

	runner := NewPipelineRunner(DefaultStages()...)
	err := runner.RunAll(context.Background(), g, "file:///pipeline.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// After pipeline, parse stage should have result
	result := g.StageResult("file:///pipeline.yaml", StageParse)
	if result == nil {
		t.Error("parse stage should have result")
	}

	// Index should be in cache
	idx := cache.Get("file:///pipeline.yaml")
	if idx == nil {
		t.Error("index should be cached after pipeline")
	}
	if !idx.IsOpenAPI() {
		t.Error("parsed index should be OpenAPI")
	}
}

func TestPipelineRunner_RunThrough(t *testing.T) {
	cache := navigator.NewIndexCache()
	fg := navigator.NewFileGraph()
	g := New(cache, fg)

	src := navigator.NewMemorySource("file:///test.yaml", []byte(`openapi: "3.0.0"
info:
  title: T
  version: "1"
paths: {}`))
	g.AddSource(src)

	runner := NewPipelineRunner(DefaultStages()...)
	err := runner.RunThrough(context.Background(), g, "file:///test.yaml", StageParse)
	if err != nil {
		t.Fatal(err)
	}

	if g.StageResult("file:///test.yaml", StageParse) == nil {
		t.Error("parse should have run")
	}
}

func TestPipelineRunner_CachingCleanStages(t *testing.T) {
	cache := navigator.NewIndexCache()
	fg := navigator.NewFileGraph()
	g := New(cache, fg)

	src := navigator.NewMemorySource("file:///test.yaml", []byte(`openapi: "3.0.0"
info:
  title: T
  version: "1"
paths: {}`))
	g.AddSource(src)

	runner := NewPipelineRunner(DefaultStages()...)

	// First run
	runner.RunAll(context.Background(), g, "file:///test.yaml")

	// Second run should skip (not dirty)
	err := runner.RunAll(context.Background(), g, "file:///test.yaml")
	if err != nil {
		t.Fatal(err)
	}
}
