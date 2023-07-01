package http

import (
	"server-template/config"
	server_template "server-template/internal/server-template"
	"server-template/pkg/errors"
	"server-template/pkg/logger"

	"github.com/gofiber/fiber/v2"
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

func (h *serverTeamplateHandlers) Ping() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userID, ok := c.Locals("userID").(int)
		if !ok {
			err := errors.ErrGetUserIDFromCtx
			h.log.Errorf("Error: %v\nauth.delivery.http.handlers.GetUUID() cannot get userID from fiber ctx", err)
			return err
		}

		h.log.Infof("UserID: %d", userID)

		result, err := h.serverTemplateUC.Ping(c.Context())
		if err != nil {
			h.log.Error(err)
			return err
		}

		return c.JSON(result)
	}
}
