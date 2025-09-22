package registry

import (
	"context"
	"os/exec"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/go-plugin"
)

// PluginCatalog provides a thread-safe structure for managing plugins, their manifests, launch details,
// and file watchers.
type PluginCatalog struct {
	mu            sync.RWMutex
	manifests     *Manifests
	pluginMap     map[string]plugin.Plugin // this is passed to each client config
	launchDetails []*PluginLaunchDetails   // these are passed to the plugin launcher
	fw            *fsnotify.Watcher
	watch         func(ctx context.Context, fw *fsnotify.Watcher)
}

// NewPluginCatalog creates and initializes a new PluginCatalog instance with the given manifests.
func NewPluginCatalog(manifests *Manifests) *PluginCatalog {
	return &PluginCatalog{
		manifests:     manifests,
		mu:            sync.RWMutex{},
		pluginMap:     make(map[string]plugin.Plugin),
		launchDetails: make([]*PluginLaunchDetails, 0),
	}
}

// GetPlugin retrieves a plugin from the catalog by its PluginName in a thread-safe manner. Returns nil if not found.
func (c *PluginCatalog) GetPlugin(name string) plugin.Plugin {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pluginMap[name]
}

// AddPlugin adds a plugin to the catalog, associating it with the specified PluginName in a thread-safe manner.
func (c *PluginCatalog) AddPlugin(name string, plugin plugin.Plugin) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pluginMap[name] = plugin
}

// GetLaunchDetails retrieves the list of PluginLaunchDetails currently stored in the PluginCatalog.
func (c *PluginCatalog) GetLaunchDetails() []*PluginLaunchDetails {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.launchDetails
}

// AddLaunchDetails adds a new PluginLaunchDetails object to the catalog in a thread-safe manner.
func (c *PluginCatalog) AddLaunchDetails(details *PluginLaunchDetails) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.launchDetails = append(c.launchDetails, details)
}

// WithFileWatcher sets the file watcher for the PluginCatalog and returns the updated instance.
func (c *PluginCatalog) WithFileWatcher(fw *fsnotify.Watcher,
	watch func(ctx context.Context, fw *fsnotify.Watcher)) *PluginCatalog {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.fw = fw
	c.watch = watch
	return c
}

// AddWatch adds a new directory path to the file watcher and returns an error if it fails.
// Uses a read lock for thread safety.
func (c *PluginCatalog) AddWatch(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.fw.Add(path)
}

// GetWatchlist retrieves the current list of file paths being watched by the file watcher in the PluginCatalog.
func (c *PluginCatalog) GetWatchlist() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.fw.WatchList()
}

// SetWatchFunc sets a function to handle filesystem watch events, synchronizing access with a mutex.
func (c *PluginCatalog) SetWatchFunc(watch func(ctx context.Context, fw *fsnotify.Watcher)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watch = watch
}

// PluginLaunchDetails represents the details required to launch a plugin including its configuration
// and execution command.
// PluginName is the identifier for the plugin.
// HandshakeConfig specifies the handshake configuration needed for the plugin communication.
// Cmd holds the execution command for running the plugin.
// AllowedProtocols lists the communication protocols supported by the plugin.
type PluginLaunchDetails struct {
	PluginName       string                  `json:"plugin_name" yaml:"plugin_name"`
	HandshakeConfig  *plugin.HandshakeConfig `json:"handshake_config" yaml:"handshake_config"`
	Cmd              *exec.Cmd               `json:"Cmd" yaml:"Cmd"`
	AllowedProtocols []plugin.Protocol       `json:"allowed_protocols" yaml:"allowed_protocols"`
	AutoMTLS         bool                    `json:"auto_mtls" yaml:"auto_mtls"`
}

// NewPluginLaunchDetails initializes a new PluginLaunchDetails instance with the specified parameters.
func NewPluginLaunchDetails(name string,
	handshakeConfig *plugin.HandshakeConfig,
	cmd *exec.Cmd,
	allowedProtocols []plugin.Protocol,
	autoMTLS bool) *PluginLaunchDetails {
	return &PluginLaunchDetails{
		PluginName:       name,
		HandshakeConfig:  handshakeConfig,
		Cmd:              cmd,
		AllowedProtocols: allowedProtocols,
		AutoMTLS:         autoMTLS,
	}
}

// Name returns the PluginName of the plugin instance.
func (p *PluginLaunchDetails) Name() string {
	return p.PluginName
}

// Handshake returns the handshake configuration for the plugin, defining authentication and protocol version details.
func (p *PluginLaunchDetails) Handshake() *plugin.HandshakeConfig {
	return p.HandshakeConfig
}

// Entrypoint returns the command (`*exec.Cmd`) associated with the plugin's launch details.
func (p *PluginLaunchDetails) Entrypoint() *exec.Cmd {
	return p.Cmd
}

// PluginAllowedProtocols returns the list of protocols allowed for the plugin.
func (p *PluginLaunchDetails) PluginAllowedProtocols() []plugin.Protocol {
	return p.AllowedProtocols
}
