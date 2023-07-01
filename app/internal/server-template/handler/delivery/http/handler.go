package http

import (
	"server-template/config"
	server_template "server-template/internal/server-template"
	"server-template/pkg/logger"
)

type serverTeamplateHandlers struct {
	cfg              *config.Config
	serverTemplateUC server_template.UseCase
	log              logger.Logger
}

func NewserverTeamplateHandlers(
	cfg *config.Config, serverTemplateUC server_template.UseCase, log logger.Logger) server_template.Handlers {
	return &serverTeamplateHandlers{cfg: cfg, serverTemplateUC: serverTemplateUC, log: log}
}
