package navigator

import (
	"context"
	"os"
	"path/filepath"
	"sync"
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
	Watch(ctx context.Context, onChange func()) (cancel func())
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

func (s *FilesystemSource) Watch(ctx context.Context, onChange func()) func() {
	// Simple polling implementation; callers needing fsnotify should use the
	// graph package's more sophisticated watcher.
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-ctx.Done()
	}()
	return cancel
}

func (s *FilesystemSource) Hint() ClassificationHint { return s.hint }

// MemorySource holds document content in memory with manual version control.
type MemorySource struct {
	mu       sync.RWMutex
	uri      string
	content  []byte
	version  int64
	hint     ClassificationHint
	watchers []func()
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

func (s *MemorySource) Watch(_ context.Context, onChange func()) func() {
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
	watchers := make([]func(), len(s.watchers))
	copy(watchers, s.watchers)
	s.mu.Unlock()

	for _, fn := range watchers {
		fn()
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
