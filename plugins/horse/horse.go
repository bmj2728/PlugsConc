package main

import (
	"PlugsConc/pkg/exten"

	"github.com/hashicorp/go-plugin"
)

type Horse struct {
}

func (h Horse) Speak() string {
	return "Neigh"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	horse := Horse{}

	pluginMap := map[string]plugin.Plugin{
		"horse": &exten.AnimalPlugin{Impl: horse},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
