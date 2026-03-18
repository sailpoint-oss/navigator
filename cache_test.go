package navigator

import (
	"sync"
	"testing"
)

func TestIndexCache_SetGet(t *testing.T) {
	c := NewIndexCache()
	idx := &Index{Document: &Document{Version: "3.0.0"}}
	c.Set("file:///a.yaml", idx)
	got := c.Get("file:///a.yaml")
	if got == nil {
		t.Fatal("expected cached index")
	}
	if got.Document.Version != "3.0.0" {
		t.Errorf("version = %q", got.Document.Version)
	}
}

func TestIndexCache_LazyBuilder(t *testing.T) {
	c := NewIndexCache()
	calls := 0
	c.SetBuilder(func(uri string) *Index {
		calls++
		return &Index{Document: &Document{Version: "built"}}
	})

	idx := c.Get("file:///lazy.yaml")
	if idx == nil || idx.Document.Version != "built" {
		t.Error("builder should produce the index")
	}
	if calls != 1 {
		t.Errorf("builder called %d times, want 1", calls)
	}

	idx2 := c.Get("file:///lazy.yaml")
	if calls != 1 {
		t.Error("second Get should use cached value")
	}
	if idx2.Document.Version != "built" {
		t.Error("cached value mismatch")
	}
}

func TestIndexCache_Invalidate(t *testing.T) {
	c := NewIndexCache()
	c.Set("file:///a.yaml", &Index{Document: &Document{Version: "1"}})
	c.Invalidate("file:///a.yaml")
	if c.Get("file:///a.yaml") != nil {
		t.Error("invalidated index should be nil without builder")
	}
}

func TestIndexCache_FindByOperationID(t *testing.T) {
	c := NewIndexCache()
	idx := newEmptyIndex()
	idx.Operations["listPets"] = &OperationRef{Path: "/pets", Method: "get"}
	c.Set("file:///api.yaml", idx)

	uri, ref := c.FindByOperationID("listPets")
	if uri == "" || ref == nil {
		t.Fatal("should find listPets")
	}
	if ref.Path != "/pets" {
		t.Errorf("path = %q", ref.Path)
	}
}

func TestIndexCache_FindRefTarget(t *testing.T) {
	c := NewIndexCache()
	content := loadFixture(t, "single/minimal.yaml")
	idx := Parse(content)
	c.Set("file:///api.yaml", idx)

	uri, val := c.FindRefTarget("#/components/schemas/Pet")
	if uri == "" || val == nil {
		t.Error("should find Pet schema across cache")
	}
}

func TestIndexCache_Concurrent(t *testing.T) {
	c := NewIndexCache()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			uri := "file:///concurrent.yaml"
			c.Set(uri, &Index{Document: &Document{Version: "v1"}})
			c.Get(uri)
			c.Invalidate(uri)
		}(i)
	}
	wg.Wait()
}

func TestIndexCache_All(t *testing.T) {
	c := NewIndexCache()
	c.Set("file:///a.yaml", &Index{})
	c.Set("file:///b.yaml", &Index{})

	all := c.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
