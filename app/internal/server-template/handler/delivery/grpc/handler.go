package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	server_template "server-template/internal/server-template"

	pb "server-template/pkg/proto"
)

type ServerTemplateHandlers struct {
	pb.UnimplementedServerTemplateServer
	serverTemplateUC server_template.UseCase
}

func NewServerTemplateHandlers(serverTemplateUC server_template.UseCase) pb.ServerTemplateServer {
	return &ServerTemplateHandlers{serverTemplateUC: serverTemplateUC}
}

func (s *ServerTemplateHandlers) Ping(ctx context.Context, empty *emptypb.Empty) (res *pb.ServerTemplateResponse, err error) {
	return
}
