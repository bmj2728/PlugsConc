package registry

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidProtocolVersion  = errors.New("invalid protocol version")
	ErrInvalidMagicCookieKey   = errors.New("invalid magic cookie key")
	ErrInvalidMagicCookieValue = errors.New("invalid magic cookie value")
)

// Manifest defines the structure for metadata about a plugin,
// including details like PluginName, type, version, and maintainer.
type Manifest struct {
	PluginData PluginData `json:"plugin" yaml:"plugin"`
	About      About      `json:"about" yaml:"about"`
	Handshake  Handshake  `json:"handshake" yaml:"handshake"`
	Security   Security   `json:"security" yaml:"security"`
}

type PluginData struct {
	Name       string `json:"name" yaml:"name"`
	Type       string `json:"type" yaml:"type"`
	Format     string `json:"format" yaml:"format"`
	Entrypoint string `json:"entrypoint" yaml:"entrypoint"`
	Language   string `json:"language" yaml:"language"`
	Version    string `json:"version" yaml:"version"`
}

type About struct {
	Description string `json:"description" yaml:"description"`
	Maintainer  string `json:"maintainer" yaml:"maintainer"`
	URL         string `json:"url" yaml:"url"`
}

// Handshake represents a structure for plugin handshake configuration with protocol version and magic cookie details.
type Handshake struct {
	ProtocolVersion  uint   `json:"protocol_version" yaml:"protocol_version"`
	MagicCookieKey   string `json:"magic_cookie_key" yaml:"magic_cookie_key"`
	MagicCookieValue string `json:"magic_cookie_value" yaml:"magic_cookie_value"`
}

// Security represents configuration related to security features, including automatic mutual TLS (Transport Layer Security).
type Security struct {
	AutoMTLS bool `json:"auto_mtls" yaml:"auto_mtls"`
}

// LoadManifest reads and parses a manifest file at the specified path, returning the parsed Manifest,
// its hash, and any error.
func LoadManifest(root, path string) (m *Manifest, entrypoint string, hash string, err error) {
	r, err := os.OpenRoot(root)
	if err != nil {
		err := errors.Join(ErrLoadingFS, err)
		hclog.Default().Error("Failed to load plugin root", logger.KeyError, err)
		return nil, "", "", err
	}
	defer func(r *os.Root) {
		err := r.Close()
		if err != nil {
			err := errors.Join(ErrClosingFS, err)
			hclog.Default().Error("Failed to close root", logger.KeyError, err)
		}
	}(r)

	rootFS := r.FS()

	f, err := fs.ReadFile(rootFS, path)
	if err != nil {
		err := errors.Join(ErrReadingFile, err)
		hclog.Default().Error("Failed to load manifest", logger.KeyError, err)
		return nil, "", "", err
	}

	hash = getMD5Hash(f)

	if err := yaml.Unmarshal(f, &m); err != nil {
		err := errors.Join(ErrYAMLUnmarshaling, err)
		hclog.Default().Error("Failed to unmarshall manifest", logger.KeyError, err)
		return nil, "", "", err
	}

	entrypoint = filepath.Join(root, m.PluginData.Entrypoint)
	_, err = exec.LookPath(entrypoint)
	if err != nil {
		hclog.Default().Error("Failed to look up entrypoint", logger.KeyError, err)
		return nil, "", "", err
	}

	return m, entrypoint, hash, nil
}

// getMD5Hash computes the MD5 hash of the given byte slice and returns it as a hexadecimal string.
func getMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func (m *Manifest) ToLaunchDetails() *PluginLaunchDetails {
	var ld PluginLaunchDetails
	ld.PluginName = m.PluginData.Name
	hc, err := m.Handshake.ToConfig()
	if err != nil {
		hclog.Default().Error("Failed to load plugin launch details", logger.KeyError, err)
		return nil
	}
	ld.HandshakeConfig = hc
	ld.Cmd = exec.Command(m.PluginData.Entrypoint)
	validFormat := AvailablePluginFormatLookup.IsValidFormat(m.PluginData.Format)
	if validFormat {
		pf := AvailablePluginFormats.GetByString(m.PluginData.Format)
		ld.AllowedProtocols = pf
	}
	ld.AutoMTLS = m.Security.AutoMTLS
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
