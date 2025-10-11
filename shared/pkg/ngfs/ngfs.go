// Package ngfs provides wrappers for various file system functions used by the host file system service
package ngfs

import "os"

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

type NGFS interface {
	ReadDir(path string) ([]os.DirEntry, error)
	Stat(path string) (os.FileInfo, error)
}
