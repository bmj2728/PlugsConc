package checksum

import (
	"crypto"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

const (
	CSFileName = "plugin.sha256"
)

var (
	ErrInvalidChecksum     = errors.New("invalid checksum file")
	ErrInvalidChecksumPath = errors.New("invalid checksum file path")
)

type SHA256File struct {
	path     string
	hexHash  string
	fileName string
}

func NewSHA256File(dir string) (*SHA256File, error) {
	// Check that path is not empty
	if dir == "" {
		return nil, ErrInvalidChecksumPath
	}
	// Get the absolute path
	aPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, ErrInvalidChecksumPath
	}
	// Create the SHA256File
	sf := &SHA256File{
		path: aPath,
	}
	return sf, nil
}

func (sf *SHA256File) Path() string {
	return sf.path
}

func (sf *SHA256File) Parse() error {
	r, err := os.OpenRoot(sf.path)
	if err != nil {
		err = errors.Join(ErrInvalidChecksumPath, err)
		hclog.Default().Error("Failed to open checksum file", logger.KeyError, err)
		return err
	}
	defer func(r *os.Root) {
		err := r.Close()
		if err != nil {
			hclog.Default().Error("Failed to close checksum file", logger.KeyError, err)
		}
	}(r)

	fileBytes, err := fs.ReadFile(r.FS(), CSFileName)
	if err != nil {
		err := errors.Join(ErrInvalidChecksum, err)
		hclog.Default().Error("Failed to read checksum file", logger.KeyError, err)
		return err
	}

	rawFields := strings.Fields(string(fileBytes))
	if len(rawFields) != 2 {
		err := errors.Join(ErrInvalidChecksum, err)
		hclog.Default().Error("Failed to parse checksum file", logger.KeyError, err)
		return err
	}

	sf.hexHash = rawFields[0]
	sf.fileName = rawFields[1]

	return nil
}

func (sf *SHA256File) SecConf() (*plugin.SecureConfig, error) {
	if sf.hexHash == "" {
		return nil, ErrInvalidChecksum
	}
	checksumBytes, err := hex.DecodeString(sf.hexHash)
	if err != nil {
		err := errors.Join(ErrInvalidChecksum, err)
		hclog.Default().Error("Failed to parse checksum file", logger.KeyError, err)
		return nil, err
	}
	return &plugin.SecureConfig{
		Checksum: checksumBytes,
		Hash:     crypto.SHA256.New(),
	}, nil
}
