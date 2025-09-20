package animal

import (
	"net/rpc"

	"github.com/hashicorp/go-hclog"
)

type RPCClient struct {
	client *rpc.Client
}

func (a *RPCClient) Speak(isLoud bool) string {
	var reply string
	err := a.client.Call("Plugin.Speak", map[string]interface{}{"isLoud": isLoud}, &reply)
	if err != nil {
		hclog.Default().Error("error calling Speak()", "error", err)
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
