package checksum

import (
	"crypto"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

const (
	CSFileExt = "sha256"
)

var (
	ErrInvalidChecksum = errors.New("invalid checksum file")
)

func LoadSHA256(root *os.Root, path string) (*plugin.SecureConfig, error) {

	fileBytes, err := fs.ReadFile(root.FS(), path)
	if err != nil {
		err := errors.Join(ErrInvalidChecksum, err)
		if err == nil {
			err = ErrInvalidChecksum
		}
		hclog.Default().Error("Failed to read file", logger.KeyError, err.Error())
		return nil, errors.Join(ErrInvalidChecksum, err)
	}

	rawFields := strings.Fields(string(fileBytes))
	if len(rawFields) == 0 {
		err := ErrInvalidChecksum
		hclog.Default().Error("Failed to parse checksum file", logger.KeyError, err.Error())
		return nil, ErrInvalidChecksum
	}

	hexHash := rawFields[0]
	if hexHash == "" {
		err := ErrInvalidChecksum
		hclog.Default().Error("Failed to parse checksum file", logger.KeyError, err.Error())
		return nil, ErrInvalidChecksum
	}

	checksumBytes, err := hex.DecodeString(hexHash)
	if err != nil {
		err := errors.Join(ErrInvalidChecksum, err)
		if err == nil {
			err = ErrInvalidChecksum
		}
		hclog.Default().Error("Failed to parse checksum file", logger.KeyError, err.Error())
		return nil, errors.Join(ErrInvalidChecksum, err)
	}

	return &plugin.SecureConfig{
		Checksum: checksumBytes,
		Hash:     crypto.SHA256.New(),
	}, nil
}
