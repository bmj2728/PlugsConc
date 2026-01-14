package filelister

import (
	"context"
	"fmt"

	filelisterv1 "github.com/bmj2728/PlugsConc/shared/protogen/filelister/v1"
	"github.com/hashicorp/go-plugin"
)

type GRPCClient struct {
	client filelisterv1.FileListerClient
	broker plugin.GRPCBroker
}

func (c *GRPCClient) List(path string) ([]string, error) {
	l, err := c.client.List(context.Background(), &filelisterv1.FileListRequest{Dir: path, HostFsBroker: c.broker.NextId()})
	if err != nil {
		return nil, err
	}
	return l.GetEntry(), nil
}

type GRPCServer struct {
	Impl   FileLister
	broker *plugin.GRPCBroker
	filelisterv1.UnimplementedFileListerServer
}

func (s *GRPCServer) List(ctx context.Context, req *filelisterv1.FileListRequest) (*filelisterv1.FileListResponse, error) {
	entries, err := s.Impl.List(req.Dir)
	if err != nil {
		eStr := fmt.Sprintf("Error: %s", err)
		return &filelisterv1.FileListResponse{Entry: entries, Error: &eStr}, err
	}
	return &filelisterv1.FileListResponse{Entry: entries}, nil
}
