package registry

import (
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

// LoadManifest reads a YAML manifest file from the given path and unmarshals its contents into a Manifest struct.
// Returns a pointer to the Manifest and an error if the file cannot be read or unmarshalled.
// If the provided path is empty, it returns nil without error.
func LoadManifest(path string) (*Manifest, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

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
