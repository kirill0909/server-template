package http

import (
	"auth-svc/internal/auth"
	"auth-svc/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapAuthRoutes(authRoutes fiber.Router, h auth.Handlers, mw *middleware.MDWManager) {
}
