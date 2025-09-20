package registry

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/hashicorp/go-hclog"
)

// ErrInvalidPluginPath is returned when the specified plugins directory path is invalid or cannot be accessed.
var (
	ErrInvalidPluginPath = errors.New("invalid plugins directory path")
	ErrLoadingFS         = errors.New("failed to load plugin files")
	ErrClosingFS         = errors.New("failed to close plugin files")
	ErrReadingFile       = errors.New("failed to read file")
	ErrYAMLUnmarshaling  = errors.New("failed to unmarshal YAML")
)

const (
	ManifestFileName = "manifest.yaml"
	ConfigFileSuffix = ".config.yaml"
)

type PluginPaths struct {
	root       string
	entrypoint string
	sha256     string
	manifest   string
	config     string
}

// LoaderErrors is a map that associates a directory with the load error that occurred during its loading process.
type LoaderErrors map[string]error

func (l LoaderErrors) add(dir string, err error) LoaderErrors {
	l[dir] = err
	return l
}

// PluginLoader is responsible for discovering, loading, and managing plugin manifests from a specified directory.
type PluginLoader struct {
	loadLogger hclog.Logger
	path       string // path to the plugins directory
	manifests  *Manifests
}

// NewPluginLoader initializes a new PluginLoader for managing plugins in the specified directory path.
func NewPluginLoader(path string, loadLogger hclog.Logger) (*PluginLoader, error) {
	if loadLogger == nil {
		loadLogger = hclog.Default()
	}
	loader := &PluginLoader{
		loadLogger: loadLogger,
		path:       path,
		manifests:  NewManifests(),
	}
	return loader, nil
}

// Load discovers, parses, and loads plugin manifests from the specified directory, returning manifests and load errors.
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
		pl.loadLogger.Error("Failed to open root", logger.KeyError, err)
		lErrs.add(pl.path, err)
		return pl.manifests, lErrs
	}
	defer func(root *os.Root) {
		err := root.Close()
		if err != nil {
			err = errors.Join(ErrClosingFS, err)
			pl.loadLogger.Error("Failed to close root", logger.KeyError, err)
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
			pl.loadLogger.Error("Failed to walk directory", logger.KeyError, err)
			absPath, pathErr := filepath.Abs(filepath.Join(pl.path, path))
			if pathErr != nil {
				pl.loadLogger.Error("Failed to get absolute path", logger.KeyError, err)
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
				pl.loadLogger.Error("Failed to get absolute path", logger.KeyError, err)
				// if there is an error getting the absolute path, try to use the relative path instead
				absPluginRoot = filepath.Join(pl.path, path)
			}
			manifest, entrypoint, hash, err := LoadManifest(absPluginRoot, ManifestFileName)
			if err != nil {
				pl.loadLogger.Error("Failed to load manifest", logger.KeyError, err)
				// if there is an error loading the manifest, Add it to the LoaderErrors map
				lErrs.add(absPluginRoot, err)
				// Add the manifest to the manifests map (nil/"") to indicate that the manifest is invalid/missing
				// this allows observability for improperly "installed" plugins
				pl.manifests.Add(absPluginRoot, NewManifestEntry(manifest, entrypoint, hash))
			}
			// Add the manifest to the manifest entry map
			pl.manifests.Add(absPluginRoot, NewManifestEntry(manifest, entrypoint, hash))
		}
		return nil
	})
	if err != nil {
		err = errors.Join(ErrLoadingFS, err)
		pl.loadLogger.Error("Failed to load plugins", logger.KeyError, err)
		lErrs.add(pl.path, err)
		return pl.manifests, lErrs
	}

	return pl.manifests, lErrs
}

// GetManifests returns a reference to the loaded plugin manifests managed by the PluginLoader.
func (pl *PluginLoader) GetManifests() *Manifests {
	return pl.manifests
}
