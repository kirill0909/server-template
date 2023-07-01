package http

import (
	"server-template/internal/auth"
	"server-template/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapAuthRoutes(authRoutes fiber.Router, h auth.Handlers, mw *middleware.MDWManager) {
}
