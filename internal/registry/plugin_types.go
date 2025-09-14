package registry

import (
	"sync"

	"github.com/bmj2728/PlugsConc/shared/pkg/animal"
	"github.com/hashicorp/go-plugin"
)

// PluginType represents a custom type used for defining various plugin classifications within the system.
type PluginType int

// AnimalPlugin represents a standard animal-related plugin type.
// AnimalGRPCPlugin represents an animal-related plugin type using gRPC.
const (
	AnimalPlugin PluginType = iota
	AnimalGRPCPlugin
)

// AvailablePluginTypes is a global instance of PluginTypes containing mappings of PluginType to their respective
// implementations.
var AvailablePluginTypes = PluginTypes{
	types: map[PluginType]plugin.Plugin{
		AnimalPlugin:     &animal.AnimalPlugin{},
		AnimalGRPCPlugin: &animal.AnimalGRPCPlugin{},
	},
	mu: sync.RWMutex{},
}

// AvailablePluginTypesLookup is a mapping of plugin type names to their corresponding PluginType values.
var AvailablePluginTypesLookup = PluginTypesLookup{
	types: map[string]PluginType{
		"animal":      AnimalPlugin,
		"animal-grpc": AnimalGRPCPlugin,
	},
	mu: sync.RWMutex{},
}

// PluginTypes provides thread-safe storage and retrieval of plugin types, mapped from PluginType to their
// implementations.
type PluginTypes struct {
	types map[PluginType]plugin.Plugin
	mu    sync.RWMutex
}

// Get retrieves the value associated with the given PluginType from the types map in a thread-safe manner.
func (pt *PluginTypes) Get(pluginType PluginType) plugin.Plugin {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	return pt.types[pluginType]
}

// GetByString retrieves the value associated with a plugin type string from the PluginTypes map if it is valid.
func (pt *PluginTypes) GetByString(pluginType string) plugin.Plugin {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	if AvailablePluginTypesLookup.IsValidPluginType(pluginType) {
		return pt.types[AvailablePluginTypesLookup.GetPluginType(pluginType)]
	} else {
		return nil
	}
}

// PluginTypesLookup is a thread-safe structure that maps string keys to PluginType objects for plugin type management.
type PluginTypesLookup struct {
	types map[string]PluginType
	mu    sync.RWMutex
}

// GetPluginType retrieves the PluginType associated with the provided pluginType key from the lookup map.
func (ptl *PluginTypesLookup) GetPluginType(pluginType string) PluginType {
	ptl.mu.RLock()
	defer ptl.mu.RUnlock()
	return ptl.types[pluginType]
}

// IsValidPluginType checks if the given plugin type string exists in the PluginTypesLookup's types map.
func (ptl *PluginTypesLookup) IsValidPluginType(pluginType string) bool {
	ptl.mu.RLock()
	defer ptl.mu.RUnlock()
	_, ok := ptl.types[pluginType]
	return ok
}
