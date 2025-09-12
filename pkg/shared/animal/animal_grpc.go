package animal

import (
	"context"

	"github.com/bmj2728/PlugsConc/pkg/protogen/animalpb"
)

type GRPCClient struct {
	client animalpb.AnimalClient
}

func (c *GRPCClient) Speak(isLoud bool) string {
	s, err := c.client.Speak(context.Background(), &animalpb.SpeakRequest{IsLoud: isLoud})
	if err != nil {
		return ""
	}
	return s.GetResp()
}

type GRPCServer struct {
	Impl Animal
	animalpb.UnimplementedAnimalServer
}

func (s *GRPCServer) Speak(_ context.Context, req *animalpb.SpeakRequest) (*animalpb.SpeakResponse, error) {
	return &animalpb.SpeakResponse{Resp: s.Impl.Speak(req.IsLoud)}, nil
}
