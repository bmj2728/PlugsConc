package registry

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidProtocolVersion  = errors.New("invalid protocol version")
	ErrInvalidMagicCookieKey   = errors.New("invalid magic cookie key")
	ErrInvalidMagicCookieValue = errors.New("invalid magic cookie value")
)

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
	ProtocolVersion  uint   `json:"protocol_version" yaml:"protocol_version"`
	MagicCookieKey   string `json:"magic_cookie_key" yaml:"magic_cookie_key"`
	MagicCookieValue string `json:"magic_cookie_value" yaml:"magic_cookie_value"`
}

// LoadManifest reads and parses a manifest file at the specified path, returning the parsed Manifest,
// its hash, and any error.
func LoadManifest(root, path string) (m *Manifest, entrypoint string, hash string, err error) {
	r, err := os.OpenRoot(root)
	if err != nil {
		err := errors.Join(ErrLoadingFS, err)
		slog.Error("Failed to load plugin root", slog.Any(logger.KeyError, err))
		return nil, "", "", err
	}
	defer func(r *os.Root) {
		err := r.Close()
		if err != nil {
			err := errors.Join(ErrClosingFS, err)
			slog.Error("Failed to close root", slog.Any(logger.KeyError, err))
		}
	}(r)

	rootFS := r.FS()

	f, err := fs.ReadFile(rootFS, path)
	if err != nil {
		err := errors.Join(ErrReadingFile, err)
		slog.Error("Failed to load manifest", slog.Any(logger.KeyError, err))
		return nil, "", "", err
	}

	hash = getMD5Hash(f)

	if err := yaml.Unmarshal(f, &m); err != nil {
		err := errors.Join(ErrYAMLUnmarshaling, err)
		slog.Error("Failed to unmarshall manifest", slog.Any(logger.KeyError, err))
		return nil, "", "", err
	}

	// todo - check if entrypoint is valid executable
	entrypoint = filepath.Join(root, m.PluginEntrypoint)

	return m, entrypoint, hash, nil
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
		slog.Group("handshake_config", slog.Int("protocol_version", int(m.Handshake.ProtocolVersion)),
			slog.String("magic_cookie_key", m.Handshake.MagicCookieKey),
			slog.String("magic_cookie_value", m.Handshake.MagicCookieValue)),
	)
}

// getMD5Hash computes the MD5 hash of the given byte slice and returns it as a hexadecimal string.
func getMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func (m *Manifest) ToLaunchDetails() *PluginLaunchDetails {
	var ld PluginLaunchDetails
	ld.name = m.PluginName
	hc, err := m.Handshake.ToConfig()
	if err != nil {
		slog.Error("Failed to load plugin launch details", slog.Any(logger.KeyError, err))
		return nil
	}
	ld.handshakeConfig = hc
	ld.cmd = exec.Command(m.PluginEntrypoint)
	validFormat := AvailablePluginFormatLookup.IsValidFormat(m.PluginFormat)
	if validFormat {
		pf := AvailablePluginFormats.GetByString(m.PluginFormat)
		ld.allowedProtocols = pf
	}
	return &ld
}

// ToConfig converts a Handshake instance into a HandshakeConfig, validating required fields for correctness.
func (h Handshake) ToConfig() (*plugin.HandshakeConfig, error) {
	if h.ProtocolVersion == 0 {
		return nil, ErrInvalidProtocolVersion
	}
	if h.MagicCookieKey == "" {
		return nil, ErrInvalidMagicCookieKey
	}
	if h.MagicCookieValue == "" {
		return nil, ErrInvalidMagicCookieValue
	}
	return &plugin.HandshakeConfig{
		ProtocolVersion:  h.ProtocolVersion,
		MagicCookieKey:   h.MagicCookieKey,
		MagicCookieValue: h.MagicCookieValue,
	}, nil
}
