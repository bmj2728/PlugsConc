package main

import (
	"PlugsConc/pkg/exten"

	"github.com/hashicorp/go-plugin"
)

type Cow struct {
}

func (c Cow) Speak() string {
	return "Moo"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	cow := Cow{}

	pluginMap := map[string]plugin.Plugin{
		"cow": &exten.AnimalPlugin{Impl: cow},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
