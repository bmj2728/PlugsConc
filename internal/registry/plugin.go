package registry

type PluginState int

/**
Plugin directory is scanned
We'll need a plugin struct that's better than the existing approach with multiple object types

**/

const (
	PluginStateUnknown PluginState = iota
	PluginStateDiscovered
	//PluginState
)

type Plugin struct {
}
