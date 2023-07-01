package server

import (
	"server-template/internal/middleware"
	serverTemplateHTTPhHandler "server-template/internal/server-template/handler/delivery/http"
	serverTemplateRedisRepository "server-template/internal/server-template/repository"
	serverTemplateUsecase "server-template/internal/server-template/usecase"

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
	authRepo := serverTemplateRedisRepository.NewServerTeamplateRedisRepo(s.redisClient, s.log, s.cfg)

	// usecases
	authUC := serverTemplateUsecase.NewServerTemplateUC(authRepo, s.log, s.cfg)

	// handlers
	authHandler := serverTemplateHTTPhHandler.NewserverTeamplateHandlers(s.cfg, authUC, s.log)

	// route groups
	apiGroup := app.Group("api")
	serverTemplates := apiGroup.Group("server_templates")

	mw := middleware.NewMDWManager(s.cfg, s.log, authUC)

	// TODO: maybe add some middlewares later
	serverTemplateHTTPhHandler.MapServerTemplateRoutes(serverTemplates, authHandler, mw)

	return nil
}
