package registry

import (
	"sync"

	"github.com/hashicorp/go-plugin"
)

// PluginFormat represents the type for defining various plugin communication formats.
type PluginFormat int

// GRPC represents a plugin format using gRPC.
// RPC represents a plugin format using RPC.
const (
	GRPC PluginFormat = iota
	RPC
)

// PluginFormats is a struct that manages a thread-safe map of PluginFormat values to their string representations.
type PluginFormats struct {
	formats map[PluginFormat][]plugin.Protocol
	mu      sync.RWMutex
}

// AvailablePluginFormats defines a mapping between PluginFormat constants and their string representations.
var AvailablePluginFormats = PluginFormats{
	formats: map[PluginFormat][]plugin.Protocol{
		GRPC: {plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		RPC:  {plugin.ProtocolNetRPC},
	},
	mu: sync.RWMutex{},
}

func (pf *PluginFormats) Get(format PluginFormat) []plugin.Protocol {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	return pf.formats[format]
}

func (pf *PluginFormats) GetByString(format string) []plugin.Protocol {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	return pf.formats[AvailablePluginFormatLookup.GetPluginFormat(format)]
}

type PluginFormatLookup struct {
	mu      sync.RWMutex
	formats map[string]PluginFormat
}

// AvailablePluginFormatLookup is a pre-initialized PluginFormatLookup containing supported plugin formats
// with thread safety.
var AvailablePluginFormatLookup = PluginFormatLookup{
	formats: map[string]PluginFormat{
		"grpc": GRPC,
		"rpc":  RPC,
	},
	mu: sync.RWMutex{},
}

// GetPluginFormat retrieves the PluginFormat associated with the given format string from the lookup.
func (pfl *PluginFormatLookup) GetPluginFormat(format string) PluginFormat {
	pfl.mu.RLock()
	defer pfl.mu.RUnlock()
	return pfl.formats[format]
}

// IsValidFormat checks if the provided format string exists as a key in the PluginFormatLookup map.
// Returns true if valid.
func (pfl *PluginFormatLookup) IsValidFormat(format string) bool {
	_, ok := pfl.formats[format]
	return ok
}
