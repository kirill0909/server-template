package grpc

import (
	"server-template/internal/auth"

	// pb "gitlab.axarea.ru/main/aiforypay/package/auth-svc-proto"
	pb "server-template/pkg/proto"
)

type AuthHandlers struct {
	pb.UnimplementedAuthServer
	authUC auth.UseCase
}

func NewAuthHandlers(authUC auth.UseCase) pb.AuthServer {
	return &AuthHandlers{authUC: authUC}
}
