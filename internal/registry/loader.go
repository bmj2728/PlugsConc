package registry

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// ErrInvalidPluginPath is returned when the specified plugins directory path is invalid or cannot be accessed.
var (
	ErrInvalidPluginPath = errors.New("invalid plugins directory path")
)

// LoaderErrors is a map that associates a directory with the load error that occurred during its loading process.
type LoaderErrors map[string]error

func (l LoaderErrors) Add(dir string, err error) LoaderErrors {
	l[dir] = err
	return l
}

// Manifests is a map that associates directory paths with their corresponding plugin metadata as Manifest instances.
type Manifests map[string]*Manifest

// Add inserts a Manifest into the Manifests map, keyed by the specified directory path.
func (m Manifests) Add(dir string, manifest *Manifest) {
	m[dir] = manifest
}

// Get retrieves the Manifest associated with the specified directory from the Manifests collection.
// Returns nil if not found.
func (m Manifests) Get(dir string) *Manifest {
	return m[dir]
}

// PluginLoader represents a structure for loading and managing plugins from a specified root directory.
// It holds metadata for each plugin in the form of manifests.
type PluginLoader struct {
	path      string
	root      *os.Root // this should be the plugins directory
	manifests Manifests
}

// NewPluginLoader initializes a PluginLoader for a given plugins directory path and validates its existence.
// Returns a new PluginLoader instance or an error if the path is invalid.
func NewPluginLoader(path string) (*PluginLoader, error) {
	root, err := os.OpenRoot(path)
	if err != nil {
		slog.Error("Failed to open plugins directory", slog.Any("err", errors.Join(ErrInvalidPluginPath, err)))
		return nil, errors.Join(ErrInvalidPluginPath, err)
	}
	loader := &PluginLoader{
		path: path,
		root: root,
	}
	return loader, nil
}

func (pl *PluginLoader) Load() (Manifests, LoaderErrors) {
	// Initialize a LoaderErrors map to store errors that occurred during plugin loading
	lErrs := make(LoaderErrors)

	// Initialize the manifests map or replace it if it already exists
	pl.manifests = make(map[string]*Manifest)

	// Walk the plugins directory pl.path - errors are loaded into lErrs
	_ = filepath.Walk(pl.path, func(path string, info fs.FileInfo, err error) error {
		// if the path is the same as the plugins directory, return nil to skip it
		if path == pl.path {
			return nil
		}
		// if there is an error and the path is a directory, add it to the LoaderErrors map
		if err != nil && info.IsDir() {
			lErrs.Add(path, err)
			return nil
		}
		// if the path is not a directory, return nil to skip it
		if !info.IsDir() {
			return nil
		}
		// if the path is a directory, attempt to load the manifest and add it to the manifests map
		manifestPath := filepath.Join(path, "manifest.yaml")
		manifest, err := LoadManifest(manifestPath)
		if err != nil {
			// if there is an error loading the manifest, add it to the LoaderErrors map
			lErrs.Add(path, err)
		}
		// add the manifest to the manifests map
		pl.manifests.Add(path, manifest)
		return nil
	})
	// return the manifests map and the LoaderErrors map
	return pl.manifests, lErrs
}

func (pl *PluginLoader) GetManifests() Manifests {
	return pl.manifests
}
