package animal

import (
	"context"
	"net/rpc"

	"github.com/bmj2728/PlugsConc/shared/protogen/animal/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Animal represents a shared interface for animals that extends the base PluginMeta interface.
// It requires implementing a Speak method, defining how the animal communicates.
type Animal interface {
	Speak(isLoud bool) string
}

/**
------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------gRPC--------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
**/

type AnimalGRPCPlugin struct {
	plugin.Plugin
	Impl Animal
}

func (ag *AnimalGRPCPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	animalv1.RegisterAnimalServer(s, &GRPCServer{Impl: ag.Impl})
	return nil
}

func (ag *AnimalGRPCPlugin) GRPCClient(ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn) (interface{}, error) {

	ac := animalv1.NewAnimalClient(c)
	return &GRPCClient{client: ac}, nil
}

/**
------------------------------------------------------------------------------------------------------------------------
------------------------------------------------------RPC---------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
**/

type AnimalPlugin struct {
	Impl Animal
}

func (ap *AnimalPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: ap.Impl}, nil
}

func (ap *AnimalPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}
