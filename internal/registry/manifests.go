package registry

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/bmj2728/PlugsConc/internal/logger"
)

// ManifestEntry represents an entry containing a plugin's manifest and associated hash for identifying integrity.
type ManifestEntry struct {
	entry      *Manifest
	entrypoint string
	hash       string
}

// NewManifestEntry creates a new ManifestEntry instance, associating a manifest with its corresponding hash.
func NewManifestEntry(manifest *Manifest, entrypoint string, hash string) *ManifestEntry {
	return &ManifestEntry{
		entry:      manifest,
		entrypoint: entrypoint,
		hash:       hash,
	}
}

// Manifest retrieves the Manifest structure associated with the current ManifestEntry instance.
func (m *ManifestEntry) Manifest() *Manifest {
	return m.entry
}

// Hash returns the hash value associated with the ManifestEntry.
func (m *ManifestEntry) Hash() string {
	return m.hash
}

func (m *ManifestEntry) Entrypoint() string {
	return m.entrypoint
}

// LogValue returns a slog.Value representing the loggable structure of the associated Manifest
// object within ManifestEntry.
func (m *ManifestEntry) LogValue() slog.Value {
	return m.entry.LogValue()
}

// Manifests is a thread-safe structure for managing a collection of ManifestEntry objects with synchronized access.
type Manifests struct {
	mu      sync.RWMutex
	entries map[string]*ManifestEntry
}

// NewManifests creates and returns a new instance of Manifests with initialized fields.
func NewManifests() *Manifests {
	return &Manifests{
		mu:      sync.RWMutex{},
		entries: make(map[string]*ManifestEntry),
	}
}

// Add inserts a ManifestEntry into the manifests map, associating it with a specified directory path.
func (m *Manifests) Add(dir string, manifest *ManifestEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[dir] = manifest
}

// GetManifests returns a clone of the current map of manifest entries ensuring thread-safe access.
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

// GetEntry retrieves a ManifestEntry for the specified directory from the Manifests collection in a thread-safe manner.
func (m *Manifests) GetEntry(dir string) *ManifestEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entries[dir]
}

// GetManifest retrieves the Manifest object corresponding to the provided directory from the manifests map.
func (m *Manifests) GetManifest(dir string) *Manifest {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entries[dir].Manifest()
}

// GetHash retrieves the hash value of the manifest entry associated with the given directory path.
func (m *Manifests) GetHash(dir string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entries[dir].Hash()
}

// LogValue generates a structured slog.Value representing the current state of all plugin manifests in the collection.
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
