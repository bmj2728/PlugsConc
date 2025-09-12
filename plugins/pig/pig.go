package main

import (
	"PlugsConc/pkg/exten"

	"github.com/hashicorp/go-plugin"
)

type Pig struct {
}

func (p Pig) Speak() string {
	return "Oink"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	pig := Pig{}

	pluginMap := map[string]plugin.Plugin{
		"pig": &exten.AnimalPlugin{Impl: pig},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
