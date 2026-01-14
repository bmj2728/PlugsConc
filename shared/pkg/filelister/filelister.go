package filelister

import (
	"context"

	filelisterv1 "github.com/bmj2728/PlugsConc/shared/protogen/filelister/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type FileLister interface {
	List(path string) ([]string, error)
}

type FileListerGRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl FileLister
}

func (f *FileListerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	filelisterv1.RegisterFileListerServer(s,
		&GRPCServer{
			Impl:   f.Impl,
			broker: broker,
		})
	return nil
}

func (f *FileListerGRPCPlugin) GRPCClient(ctx context.Context,
	broker *plugin.GRPCBroker,
	conn *grpc.ClientConn) (interface{}, error) {
	flc := filelisterv1.NewFileListerClient(conn)
	return &GRPCClient{
		client: flc}, nil
}
