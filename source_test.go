package navigator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFilesystemSource_Read(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, []byte("openapi: '3.0.0'"), 0644); err != nil {
		t.Fatal(err)
	}

	src := NewFilesystemSource("file:///test.yaml", path)
	content, version, err := src.Read(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "openapi: '3.0.0'" {
		t.Errorf("content = %q", string(content))
	}
	if version == 0 {
		t.Error("version should be non-zero")
	}
}

func TestFilesystemSource_URI(t *testing.T) {
	src := NewFilesystemSource("file:///api.yaml", "/api.yaml")
	if src.URI() != "file:///api.yaml" {
		t.Errorf("URI = %q", src.URI())
	}
}

func TestMemorySource_ReadAndUpdate(t *testing.T) {
	src := NewMemorySource("file:///mem.yaml", []byte("v1"))

	content, version, err := src.Read(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "v1" || version != 1 {
		t.Errorf("initial: content=%q version=%d", string(content), version)
	}

	src.Update([]byte("v2"))
	content, version, err = src.Read(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "v2" || version != 2 {
		t.Errorf("updated: content=%q version=%d", string(content), version)
	}
}

func TestMemorySource_Watch(t *testing.T) {
	src := NewMemorySource("file:///watch.yaml", []byte("initial"))

	called := make(chan struct{}, 1)
	src.Watch(context.Background(), func() {
		called <- struct{}{}
	})

	src.Update([]byte("changed"))

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Error("watcher not called")
	}
}
