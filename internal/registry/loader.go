package registry

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/bmj2728/PlugsConc/internal/logger"
)

// ErrInvalidPluginPath is returned when the specified plugins directory path is invalid or cannot be accessed.
var (
	ErrInvalidPluginPath = errors.New("invalid plugins directory path")
	ErrLoadingFS         = errors.New("failed to load plugin files")
	ErrClosingFS         = errors.New("failed to close plugin files")
	ErrReadingFile       = errors.New("failed to read file")
	ErrYAMLUnmarshaling  = errors.New("failed to unmarshal YAML")
)

// LoaderErrors is a map that associates a directory with the load error that occurred during its loading process.
type LoaderErrors map[string]error

func (l LoaderErrors) Add(dir string, err error) LoaderErrors {
	l[dir] = err
	return l
}

func (l LoaderErrors) LogValue() slog.Value {
	var formatted strings.Builder
	if len(l) == 0 {
		return slog.AnyValue("")
	}
	formatted.WriteString("Plugin Loading Errors:\n")
	for d, e := range l {
		entry := fmt.Sprintf("%s: %s\n", d, e.Error())
		formatted.WriteString(entry)
	}
	return slog.GroupValue(slog.String(logger.KeyPluginLoadErrors, formatted.String()))
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

func (m Manifests) LogValue() slog.Value {
	var formatted strings.Builder
	if len(m) == 0 {
		return slog.AnyValue("")
	}
	formatted.WriteString("Plugin Manifests:\n")
	for d, m := range m {
		entry := fmt.Sprintf("%s: %s\n", d, m.LogValue().String())
		formatted.WriteString(entry)
	}
	return slog.GroupValue(slog.String(logger.KeyPluginMap, formatted.String()))
}

// PluginLoader represents a structure for loading and managing plugins from a specified root directory.
// It holds metadata for each plugin in the form of manifests.
type PluginLoader struct {
	path      string // path to the plugins directory
	manifests Manifests
}

// NewPluginLoader initializes a PluginLoader for a given plugins directory path and validates its existence.
// Returns a new PluginLoader instance or an error if the path is invalid.
func NewPluginLoader(path string) (*PluginLoader, error) {
	loader := &PluginLoader{
		path: path,
	}
	return loader, nil
}

func (pl *PluginLoader) Load() (Manifests, LoaderErrors) {
	// Initialize a LoaderErrors map to store errors that occurred during plugin loading
	lErrs := make(LoaderErrors)

	// Initialize the manifests map or replace it if it already exists
	pl.manifests = make(map[string]*Manifest)

	//root, err := os.OpenRoot(pl.path)
	//if err != nil {
	//	lErrs.Add(pl.path, errors.Join(ErrInvalidPluginPath, err))
	//	return pl.manifests, lErrs
	//}
	//defer func(root *os.Root) {
	//	err := root.Close()
	//	if err != nil {
	//		slog.Error("Failed to close root", slog.Any("err", err))
	//	}
	//}(root)
	//
	//pluginFS := root.FS()
	//
	//err = fs.WalkDir(pluginFS, ".", func(path string, d fs.DirEntry, err error) error {
	//	if path == pl.path {
	//		return nil
	//	}
	//	if err != nil && d.IsDir() {
	//		lErrs.Add(path, errors.Join(ErrInvalidPluginPath, err))
	//		return err
	//	}
	//	if !d.IsDir() {
	//		return nil
	//	}
	//	if d.IsDir() {
	//		manifestPath := filepath.Join(path, "manifest.yaml")
	//		manifest, err := LoadManifest(path, manifestPath)
	//		if err != nil {
	//			// if there is an error loading the manifest, add it to the LoaderErrors map
	//			lErrs.Add(path, err)
	//		}
	//		// add the manifest to the manifests map
	//		pl.manifests.Add(path, manifest)
	//	}
	//	return nil
	//})
	//if err != nil {
	//	return pl.manifests, lErrs
	//}

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
