// Package ngfs provides wrappers for various file system functions used by the host file system service
package ngfs

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bmj2728/PlugsConc/internal/logger"
	filesystemv1 "github.com/bmj2728/PlugsConc/shared/protogen/filesystem/v1"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// The host service will receive a request over gRPC.
// We then need to validate that the requestor has permission to perform the action.
// We also need to validate authority to operator on the file/directory.
// This should be achieved elsewhere.
// This package will be used if we need to specifically override functions or provide helpers for more complex
// operations.
// for instance, consider a ReadDir example:
//// BasicReadDir reads the contents of the directory specified by the given path and
//func BasicReadDir(path string) ([]os.DirEntry, error) {
//	return os.ReadDir(path)
//}
//
//// BetterReadDir opens a new root to help protect against path traversal attacks, then read the directory.
//func BetterReadDir(path string) ([]os.DirEntry, error) {
//	r, err := os.OpenRoot(path)
//	if err != nil {
//		hclog.Default().Error("Failed to open root", logger.KeyError, err)
//		return nil, err
//	}
//	defer func(r *os.Root) {
//		err := r.Close()
//		if err != nil {
//			hclog.Default().Error("Failed to close root", logger.KeyError, err)
//		}
//	}(r)
//	// Read the directory, returning the slice of DirEntry and an error, close the root
//	return fs.ReadDir(r.FS(), ".")
//}

type NGFS struct {
	filesystemv1.UnimplementedFileSystemServer
	fsLogger hclog.Logger
}

func NewNGFS() *NGFS {
	return &NGFS{
		fsLogger: logger.DefaultLogger().Named("ngfs"),
	}
}

func (N *NGFS) ReadDir(ctx context.Context, request *filesystemv1.ReadDirRequest) (*filesystemv1.ReadDirResponse, error) {
	r, err := os.OpenRoot(request.Path)
	if err != nil {
		N.fsLogger.Error("Failed to open root", logger.KeyError, err)
		return nil, err
	}
	defer func(r *os.Root) {
		err := r.Close()
		if err != nil {
			N.fsLogger.Error("Failed to close root", logger.KeyError, err)
		}
	}(r)
	entries, err := fs.ReadDir(r.FS(), ".")
	if err != nil {
		N.fsLogger.Error("Failed to read directory", logger.KeyError, err)
		return nil, err
	}
	processedEntries := make([]*filesystemv1.DirEntry, len(entries))
	for i, entry := range entries {
		processedEntries[i] = &filesystemv1.DirEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		}
	}
	return &filesystemv1.ReadDirResponse{
		Entries: processedEntries,
	}, nil
}

func (N *NGFS) Stat(ctx context.Context, request *filesystemv1.StatRequest) (*filesystemv1.StatResponse, error) {
	base, file := filepath.Split(request.Path)
	r, err := os.OpenRoot(base)
	if err != nil {
		N.fsLogger.Error("Failed to open root", logger.KeyError, err)
		return nil, err
	}
	defer func(r *os.Root) {
		err := r.Close()
		if err != nil {
			N.fsLogger.Error("Failed to close root", logger.KeyError, err)
		}
	}(r)
	info, err := fs.Stat(r.FS(), file)
	if err != nil {
		N.fsLogger.Error("Failed to stat file", logger.KeyError, err)
		return nil, err
	}
	return &filesystemv1.StatResponse{
		Info: &filesystemv1.FileInfo{
			Name:    info.Name(),
			Size:    uint64(info.Size()),
			Mode:    uint32(info.Mode()),
			ModTime: timestamppb.New(info.ModTime()),
			IsDir:   info.IsDir(),
			Sys:     nil, // we don't need this yet
		},
	}, nil
}
