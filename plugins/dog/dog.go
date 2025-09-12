package main

import (
	"PlugsConc/pkg/exten"

	"github.com/hashicorp/go-plugin"
)

type Dog struct {
}

func (d Dog) Speak() string {
	return "Woof"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	dog := Dog{}

	pluginMap := map[string]plugin.Plugin{
		"animal": &exten.AnimalPlugin{Impl: dog},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
