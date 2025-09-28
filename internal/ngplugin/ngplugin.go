package ngplugin

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bmj2728/PlugsConc/internal/checksum"
	"github.com/bmj2728/PlugsConc/internal/registry"
	"github.com/hashicorp/go-plugin"
)

type NGPlugin struct {
	dir        string
	files      PluginFiles          // plugin's directory
	state      registry.PluginState // plugin's current PluginState
	manifest   *registry.Manifest   // plugin's Manifest
	entrypoint *exec.Cmd            // plugin's launch command
	checksum   *plugin.SecureConfig // import of hash from entrypoint.sha256
}

type PluginFiles struct {
	manifestFile string
	binaryFile   string
	checksumFile string
}

func NewPluginFiles(dir string, bin string) PluginFiles {

	mf := filepath.Join(dir, "manifest.yaml")
	bf := filepath.Join(dir, bin)
	sha256 := strings.Join([]string{bf, checksum.CSFileExt}, ".")
	cf := filepath.Join(dir, sha256)

	return PluginFiles{
		manifestFile: mf,
		binaryFile:   bf,
		checksumFile: cf,
	}
}
