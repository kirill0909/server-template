package http

import (
	"server-template/config"
	"server-template/internal/auth"
	"server-template/pkg/logger"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	log    logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, log: log}
}
