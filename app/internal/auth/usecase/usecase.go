package usecase

import (
	"auth-svc/config"
	"auth-svc/internal/auth"
	models "auth-svc/internal/models/auth"
	"auth-svc/pkg/logger"
	"context"

	"github.com/google/uuid"
	pb "gitlab.axarea.ru/main/aiforypay/package/auth-svc-proto"
	"go.opentelemetry.io/otel"
)

type AuthUC struct {
	authRedisRepo auth.RedisRepo
	log           logger.Logger
	cfg           *config.Config
}

func NewAuthUC(
	authRedisRepo auth.RedisRepo,
	log logger.Logger,
	cfg *config.Config,
) auth.UseCase {
	return &AuthUC{
		authRedisRepo: authRedisRepo,
		log:           log,
		cfg:           cfg,
	}
}

func (u *AuthUC) GetSession(ctx context.Context, accessToken string, sessionType uint8) (models.Session, error) {
	return u.authRedisRepo.GetSession(ctx, accessToken, sessionType)
}

func (u *AuthUC) CheckUUIDValid(
	ctx context.Context,
	req *pb.CheckUUIDValidRequest) (res *pb.CheckUUIDValidResponse, err error) {
	ctx, span := otel.Tracer("").Start(ctx, "authUC.CheckUUIDValid")
	defer span.End()

	return u.authRedisRepo.CheckUUIDValid(ctx, req)
}

func (u *AuthUC) CreateWSToken(ctx context.Context, userID int, sessionType uint8) (res models.GetUUIDResponse, err error) {
	ctx, span := otel.Tracer("").Start(ctx, "AuthUC.CreateWSToken")
	defer span.End()

	res.UUID = uuid.New().String()
	if err = u.authRedisRepo.SetUUID(ctx, res.UUID, userID, sessionType); err != nil {
		return models.GetUUIDResponse{}, err
	}
	return res, nil
}
