package auth

import (
	models "auth-svc/internal/models/auth"
	"context"

	pb "gitlab.axarea.ru/main/aiforypay/package/auth-svc-proto"
)

type UseCase interface {
	GetSession(ctx context.Context, accessToken string, sessionType uint8) (models.Session, error)
	CheckUUIDValid(ctx context.Context, req *pb.CheckUUIDValidRequest) (res *pb.CheckUUIDValidResponse, err error)
	CreateWSToken(ctx context.Context, userID int, sessionType uint8) (res models.GetUUIDResponse, err error)
}
