package server

import (
	"auth-svc/config"
	"auth-svc/pkg/httperrors"
	"auth-svc/pkg/logger"
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	fiber       *fiber.App
	cfg         *config.Config
	redisClient *redis.Client
	log         logger.Logger
}

func NewServer(cfg *config.Config, redisClient *redis.Client, log logger.Logger) *Server {
	return &Server{
		fiber:       fiber.New(fiber.Config{ErrorHandler: httperrors.Init(cfg, log), DisableStartupMessage: true}),
		cfg:         cfg,
		redisClient: redisClient,
		log:         log,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, span := otel.Tracer("").Start(ctx, "Server.Run")
	if err := s.MapHandlers(ctx, s.fiber); err != nil {
		s.log.Fatalf("Cannot map handlers: %s", err.Error())
	}
	if err := s.fiber.Listen(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.HTTPPort)); err != nil {
		s.log.Fatalf("Error starting Server: ", err)
	}
	span.End()

	return nil
}

func (s *Server) Shutdown() (err error) {
	if err = s.fiber.Shutdown(); err != nil {
		return
	}
	return
}
