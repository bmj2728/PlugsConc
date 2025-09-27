package main

import (
	"github.com/bmj2728/PlugsConc/shared/pkg/animal"

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
	MagicCookieKey:   "DOG_PLUGIN",
	MagicCookieValue: "2ggRd5S9bhHottawB6eXwghiOAhekGORmOfIczh5b1D3AYlmrRWIXdbqwDHDJmjq",
}

func main() {
	dog := Dog{}

	pluginMap := map[string]plugin.Plugin{
		"dog-grpc": &animal.AnimalGRPCPlugin{Impl: dog},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
