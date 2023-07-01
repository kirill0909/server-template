package grpc

import (
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
