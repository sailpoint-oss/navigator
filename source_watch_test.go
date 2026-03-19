package navigator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type watchResult struct {
	URI   string
	Event WatchEvent
}

func TestFilesystemSource_Watch_Create(t *testing.T) {
	dir := t.TempDir()
	src := NewFilesystemSource("file:///root.yaml", filepath.Join(dir, "root.yaml"))
	src.path = filepath.Join(dir, "root.yaml")

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	results := make(chan watchResult, 10)
	cancel := src.Watch(ctx, func(uri string, event WatchEvent) {
		results <- watchResult{URI: uri, Event: event}
	})
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	path := filepath.Join(dir, "pet.yaml")
	if err := os.WriteFile(path, []byte("type: object"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-results:
		if r.Event != WatchCreate {
			t.Errorf("expected WatchCreate, got %d", r.Event)
		}
		if r.URI != PathToURI(path) {
			t.Errorf("URI = %q, want %q", r.URI, PathToURI(path))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for create event")
	}
}

func TestFilesystemSource_Watch_Modify(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "api.yaml")
	if err := os.WriteFile(path, []byte("openapi: '3.0.0'"), 0644); err != nil {
		t.Fatal(err)
	}

	src := NewFilesystemSource("file:///api.yaml", path)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	results := make(chan watchResult, 10)
	cancel := src.Watch(ctx, func(uri string, event WatchEvent) {
		results <- watchResult{URI: uri, Event: event}
	})
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	if err := os.WriteFile(path, []byte("openapi: '3.1.0'"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-results:
		if r.Event != WatchCreate && r.Event != WatchModify {
			t.Errorf("expected WatchCreate or WatchModify, got %d", r.Event)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for modify event")
	}
}

func TestFilesystemSource_Watch_Delete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "remove-me.json")
	if err := os.WriteFile(path, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	src := NewFilesystemSource("file:///remove-me.json", path)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	results := make(chan watchResult, 10)
	cancel := src.Watch(ctx, func(uri string, event WatchEvent) {
		results <- watchResult{URI: uri, Event: event}
	})
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-results:
		if r.Event != WatchDelete {
			t.Errorf("expected WatchDelete, got %d", r.Event)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for delete event")
	}
}

func TestFilesystemSource_Watch_IgnoresIrrelevantFiles(t *testing.T) {
	dir := t.TempDir()
	src := NewFilesystemSource("file:///root.yaml", filepath.Join(dir, "root.yaml"))
	src.path = filepath.Join(dir, "root.yaml")

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	results := make(chan watchResult, 10)
	cancel := src.Watch(ctx, func(uri string, event WatchEvent) {
		results <- watchResult{URI: uri, Event: event}
	})
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	if err := os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-results:
		t.Errorf("unexpected event for non-relevant file: %+v", r)
	case <-time.After(300 * time.Millisecond):
		// expected: no event
	}
}

func TestFilesystemSource_Watch_Subdirectory(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "schemas")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}

	src := NewFilesystemSource("file:///root.yaml", filepath.Join(dir, "root.yaml"))
	src.path = filepath.Join(dir, "root.yaml")

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	results := make(chan watchResult, 10)
	cancel := src.Watch(ctx, func(uri string, event WatchEvent) {
		results <- watchResult{URI: uri, Event: event}
	})
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	path := filepath.Join(sub, "pet.yaml")
	if err := os.WriteFile(path, []byte("type: object"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-results:
		if r.Event != WatchCreate {
			t.Errorf("expected WatchCreate, got %d", r.Event)
		}
		if r.URI != PathToURI(path) {
			t.Errorf("URI = %q, want %q", r.URI, PathToURI(path))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for subdirectory create event")
	}
}

func TestFilesystemSource_Watch_NewSubdirectory(t *testing.T) {
	dir := t.TempDir()
	src := NewFilesystemSource("file:///root.yaml", filepath.Join(dir, "root.yaml"))
	src.path = filepath.Join(dir, "root.yaml")

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	results := make(chan watchResult, 10)
	cancel := src.Watch(ctx, func(uri string, event WatchEvent) {
		results <- watchResult{URI: uri, Event: event}
	})
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	sub := filepath.Join(dir, "newdir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}

	time.Sleep(150 * time.Millisecond)

	path := filepath.Join(sub, "schema.yml")
	if err := os.WriteFile(path, []byte("type: string"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-results:
		if r.URI != PathToURI(path) {
			t.Errorf("URI = %q, want %q", r.URI, PathToURI(path))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for event in newly created subdirectory")
	}
}

func TestFilesystemSource_Watch_CancelStops(t *testing.T) {
	dir := t.TempDir()
	src := NewFilesystemSource("file:///root.yaml", filepath.Join(dir, "root.yaml"))
	src.path = filepath.Join(dir, "root.yaml")

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	results := make(chan watchResult, 10)
	cancel := src.Watch(ctx, func(uri string, event WatchEvent) {
		results <- watchResult{URI: uri, Event: event}
	})

	cancel()
	time.Sleep(50 * time.Millisecond)

	if err := os.WriteFile(filepath.Join(dir, "after-cancel.yaml"), []byte("x: 1"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-results:
		t.Errorf("received event after cancel: %+v", r)
	case <-time.After(300 * time.Millisecond):
		// expected: no event
	}
}
