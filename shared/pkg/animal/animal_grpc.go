package animal

import (
	"context"

	"github.com/bmj2728/PlugsConc/shared/protogen/animal/v1"
)

type GRPCClient struct {
	client animalv1.AnimalClient
}

func (c *GRPCClient) Speak(isLoud bool) string {
	s, err := c.client.Speak(context.Background(), &animalv1.SpeakRequest{IsLoud: isLoud})
	if err != nil {
		return ""
	}
	return s.GetResp()
}

type GRPCServer struct {
	Impl Animal
	animalv1.UnimplementedAnimalServer
}

func (s *GRPCServer) Speak(_ context.Context, req *animalv1.SpeakRequest) (*animalv1.SpeakResponse, error) {
	return &animalv1.SpeakResponse{Resp: s.Impl.Speak(req.IsLoud)}, nil
}
