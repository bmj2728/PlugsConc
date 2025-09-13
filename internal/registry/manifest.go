package registry

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/fs"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// Provides a way to load a manifest file from a given path. Manifests are used to describe a plugin.
// Manifests can be written in YAML format or JSON format.

// Manifest defines the structure for metadata about a plugin,
// including details like name, type, version, and maintainer.
type Manifest struct {
	PluginName        string    `json:"plugin_name" yaml:"plugin_name"`
	PluginType        string    `json:"plugin_type" yaml:"plugin_type"`
	PluginFormat      string    `json:"plugin_format" yaml:"plugin_format"`
	PluginLanguage    string    `json:"plugin_language" yaml:"plugin_language"`
	PluginEntrypoint  string    `json:"plugin_entrypoint" yaml:"plugin_entrypoint"`
	PluginVersion     string    `json:"plugin_version" yaml:"plugin_version"`
	PluginDescription string    `json:"plugin_description" yaml:"plugin_description"`
	PluginMaintainer  string    `json:"plugin_maintainer" yaml:"plugin_maintainer"`
	PluginURL         string    `json:"plugin_url" yaml:"plugin_url"`
	Handshake         Handshake `json:"handshake" yaml:"handshake"`
}

// Handshake represents a structure for plugin handshake configuration with protocol version and magic cookie details.
type Handshake struct {
	ProtocolVersion  int    `json:"protocol_version" yaml:"protocol_version"`
	MagicCookieKey   string `json:"magic_cookie_key" yaml:"magic_cookie_key"`
	MagicCookieValue string `json:"magic_cookie_value" yaml:"magic_cookie_value"`
}

// LoadManifest reads a manifest file from a specified root and path, parses its YAML content, and returns the Manifest.
// Returns an error if the root cannot be opened, the file cannot be read, or the YAML is invalid.
func LoadManifest(root, path string) (*Manifest, string, error) {
	r, err := os.OpenRoot(root)
	if err != nil {
		err := errors.Join(ErrLoadingFS, err)
		slog.Error("Failed to load plugin root", slog.Any("err", err))
		return nil, "", err
	}
	defer func(r *os.Root) {
		err := r.Close()
		if err != nil {
			err := errors.Join(ErrClosingFS, err)
			slog.Error("Failed to close root", slog.Any("err", err))
		}
	}(r)

	rootFS := r.FS()

	f, err := fs.ReadFile(rootFS, path)
	if err != nil {
		err := errors.Join(ErrReadingFile, err)
		slog.Error("Failed to load manifest", slog.Any("err", err))
		return nil, "", err
	}

	hash := getMD5Hash(f)

	var manifest Manifest
	if err := yaml.Unmarshal(f, &manifest); err != nil {
		err := errors.Join(ErrYAMLUnmarshaling, err)
		slog.Error("Failed to unmarshall manifest", slog.Any("err", err))
		return nil, hash, err
	}

	return &manifest, hash, nil
}

// LogValue converts the Manifest's metadata into a structured slog.Value for logging purposes.
func (m *Manifest) LogValue() slog.Value {
	return slog.GroupValue(slog.String("name", m.PluginName),
		slog.String("version", m.PluginVersion),
		slog.String("type", m.PluginType),
		slog.String("format", m.PluginFormat),
		slog.String("language", m.PluginLanguage),
		slog.String("entrypoint", m.PluginEntrypoint),
		slog.String("description", m.PluginDescription),
		slog.String("maintainer", m.PluginMaintainer),
		slog.String("url", m.PluginURL),
		slog.Group("handshake_config", slog.Int("protocol_version", m.Handshake.ProtocolVersion),
			slog.String("magic_cookie_key", m.Handshake.MagicCookieKey),
			slog.String("magic_cookie_value", m.Handshake.MagicCookieValue)),
	)
}

func getMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
