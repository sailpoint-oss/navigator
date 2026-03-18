package navigator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileRole classifies a discovered file.
type FileRole int

const (
	RoleUnknown  FileRole = iota
	RoleRoot              // A root OpenAPI document
	RoleFragment          // A partial document / fragment
)

// DiscoveredFile represents a file found during workspace scanning.
type DiscoveredFile struct {
	Path  string
	URI   string
	Role  FileRole
	MTime int64
}

// Discovery scans a workspace for OpenAPI files and classifies them.
type Discovery struct {
	mu      sync.RWMutex
	files   map[string]*DiscoveredFile
	roots   []string
	exclude []string
}

// NewDiscovery creates a new Discovery with the given exclude patterns.
func NewDiscovery(exclude []string) *Discovery {
	return &Discovery{
		files:   make(map[string]*DiscoveredFile),
		exclude: exclude,
	}
}

// Scan walks the workspace root and discovers all OpenAPI files.
func (d *Discovery) Scan(workspaceRoot string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.files = make(map[string]*DiscoveredFile)
	d.roots = nil

	return filepath.Walk(workspaceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "node_modules" || base == ".git" || base == "vendor" {
				return filepath.SkipDir
			}
			for _, ex := range d.exclude {
				if matched, _ := filepath.Match(ex, base); matched {
					return filepath.SkipDir
				}
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		role := classifyFromDisk(path)
		uri := PathToURI(path)
		df := &DiscoveredFile{
			Path:  path,
			URI:   uri,
			Role:  role,
			MTime: info.ModTime().UnixNano(),
		}
		d.files[uri] = df
		if role == RoleRoot {
			d.roots = append(d.roots, uri)
		}
		return nil
	})
}

// Roots returns all discovered root document URIs.
func (d *Discovery) Roots() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]string, len(d.roots))
	copy(result, d.roots)
	return result
}

// AllFiles returns all discovered files.
func (d *Discovery) AllFiles() []*DiscoveredFile {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]*DiscoveredFile, 0, len(d.files))
	for _, f := range d.files {
		result = append(result, f)
	}
	return result
}

// FileByURI returns the discovered file for a URI, if any.
func (d *Discovery) FileByURI(uri string) *DiscoveredFile {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.files[NormalizeURI(uri)]
}

// FileByPath returns the discovered file for a path, if any.
func (d *Discovery) FileByPath(path string) *DiscoveredFile {
	return d.FileByURI(PathToURI(path))
}

// IsRoot returns true if the URI is a known root document.
func (d *Discovery) IsRoot(uri string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, r := range d.roots {
		if r == uri {
			return true
		}
	}
	return false
}

// UpdateFile re-classifies a file and returns the updated DiscoveredFile.
func (d *Discovery) UpdateFile(path string) *DiscoveredFile {
	d.mu.Lock()
	defer d.mu.Unlock()

	uri := PathToURI(path)
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	role := classifyFromDisk(path)
	df := &DiscoveredFile{
		Path:  path,
		URI:   uri,
		Role:  role,
		MTime: info.ModTime().UnixNano(),
	}
	d.files[uri] = df

	d.roots = nil
	for _, f := range d.files {
		if f.Role == RoleRoot {
			d.roots = append(d.roots, f.URI)
		}
	}
	return df
}

// RemoveFile removes a file from the discovery state.
func (d *Discovery) RemoveFile(path string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	uri := PathToURI(path)
	delete(d.files, uri)

	d.roots = nil
	for _, f := range d.files {
		if f.Role == RoleRoot {
			d.roots = append(d.roots, f.URI)
		}
	}
}

func classifyFromDisk(path string) FileRole {
	data, err := os.ReadFile(path)
	if err != nil {
		return RoleUnknown
	}
	return classifyByContent(data)
}

func classifyByContent(data []byte) FileRole {
	const maxScan = 2048
	if len(data) > maxScan {
		data = data[:maxScan]
	}

	if bytes.Contains(data, []byte("openapi:")) || bytes.Contains(data, []byte("\"openapi\"")) ||
		bytes.Contains(data, []byte("swagger:")) || bytes.Contains(data, []byte("\"swagger\"")) {
		return RoleRoot
	}

	return RoleFragment
}
