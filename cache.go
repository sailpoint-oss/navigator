package navigator

import (
	"sync"
)

// IndexCache provides a workspace-wide, thread-safe cache of per-document indexes
// with an optional lazy builder for on-demand construction.
type IndexCache struct {
	mu      sync.RWMutex
	indexes map[string]*Index
	builder func(uri string) *Index
}

// NewIndexCache creates a new empty index cache.
func NewIndexCache() *IndexCache {
	return &IndexCache{
		indexes: make(map[string]*Index),
	}
}

// SetBuilder registers a fallback function that builds the index on-demand
// when Get finds no cached entry.
func (c *IndexCache) SetBuilder(fn func(uri string) *Index) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.builder = fn
}

// Get returns the cached index for a URI. If no cached entry exists and a
// builder has been registered, it builds, caches, and returns the index.
func (c *IndexCache) Get(uri string) *Index {
	norm := NormalizeURI(uri)
	c.mu.RLock()
	idx := c.indexes[norm]
	builder := c.builder
	c.mu.RUnlock()

	if idx != nil {
		return idx
	}
	if builder == nil {
		return nil
	}

	idx = builder(uri)
	if idx != nil {
		c.Set(uri, idx)
	}
	return idx
}

// Set stores an index for a URI.
func (c *IndexCache) Set(uri string, idx *Index) {
	norm := NormalizeURI(uri)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.indexes[norm] = idx
}

// Invalidate removes the cached index for a URI, forcing rebuild on next Get.
func (c *IndexCache) Invalidate(uri string) {
	norm := NormalizeURI(uri)
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.indexes, norm)
}

// Delete removes the index for a URI.
func (c *IndexCache) Delete(uri string) {
	c.Invalidate(uri)
}

// All returns a snapshot of all cached indexes.
func (c *IndexCache) All() map[string]*Index {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]*Index, len(c.indexes))
	for k, v := range c.indexes {
		result[k] = v
	}
	return result
}

// Len returns the number of cached indexes.
func (c *IndexCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.indexes)
}

// FindByOperationID searches all cached indexes for an operationId.
func (c *IndexCache) FindByOperationID(opID string) (string, *OperationRef) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for uri, idx := range c.indexes {
		if ref, ok := idx.Operations[opID]; ok {
			return uri, ref
		}
	}
	return "", nil
}

// FindRefTarget searches all cached indexes for a $ref target.
func (c *IndexCache) FindRefTarget(ref string) (string, interface{}) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for uri, idx := range c.indexes {
		if result, err := idx.ResolveRef(ref); err == nil {
			return uri, result
		}
	}
	return "", nil
}
