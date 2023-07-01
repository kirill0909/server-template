package grpc

import (
	"auth-svc/internal/auth"
	"context"

	pb "gitlab.axarea.ru/main/aiforypay/package/auth-svc-proto"

	"go.opentelemetry.io/otel"
)

type AuthHandlers struct {
	pb.UnimplementedAuthServer
	authUC auth.UseCase
}

func NewAuthHandlers(authUC auth.UseCase) pb.AuthServer {
	return &AuthHandlers{authUC: authUC}
}

func (h *AuthHandlers) CheckUUIDValid(
	ctx context.Context,
	req *pb.CheckUUIDValidRequest) (res *pb.CheckUUIDValidResponse, err error) {
	ctx, span := otel.Tracer("").Start(ctx, "AuthHandlers.CheckUUIDValid")
	defer span.End()

	return h.authUC.CheckUUIDValid(ctx, req)
}
