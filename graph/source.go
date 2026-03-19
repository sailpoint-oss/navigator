package graph

import (
	"context"
	"fmt"

	navigator "github.com/sailpoint-oss/navigator"
)

// LSPSource adapts an external LSP document provider into navigator's
// DocumentSource interface. Telescope uses this to bridge gossip's
// document overlay into the workspace graph.
type LSPSource struct {
	uri      string
	provider LSPDocumentProvider
	hint     navigator.ClassificationHint
}

// LSPDocumentProvider abstracts gossip's document.Store for reading document
// content without importing gossip in the core package. The signature matches
// telescope's gossip-style overlay API.
type LSPDocumentProvider interface {
	Content(uri string) (text string, version int32, ok bool)
}

// NewLSPSource creates a source backed by an LSP document provider.
func NewLSPSource(uri string, provider LSPDocumentProvider, hint ...navigator.ClassificationHint) *LSPSource {
	s := &LSPSource{uri: uri, provider: provider}
	if len(hint) > 0 {
		s.hint = hint[0]
	}
	return s
}

func (s *LSPSource) URI() string { return s.uri }

func (s *LSPSource) Read(_ context.Context) ([]byte, int64, error) {
	text, version, ok := s.provider.Content(s.uri)
	if !ok {
		return nil, 0, fmt.Errorf("document not found: %s", s.uri)
	}
	return []byte(text), int64(version), nil
}

func (s *LSPSource) Watch(_ context.Context, _ func(string, navigator.WatchEvent)) func() {
	return func() {}
}

func (s *LSPSource) Hint() navigator.ClassificationHint { return s.hint }
