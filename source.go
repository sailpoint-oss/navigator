package navigator

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchEvent classifies the type of filesystem change observed.
type WatchEvent int

const (
	WatchCreate WatchEvent = iota
	WatchModify
	WatchDelete
)

// ClassificationHint provides optional metadata that helps classify a document.
type ClassificationHint struct {
	IsOpenAPI      bool
	IsFragment     bool
	OpenAPIVersion string
	LanguageID     string
	Skip           bool
}

// DocumentSource is the abstraction for how documents enter the workspace.
type DocumentSource interface {
	URI() string
	Read(ctx context.Context) (content []byte, version int64, err error)
	Watch(ctx context.Context, onChange func(uri string, event WatchEvent)) (cancel func())
	Hint() ClassificationHint
}

// FilesystemSource reads a document from disk with version derived from mtime.
type FilesystemSource struct {
	uri  string
	path string
	hint ClassificationHint
}

// NewFilesystemSource creates a source backed by a file on disk.
func NewFilesystemSource(uri, path string) *FilesystemSource {
	return &FilesystemSource{uri: uri, path: path}
}

// NewFilesystemSourceWithHint creates a filesystem source with a classification hint.
func NewFilesystemSourceWithHint(path string, hint ClassificationHint) *FilesystemSource {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return &FilesystemSource{
		uri:  PathToURI(abs),
		path: abs,
		hint: hint,
	}
}

func (s *FilesystemSource) URI() string { return s.uri }

func (s *FilesystemSource) Read(_ context.Context) ([]byte, int64, error) {
	info, err := os.Stat(s.path)
	if err != nil {
		return nil, 0, err
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, 0, err
	}
	return data, info.ModTime().UnixNano(), nil
}

func (s *FilesystemSource) Watch(ctx context.Context, onChange func(uri string, event WatchEvent)) func() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ctx, cancel := context.WithCancel(ctx)
		go func() { <-ctx.Done() }()
		return cancel
	}

	dir := filepath.Dir(s.path)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "node_modules" || base == ".git" || base == "vendor" {
				return filepath.SkipDir
			}
			_ = watcher.Add(path)
		}
		return nil
	})

	ctx, cancel := context.WithCancel(ctx)
	go watchLoop(ctx, watcher, onChange)

	return func() {
		cancel()
		watcher.Close()
	}
}

func isRelevantExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml" || ext == ".json"
}

func watchLoop(ctx context.Context, watcher *fsnotify.Watcher, onChange func(uri string, event WatchEvent)) {
	const debounce = 100 * time.Millisecond
	pending := make(map[string]WatchEvent)
	timer := time.NewTimer(debounce)
	timer.Stop()

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case ev, ok := <-watcher.Events:
			if !ok {
				return
			}
			handleFSEvent(watcher, ev, pending)
			timer.Reset(debounce)
		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		case <-timer.C:
			for path, evt := range pending {
				onChange(PathToURI(path), evt)
			}
			pending = make(map[string]WatchEvent)
		}
	}
}

func handleFSEvent(watcher *fsnotify.Watcher, ev fsnotify.Event, pending map[string]WatchEvent) {
	path := ev.Name

	if ev.Has(fsnotify.Create) {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			_ = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if fi.IsDir() {
					base := filepath.Base(p)
					if base == "node_modules" || base == ".git" || base == "vendor" {
						return filepath.SkipDir
					}
					_ = watcher.Add(p)
					return nil
				}
				if isRelevantExt(p) {
					pending[p] = WatchCreate
				}
				return nil
			})
			return
		}
	}

	if !isRelevantExt(path) {
		return
	}

	switch {
	case ev.Has(fsnotify.Create):
		pending[path] = WatchCreate
	case ev.Has(fsnotify.Write):
		if _, already := pending[path]; !already {
			pending[path] = WatchModify
		}
	case ev.Has(fsnotify.Remove) || ev.Has(fsnotify.Rename):
		pending[path] = WatchDelete
	}
}

func (s *FilesystemSource) Hint() ClassificationHint { return s.hint }

// MemorySource holds document content in memory with manual version control.
type MemorySource struct {
	mu       sync.RWMutex
	uri      string
	content  []byte
	version  int64
	hint     ClassificationHint
	watchers []func(uri string, event WatchEvent)
}

// NewMemorySource creates a source backed by in-memory content.
func NewMemorySource(uri string, content []byte) *MemorySource {
	return &MemorySource{
		uri:     uri,
		content: content,
		version: 1,
	}
}

// NewSyntheticSource creates a source backed by in-memory content with a
// classification hint. This is the name telescope uses for programmatic sources.
func NewSyntheticSource(uri string, content []byte, hint ClassificationHint) *MemorySource {
	return &MemorySource{
		uri:     uri,
		content: append([]byte(nil), content...),
		version: 1,
		hint:    hint,
	}
}

// SyntheticSource is a type alias for MemorySource, matching telescope's naming.
type SyntheticSource = MemorySource

func (s *MemorySource) URI() string { return s.uri }

func (s *MemorySource) Read(_ context.Context) ([]byte, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c := make([]byte, len(s.content))
	copy(c, s.content)
	return c, s.version, nil
}

func (s *MemorySource) Watch(_ context.Context, onChange func(uri string, event WatchEvent)) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.watchers = append(s.watchers, onChange)
	return func() {}
}

func (s *MemorySource) Hint() ClassificationHint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hint
}

// Update replaces the content and increments the version, notifying all watchers.
func (s *MemorySource) Update(content []byte) {
	s.mu.Lock()
	s.content = content
	s.version++
	watchers := make([]func(uri string, event WatchEvent), len(s.watchers))
	copy(watchers, s.watchers)
	uri := s.uri
	s.mu.Unlock()

	for _, fn := range watchers {
		fn(uri, WatchModify)
	}
}

// FilesystemSourceFromPath creates a FilesystemSource, deriving the URI from the path.
func FilesystemSourceFromPath(path string) *FilesystemSource {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return &FilesystemSource{
		uri:  PathToURI(abs),
		path: abs,
	}
}
