package middleware

import (
	"auth-svc/config"
	"auth-svc/internal/auth"
	"auth-svc/pkg/logger"
)

type MDWManager struct {
	cfg    *config.Config
	log    logger.Logger
	authUC auth.UseCase
}

func NewMDWManager(
	cfg *config.Config,
	log logger.Logger,
	authUC auth.UseCase,
) *MDWManager {
	return &MDWManager{
		cfg:    cfg,
		log:    log,
		authUC: authUC,
	}
}
