package exten

import (
	"log/slog"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Animal represents a exten interface for animals that extends the base PluginMeta interface.
// It requires implementing a Speak method, defining how the animal communicates.
type Animal interface {
	Speak() string
}

type AnimalRPC struct {
	client *rpc.Client
}

func (a *AnimalRPC) Speak() string {
	var reply string
	err := a.client.Call("Plugin.Speak", new(interface{}), &reply)
	if err != nil {
		slog.With(slog.String("error", err.Error())).Error("error calling Animal.Speak")
	}
	return reply
}

type AnimalRPCServer struct {
	Impl Animal
}

func (arp *AnimalRPCServer) Speak(_ interface{}, resp *string) error {
	*resp = arp.Impl.Speak()
	return nil
}

type AnimalPlugin struct {
	Impl Animal
}

func (ap *AnimalPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &AnimalRPCServer{Impl: ap.Impl}, nil
}

func (ap *AnimalPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &AnimalRPC{client: c}, nil
}
