package middleware

import (
	"server-template/config"
	server_template "server-template/internal/server-template"
	"server-template/pkg/logger"
)

type MDWManager struct {
	cfg              *config.Config
	log              logger.Logger
	serverTemplateUC server_template.UseCase
}

func NewMDWManager(
	cfg *config.Config,
	log logger.Logger,
	serverTemplateUC server_template.UseCase,
) *MDWManager {
	return &MDWManager{
		cfg:              cfg,
		log:              log,
		serverTemplateUC: serverTemplateUC,
	}
}
