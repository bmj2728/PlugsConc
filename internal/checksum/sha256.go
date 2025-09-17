package checksum

import (
	"crypto"
	"encoding/hex"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/hashicorp/go-plugin"
)

const (
	ChecksumFileExt = ".sha256"
)

var (
	ErrInvalidChecksum = errors.New("invalid checksum file")
)

func LoadSHA256(root *os.Root, path string) (*plugin.SecureConfig, error) {

	fileBytes, err := fs.ReadFile(root.FS(), path)
	if err != nil {
		err := errors.Join(ErrInvalidChecksum, err)
		slog.Error("Failed to read file", slog.Any(logger.KeyError, err))
		return nil, errors.Join(ErrInvalidChecksum, err)
	}

	rawFields := strings.Fields(string(fileBytes))
	if len(rawFields) == 0 {
		err := ErrInvalidChecksum
		slog.Error("Failed to parse checksum file", slog.Any(logger.KeyError, err))
		return nil, ErrInvalidChecksum
	}

	hexHash := rawFields[0]
	if hexHash == "" {
		err := ErrInvalidChecksum
		slog.Error("Failed to parse checksum file", slog.Any(logger.KeyError, err))
		return nil, ErrInvalidChecksum
	}

	checksumBytes, err := hex.DecodeString(hexHash)
	if err != nil {
		err := errors.Join(ErrInvalidChecksum, err)
		slog.Error("Failed to parse checksum file", slog.Any(logger.KeyError, err))
		return nil, errors.Join(ErrInvalidChecksum, err)
	}

	return &plugin.SecureConfig{
		Checksum: checksumBytes,
		Hash:     crypto.SHA256.New(),
	}, nil
}
