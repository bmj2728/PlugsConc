package main

import (
	"github.com/bmj2728/PlugsConc/shared/pkg/animal"

	"github.com/hashicorp/go-plugin"
)

type Cow struct {
}

func (c Cow) Speak(isLoud bool) string {
	if isLoud {
		return "MOO!"
	} else {
		return "Moo"
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	cow := Cow{}

	pluginMap := map[string]plugin.Plugin{
		"cow": &animal.AnimalPlugin{Impl: cow},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
