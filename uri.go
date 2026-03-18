package navigator

import (
	"net/url"
	"path/filepath"
	"strings"
)

// NormalizeURI canonicalizes a file:// URI by cleaning the path and
// removing query/fragment components.
func NormalizeURI(uri string) string {
	if uri == "" {
		return ""
	}
	if !strings.HasPrefix(uri, "file://") {
		return uri
	}
	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	u.Host = ""
	u.RawQuery = ""
	u.Fragment = ""
	u.Path = filepath.Clean(u.Path)
	return u.String()
}

// PathToURI converts a filesystem path to a file:// URI.
func PathToURI(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	abs = filepath.ToSlash(abs)
	if !strings.HasPrefix(abs, "/") {
		abs = "/" + abs
	}
	return "file://" + abs
}

// URIToPath converts a file:// URI to a filesystem path.
func URIToPath(uri string) string {
	if !strings.HasPrefix(uri, "file://") {
		return uri
	}
	u, err := url.Parse(uri)
	if err != nil {
		return strings.TrimPrefix(uri, "file://")
	}
	return filepath.FromSlash(u.Path)
}

// ResolveRelativeURI resolves a relative path against a base URI.
func ResolveRelativeURI(baseURI, relPath string) string {
	if strings.Contains(relPath, "://") {
		return relPath
	}
	// Split off fragment from relPath
	fragment := ""
	if i := strings.Index(relPath, "#"); i >= 0 {
		fragment = relPath[i:]
		relPath = relPath[:i]
	}
	if relPath == "" {
		return baseURI + fragment
	}

	basePath := URIToPath(baseURI)
	baseDir := filepath.Dir(basePath)
	resolved := filepath.Join(baseDir, relPath)
	resolved = filepath.Clean(resolved)

	return PathToURI(resolved) + fragment
}

// SplitRefURI splits a $ref value into its file URI and JSON pointer fragment.
// For local refs like "#/components/schemas/Pet", fileURI is empty.
// For external refs like "./file.yaml#/components/schemas/Pet", both are returned.
func SplitRefURI(ref string) (fileURI, fragment string) {
	if strings.HasPrefix(ref, "#") {
		return "", ref
	}
	if i := strings.Index(ref, "#"); i >= 0 {
		return ref[:i], ref[i:]
	}
	return ref, ""
}
