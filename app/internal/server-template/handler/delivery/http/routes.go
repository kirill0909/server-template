package http

import (
	"server-template/internal/middleware"
	server_template "server-template/internal/server-template"

	"github.com/gofiber/fiber/v2"
)

func MapServerTemplateRoutes(serverTemplateRoutes fiber.Router, h server_template.Handlers, mw *middleware.MDWManager) {
	serverTemplateRoutes.Get("/ping", mw.BrokerAuthedMiddleware(), h.Ping())
}
