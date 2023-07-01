package usecase

import (
	"auth-svc/config"
	"auth-svc/internal/auth"
	models "auth-svc/internal/models/auth"
	"auth-svc/pkg/logger"
	"context"
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
