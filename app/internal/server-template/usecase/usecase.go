package usecase

import (
	"context"
	"server-template/config"
	"server-template/internal/models/auth"
	"server-template/internal/models/templates"
	server_template "server-template/internal/server-template"
	"server-template/pkg/logger"

	pb "server-template/pkg/proto"
)

type ServerTemplateUC struct {
	serverTemplateRedisRepo server_template.RedisRepo
	log                     logger.Logger
	cfg                     *config.Config
}

func NewServerTemplateUC(serverTemplateRedisRepo server_template.RedisRepo, log logger.Logger, cfg *config.Config) server_template.UseCase {
	return &ServerTemplateUC{serverTemplateRedisRepo: serverTemplateRedisRepo, log: log, cfg: cfg}
}

func (u *ServerTemplateUC) GetSession(ctx context.Context, accessToken string, sessionType uint8) (auth.Session, error) {
	return u.serverTemplateRedisRepo.GetSession(ctx, accessToken, sessionType)
}

func (u *ServerTemplateUC) HTTPPing(ctx context.Context) (result templates.Templates, err error) {
	return u.serverTemplateRedisRepo.HTTPPing(ctx)
}

func (u *ServerTemplateUC) GRPCPing(ctx context.Context) (res *pb.ServerTemplateResponse, err error) {
	return u.serverTemplateRedisRepo.GRPCPing(ctx)
}
