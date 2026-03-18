package graph

import (
	"sync"
	"testing"
	"time"

	navigator "github.com/sailpoint-oss/navigator"
)

func TestSnapshot_Build(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	src := navigator.NewMemorySource("file:///a.yaml", []byte("test"))
	g.AddSource(src)
	g.SetRoot("file:///a.yaml")

	sm := NewSnapshotManager(g)
	snap := sm.Build()

	if snap.ID != 1 {
		t.Errorf("snap ID = %d, want 1", snap.ID)
	}
	if len(snap.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(snap.Nodes))
	}
	if len(snap.Roots) != 1 {
		t.Errorf("expected 1 root, got %d", len(snap.Roots))
	}
}

func TestSnapshot_Immutability(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	src := navigator.NewMemorySource("file:///a.yaml", []byte("v1"))
	g.AddSource(src)

	sm := NewSnapshotManager(g)
	snap1 := sm.Build()

	// Modify graph after snapshot
	g.AddSource(navigator.NewMemorySource("file:///b.yaml", []byte("v2")))

	// Snapshot should not reflect new node
	if len(snap1.Nodes) != 1 {
		t.Error("snapshot should be immutable")
	}

	snap2 := sm.Build()
	if len(snap2.Nodes) != 2 {
		t.Errorf("new snapshot should have 2 nodes, got %d", len(snap2.Nodes))
	}
}

func TestSnapshot_IDIncrement(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	sm := NewSnapshotManager(g)

	snap1 := sm.Build()
	snap2 := sm.Build()
	snap3 := sm.Build()

	if snap2.ID != snap1.ID+1 || snap3.ID != snap2.ID+1 {
		t.Errorf("IDs should increment: %d, %d, %d", snap1.ID, snap2.ID, snap3.ID)
	}
}

func TestSnapshot_Callbacks(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	sm := NewSnapshotManager(g)

	called := make(chan int64, 1)
	sm.OnSnapshot(func(s *Snapshot) {
		called <- s.ID
	})

	sm.Build()

	select {
	case id := <-called:
		if id != 1 {
			t.Errorf("callback received ID %d, want 1", id)
		}
	case <-time.After(time.Second):
		t.Error("callback not called")
	}
}

func TestSnapshot_Current(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	sm := NewSnapshotManager(g)

	if sm.Current() != nil {
		t.Error("should be nil before first build")
	}

	snap := sm.Build()
	if sm.Current() != snap {
		t.Error("Current should return last built snapshot")
	}
}

func TestSnapshot_ConcurrentReads(t *testing.T) {
	g := New(navigator.NewIndexCache(), navigator.NewFileGraph())
	for i := 0; i < 10; i++ {
		g.AddSource(navigator.NewMemorySource("file:///"+string(rune('a'+i))+".yaml", nil))
	}

	sm := NewSnapshotManager(g)
	sm.Build()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snap := sm.Current()
			if snap == nil {
				return
			}
			_ = len(snap.Nodes)
			_ = snap.Roots
		}()
	}
	wg.Wait()
}
