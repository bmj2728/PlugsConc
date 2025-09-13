package animal

import (
	"log/slog"
	"net/rpc"
)

type RPCClient struct {
	client *rpc.Client
}

func (a *RPCClient) Speak(isLoud bool) string {
	var reply string
	err := a.client.Call("Plugin.Speak", map[string]interface{}{"isLoud": isLoud}, &reply)
	if err != nil {
		slog.With(slog.String("error", err.Error())).Error("error calling Animal.Speak")
	}
	return reply
}

type RPCServer struct {
	Impl Animal
}

func (arp *RPCServer) Speak(args map[string]interface{}, resp *string) error {
	*resp = arp.Impl.Speak(args["isLoud"].(bool))
	return nil
}
