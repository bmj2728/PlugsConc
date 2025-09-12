package main

import (
	"github.com/bmj2728/PlugsConc/pkg/shared/animal"

	"github.com/hashicorp/go-plugin"
)

type Cat struct {
}

func (c Cat) Speak(isLoud bool) string {
	if isLoud {
		return "MEOW!"
	} else {
		return "Meow"
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	cat := Cat{}

	pluginMap := map[string]plugin.Plugin{
		"cat": &animal.AnimalPlugin{Impl: cat},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
