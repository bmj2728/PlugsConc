package plugin

// Plugin defines a contract for plugins with a name, version, and type properties.
type Plugin interface {
	Name() string
	Version() string
	Type() string
}

// AnimalPlugin represents a plugin interface for animals that extends the base Plugin interface.
// It requires implementing a Speak method, defining how the animal communicates.
type AnimalPlugin interface {
	Plugin
	Speak()
}
