package navigator

import (
	"context"
	"os"
	"strings"
)

// Project is a thin facade over workspace-wide shared state, providing a
// single-root lens into the IndexCache, FileGraph, and CrossFileResolver.
type Project struct {
	rootURI  string
	cache    *IndexCache
	graph    *FileGraph
	resolver *CrossFileResolver
	disc     *Discovery
}

// ProjectOption configures project creation.
type ProjectOption func(*projectConfig)

type projectConfig struct {
	cache    *IndexCache
	graph    *FileGraph
	disc     *Discovery
	exclude  []string
}

// WithCache sets a shared IndexCache for the project.
func WithCache(c *IndexCache) ProjectOption {
	return func(cfg *projectConfig) { cfg.cache = c }
}

// WithGraph sets a shared FileGraph for the project.
func WithGraph(g *FileGraph) ProjectOption {
	return func(cfg *projectConfig) { cfg.graph = g }
}

// WithDiscovery sets a shared Discovery for the project.
func WithDiscovery(d *Discovery) ProjectOption {
	return func(cfg *projectConfig) { cfg.disc = d }
}

// WithExclude sets directory exclude patterns for discovery.
func WithExclude(patterns []string) ProjectOption {
	return func(cfg *projectConfig) { cfg.exclude = patterns }
}

// OpenProject opens a multi-file OpenAPI project rooted at the given path.
// Only the root is loaded initially; transitive deps load lazily via the
// IndexCache builder.
func OpenProject(rootPath string, opts ...ProjectOption) (*Project, error) {
	cfg := &projectConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.cache == nil {
		cfg.cache = NewIndexCache()
	}
	if cfg.graph == nil {
		cfg.graph = NewFileGraph()
	}
	if cfg.disc == nil {
		cfg.disc = NewDiscovery(cfg.exclude)
	}

	abs, err := resolveAbsPath(rootPath)
	if err != nil {
		return nil, err
	}
	rootURI := PathToURI(abs)

	cache := cfg.cache
	graph := cfg.graph
	resolver := NewCrossFileResolver(cache, graph)

	cache.SetBuilder(func(uri string) *Index {
		path := URIToPath(uri)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		idx := ParseURI(uri, data)
		if idx != nil {
			UpdateGraphFromIndex(graph, uri, idx)
		}
		return idx
	})

	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}
	rootIdx := ParseURI(rootURI, data)
	if rootIdx == nil {
		return nil, &ParseError{Path: rootPath}
	}
	cache.Set(rootURI, rootIdx)
	UpdateGraphFromIndex(graph, rootURI, rootIdx)

	return &Project{
		rootURI:  rootURI,
		cache:    cache,
		graph:    graph,
		resolver: resolver,
		disc:     cfg.disc,
	}, nil
}

// OpenProjectFromIndex builds a project from a pre-parsed root index.
func OpenProjectFromIndex(rootURI string, root *Index, opts ...ProjectOption) *Project {
	cfg := &projectConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.cache == nil {
		cfg.cache = NewIndexCache()
	}
	if cfg.graph == nil {
		cfg.graph = NewFileGraph()
	}
	if cfg.disc == nil {
		cfg.disc = NewDiscovery(cfg.exclude)
	}

	cache := cfg.cache
	graph := cfg.graph
	resolver := NewCrossFileResolver(cache, graph)

	cache.Set(rootURI, root)
	UpdateGraphFromIndex(graph, rootURI, root)

	cache.SetBuilder(func(uri string) *Index {
		path := URIToPath(uri)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		idx := ParseURI(uri, data)
		if idx != nil {
			UpdateGraphFromIndex(graph, uri, idx)
		}
		return idx
	})

	return &Project{
		rootURI:  rootURI,
		cache:    cache,
		graph:    graph,
		resolver: resolver,
		disc:     cfg.disc,
	}
}

// Index returns the root document's index.
func (p *Project) Index() *Index {
	return p.cache.Get(p.rootURI)
}

// RootURI returns the root document URI.
func (p *Project) RootURI() string {
	return p.rootURI
}

// IndexFor returns the index for any file in the project (lazy-loads on miss).
func (p *Project) IndexFor(uri string) *Index {
	return p.cache.Get(uri)
}

