package main

import (
	"github.com/bmj2728/PlugsConc/pkg/shared/animal"

	"github.com/hashicorp/go-plugin"
)

type Pig struct {
}

func (p Pig) Speak(isLoud bool) string {
	if isLoud {
		return "OINK!"
	} else {
		return "Oink"
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	pig := Pig{}

	pluginMap := map[string]plugin.Plugin{
		"pig": &animal.AnimalPlugin{Impl: pig},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
