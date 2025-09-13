package registry

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/bmj2728/PlugsConc/internal/logger"
)

// ManifestEntry represents an entry that wraps a Manifest with an associated hash for identifier purposes.
type ManifestEntry struct {
	entry *Manifest
	hash  string
}

// NewManifestEntry creates a new ManifestEntry with the provided Manifest object and hash string.
func NewManifestEntry(manifest *Manifest, hash string) *ManifestEntry {
	return &ManifestEntry{
		entry: manifest,
		hash:  hash,
	}
}

// Manifest returns the Manifest associated with the ManifestEntry instance.
func (m *ManifestEntry) Manifest() *Manifest {
	return m.entry
}

// Hash returns the hash value of the ManifestEntry.
func (m *ManifestEntry) Hash() string {
	return m.hash
}

// LogValue returns a slog.Value representation of the manifest entry by delegating to the underlying Manifest instance.
func (m *ManifestEntry) LogValue() slog.Value {
	return m.entry.LogValue()
}

// Manifests is a map that associates directory paths with their corresponding plugin metadata as Manifest instances.
type Manifests struct {
	mu      sync.RWMutex
	entries map[string]*ManifestEntry
}

func NewManifests() *Manifests {
	return &Manifests{
		mu:      sync.RWMutex{},
		entries: make(map[string]*ManifestEntry),
	}
}

// add inserts a Manifest into the Manifests map, keyed by the specified directory path.
func (m *Manifests) add(dir string, manifest *ManifestEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[dir] = manifest
}

// GetManifests returns a thread-safe copy of all manifest entries in the Manifests map.
func (m *Manifests) GetManifests() map[string]*ManifestEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	clone := make(map[string]*ManifestEntry)
	for k, v := range m.entries {
		entry := *v
		clone[k] = &entry
	}
	return clone
}

// GetEntry retrieves the ManifestEntry associated with the specified directory from the Manifests collection.
// Returns nil if not found.
func (m *Manifests) GetEntry(dir string) *ManifestEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entries[dir]
}

func (m *Manifests) GetManifest(dir string) *Manifest {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entries[dir].Manifest()
}

func (m *Manifests) GetHash(dir string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entries[dir].Hash()
}

func (m *Manifests) LogValue() slog.Value {
	var formatted strings.Builder
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.entries) == 0 {
		return slog.AnyValue("")
	}
	formatted.WriteString("Plugin Manifests:\n")
	for d, e := range m.entries {
		entry := fmt.Sprintf("%s: %s\n", d, e.Manifest().LogValue().String())
		formatted.WriteString(entry)
	}
	return slog.GroupValue(slog.String(logger.KeyPluginMap, formatted.String()))
}
