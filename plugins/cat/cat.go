package main

import (
	"PlugsConc/pkg/exten"

	"github.com/hashicorp/go-plugin"
)

type Cat struct {
}

func (c Cat) Speak() string {
	return "Meow"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	cat := Cat{}

	pluginMap := map[string]plugin.Plugin{
		"cat": &exten.AnimalPlugin{Impl: cat},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
