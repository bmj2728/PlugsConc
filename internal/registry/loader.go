package registry

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
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

func (l LoaderErrors) add(dir string, err error) LoaderErrors {
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

// PluginLoader is responsible for discovering, loading, and managing plugin manifests from a specified directory.
type PluginLoader struct {
	path      string // path to the plugins directory
	manifests *Manifests
}

// NewPluginLoader creates a new instance of PluginLoader for managing plugins in the specified directory path.
func NewPluginLoader(path string) (*PluginLoader, error) {
	loader := &PluginLoader{
		path:      path,
		manifests: NewManifests(),
	}
	return loader, nil
}

// Load initializes and loads plugin manifests from the configured root directory, handling errors during the process.
// Returns maps of plugin manifests and errors encountered during loading, keyed by the directory path.
func (pl *PluginLoader) Load() (*Manifests, LoaderErrors) {
	// Initialize a LoaderErrors map to store errors that occurred during plugin loading
	lErrs := make(LoaderErrors)

	// Initialize the manifests map if it is nil
	if pl.manifests == nil {
		pl.manifests = NewManifests()
	}

	root, err := os.OpenRoot(pl.path)
	if err != nil {
		err = errors.Join(ErrInvalidPluginPath, err)
		slog.Error("Failed to open root", slog.Any("err", err))
		lErrs.add(pl.path, err)
		return pl.manifests, lErrs
	}
	defer func(root *os.Root) {
		err := root.Close()
		if err != nil {
			err = errors.Join(ErrClosingFS, err)
			slog.Error("Failed to close root", slog.Any("err", err))
			lErrs.add(pl.path, err)
		}
	}(root)

	pluginsFS := root.FS()

	err = fs.WalkDir(pluginsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}
		if err != nil && d.IsDir() {
			err = errors.Join(ErrInvalidPluginPath, err)
			slog.Error("Failed to walk directory", slog.Any("err", err))
			absPath, pathErr := filepath.Abs(filepath.Join(pl.path, path))
			if pathErr != nil {
				slog.Error("Failed to get absolute path", slog.Any("err", err))
			}
			if absPath != "" {
				lErrs.add(absPath, err)
			} else {
				lErrs.add(path, err)
			}
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.IsDir() {
			absPluginRoot, absErr := filepath.Abs(filepath.Join(pl.path, path))
			if absErr != nil {
				slog.Error("Failed to get absolute path", slog.Any("err", err))
				// if there is an error getting the absolute path, try to use the relative path instead
				absPluginRoot = filepath.Join(pl.path, path)
			}
			manifest, hash, err := LoadManifest(absPluginRoot, "manifest.yaml")
			if err != nil {
				slog.Error("Failed to load manifest", slog.Any("err", err))
				// if there is an error loading the manifest, add it to the LoaderErrors map
				lErrs.add(absPluginRoot, err)
				// add the manifest to the manifests map (nil/"") to indicate that the manifest is invalid/missing
				// this allows observability for improperly "installed" plugins
				pl.manifests.add(absPluginRoot, NewManifestEntry(manifest, hash))
			}
			// add the manifest to the manifests map
			// TODO add md5 hashing second parameter should be the manifest contents hashed
			pl.manifests.add(absPluginRoot, NewManifestEntry(manifest, hash))
		}
		return nil
	})
	if err != nil {
		err = errors.Join(ErrLoadingFS, err)
		slog.Error("Failed to load plugins", slog.Any("err", err))
		lErrs.add(pl.path, err)
		return pl.manifests, lErrs
	}

	return pl.manifests, lErrs
}

// GetManifests retrieves the collection of plugin manifests currently loaded by the PluginLoader.
func (pl *PluginLoader) GetManifests() *Manifests {
	return pl.manifests
}
