package server

import (
	authHTTPhHandler "server-template/internal/auth/handler/delivery/http"
	authRedisRepository "server-template/internal/auth/repository"
	authUsecase "server-template/internal/auth/usecase"
	"server-template/internal/middleware"

	"context"
	"fmt"
	"server-template/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.opentelemetry.io/otel"
)

func (s *Server) MapHandlers(ctx context.Context, app *fiber.App) error {
	_, span := otel.Tracer("").Start(ctx, "Server.MapHandlers")
	defer span.End()

	// health check
	s.fiber.Use(func(c *fiber.Ctx) error {
		if c.OriginalURL() == "/health_check" {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.Next()
	})

	// fiber logger
	loggerCfg := logger.ConfigDefault
	loggerCfg.Format = fmt.Sprintf(
		"${time} | ${status} | ${latency} | ${method} | ${path} | ${respHeader:%s}\n",
		utils.TraceIDHeader,
	)
	app.Use(logger.New(loggerCfg))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://broker.aifory.io,https://admin.aifory.io",
		AllowHeaders:     "Accept,Accept-Language,Content-Language,Content-Type,fingerprint,User-Agent",
		AllowCredentials: true,
	}))

	// repos
	authRepo := authRedisRepository.NewAuthRedisRepo(s.redisClient, s.log, s.cfg)

	// usecases
	authUC := authUsecase.NewAuthUC(authRepo, s.log, s.cfg)

	// handlers
	authHandler := authHTTPhHandler.NewAuthHandlers(s.cfg, authUC, s.log)

	// route groups
	apiGroup := app.Group("api")
	authGroup := apiGroup.Group("auth")

	mw := middleware.NewMDWManager(s.cfg, s.log, authUC)

	// TODO: maybe add some middlewares later
	authHTTPhHandler.MapAuthRoutes(authGroup, authHandler, mw)

	return nil
}