// Resolve resolves a $ref from one file to its target.
func (p *Project) Resolve(fromURI, ref string) (*ResolveResult, error) {
	return p.resolver.Resolve(fromURI, ref)
}

// LoadAll eagerly loads all transitive dependencies from the root.
func (p *Project) LoadAll() error {
	return p.loadTransitive(p.rootURI, make(map[string]bool))
}

func (p *Project) loadTransitive(uri string, visited map[string]bool) error {
	if visited[uri] {
		return nil
	}
	visited[uri] = true

	idx := p.cache.Get(uri)
	if idx == nil {
		return nil
	}

	targets := CollectExternalRefTargets(uri, idx)
	for _, target := range targets {
		if err := p.loadTransitive(target, visited); err != nil {
			return err
		}
	}
	return nil
}

// Graph returns the workspace-wide FileGraph for break-glass access.
func (p *Project) Graph() *FileGraph { return p.graph }

// Cache returns the workspace-wide IndexCache for break-glass access.
func (p *Project) Cache() *IndexCache { return p.cache }

// Discovery returns the workspace Discovery for break-glass access.
func (p *Project) Discovery() *Discovery { return p.disc }

// Resolver returns the CrossFileResolver for break-glass access.
func (p *Project) Resolver() *CrossFileResolver { return p.resolver }

// AllURIs returns URIs of all files in this project's transitive closure.
func (p *Project) AllURIs() []string {
	deps := p.graph.TransitiveDependenciesOf(p.rootURI)
	result := make([]string, 0, len(deps)+1)
	result = append(result, p.rootURI)
	result = append(result, deps...)
	return result
}

// ContainsFile returns true if the URI is part of this project's transitive closure.
func (p *Project) ContainsFile(uri string) bool {
	if uri == p.rootURI {
		return true
	}
	for _, dep := range p.graph.TransitiveDependenciesOf(p.rootURI) {
		if dep == uri {
			return true
		}
	}
	return false
}

// OnFileChanged invalidates the file and re-extracts its external refs.
// Returns the URIs that were transitively affected.
func (p *Project) OnFileChanged(uri string) []string {
	p.cache.Invalidate(uri)
	idx := p.cache.Get(uri) // triggers rebuild via builder
	if idx != nil {
		UpdateGraphFromIndex(p.graph, uri, idx)
	}
	return p.graph.TransitiveDependentsOf(uri)
}

// ParseError indicates that a file could not be parsed.
type ParseError struct {
	Path string
}

func (e *ParseError) Error() string {
	return "navigator: failed to parse " + e.Path
}

func resolveAbsPath(path string) (string, error) {
	if strings.HasPrefix(path, "/") {
		return path, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd + "/" + path, nil
}

// Workspace holds shared state for a multi-root workspace.
type Workspace struct {
	Cache     *IndexCache
	Graph     *FileGraph
	Discovery *Discovery
}

// WorkspaceOption configures workspace creation.
type WorkspaceOption func(*workspaceConfig)

type workspaceConfig struct {
	exclude []string
}

// WithWorkspaceExclude sets directory exclude patterns.
func WithWorkspaceExclude(patterns []string) WorkspaceOption {
	return func(cfg *workspaceConfig) { cfg.exclude = patterns }
}

// NewWorkspace creates workspace-wide shared infrastructure that multiple
// Projects or a WorkspaceGraph can share.
func NewWorkspace(opts ...WorkspaceOption) *Workspace {
	cfg := &workspaceConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	cache := NewIndexCache()
	graph := NewFileGraph()
	disc := NewDiscovery(cfg.exclude)

	cache.SetBuilder(func(uri string) *Index {
		path := URIToPath(uri)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		idx := ParseURI(uri, data)
		if idx != nil {
			UpdateGraphFromIndex(graph, uri, idx)
		}
		return idx
	})

	return &Workspace{
		Cache:     cache,
		Graph:     graph,
		Discovery: disc,
	}
}

// OpenProject creates a Project lens over this workspace for the given root.
func (w *Workspace) OpenProject(rootPath string, opts ...ProjectOption) (*Project, error) {
	opts = append(opts, WithCache(w.Cache), WithGraph(w.Graph), WithDiscovery(w.Discovery))
	return OpenProject(rootPath, opts...)
}

// contextKey suppresses linter warnings about context usage.
var _ = context.Background
