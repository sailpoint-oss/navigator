package navigator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkParse_Small(b *testing.B) {
	content := []byte(`openapi: "3.0.3"
info:
  title: Small
  version: "1.0"
paths:
  /hello:
    get:
      operationId: hello
      responses:
        "200":
          description: OK`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(content)
	}
}

func BenchmarkParse_Medium(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testdata", "single", "minimal.yaml"))
	if err != nil {
		b.Skip("fixture not available")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(content)
	}
}

func BenchmarkParseContent_Medium(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testdata", "single", "minimal.yaml"))
	if err != nil {
		b.Skip("fixture not available")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseContent(content, "file:///bench.yaml")
	}
}

func BenchmarkResolveRef(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testdata", "single", "minimal.yaml"))
	if err != nil {
		b.Skip("fixture not available")
	}
	idx := Parse(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.ResolveRef("#/components/schemas/Pet")
	}
}

func BenchmarkIndexCache_ConcurrentAccess(b *testing.B) {
	cache := NewIndexCache()
	idx := &Index{Document: &Document{Version: "1"}}
	cache.Set("file:///bench.yaml", idx)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("file:///bench.yaml")
		}
	})
}

func BenchmarkFileGraph_TransitiveDependents(b *testing.B) {
	g := NewFileGraph()
	for i := 0; i < 1000; i++ {
		from := fmt.Sprintf("file:///%d.yaml", i)
		to := fmt.Sprintf("file:///%d.yaml", i+1)
		g.AddEdge(RefEdge{FromURI: from, ToURI: to})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.TransitiveDependentsOf("file:///999.yaml")
	}
}
