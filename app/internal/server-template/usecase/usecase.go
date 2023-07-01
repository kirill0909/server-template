package usecase

import (
	"context"
	"server-template/config"
	models "server-template/internal/models/auth"
	server_template "server-template/internal/server-template"
	"server-template/pkg/logger"
)

type ServerTemplateUC struct {
	serverTemplateRedisRepo server_template.RedisRepo
	log                     logger.Logger
	cfg                     *config.Config
}

func NewServerTemplateUC(serverTemplateRedisRepo server_template.RedisRepo, log logger.Logger, cfg *config.Config) server_template.UseCase {
	return &ServerTemplateUC{serverTemplateRedisRepo: serverTemplateRedisRepo, log: log, cfg: cfg}
}

func (u *ServerTemplateUC) GetSession(ctx context.Context, accessToken string, sessionType uint8) (models.Session, error) {
	return u.serverTemplateRedisRepo.GetSession(ctx, accessToken, sessionType)
}
