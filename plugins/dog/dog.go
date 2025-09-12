package main

import (
	"github.com/bmj2728/PlugsConc/pkg/shared/animal"

	"github.com/hashicorp/go-plugin"
)

type Dog struct {
}

func (d Dog) Speak(isLoud bool) string {
	//return "Woof"
	if isLoud {
		return "WOOF!"
	} else {
		return "Woof"
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ANIMAL_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	dog := Dog{}

	pluginMap := map[string]plugin.Plugin{
		"dog": &animal.AnimalPlugin{Impl: dog},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		//GRPCServer: plugin.DefaultGRPCServer, // add this line to enable grpc
	})
}
