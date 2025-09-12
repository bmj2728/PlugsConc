package main

import (
	"github.com/bmj2728/PlugsConc/pkg/shared/animal"

	"github.com/hashicorp/go-plugin"
)

type Horse struct {
}

func (h Horse) Speak(isLoud bool) string {
	if isLoud {
		return "NEIGH!"
	} else {
		return "Neigh"
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	horse := Horse{}

	pluginMap := map[string]plugin.Plugin{
		"horse": &animal.AnimalPlugin{Impl: horse},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
