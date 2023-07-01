package middleware

import (
	"server-template/config"
	"server-template/internal/auth"
	"server-template/pkg/logger"
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
