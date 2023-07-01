package http

import (
	"auth-svc/config"
	"auth-svc/internal/auth"
	customErrors "auth-svc/pkg/errors"
	"auth-svc/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	log    logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, log: log}
}

func (h *authHandlers) CreateAdminWsToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := otel.Tracer("").Start(c.Context(), "AuthHandlers.GetUUID")
		defer span.End()

		userID, ok := c.Locals("userID").(int)
		if !ok {
			err := customErrors.ErrGetUserIDFromCtx
			h.log.Errorf("Error: %v\nauth.delivery.http.handlers.GetUUID() cannot get userID from fiber ctx", err)
			return err
		}

		result, err := h.authUC.CreateWSToken(ctx, userID, auth.AdminSessionTypeID)
		if err != nil {
			h.log.Error(err)
			return err
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

func (h *authHandlers) CreateBrokerWsToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := otel.Tracer("").Start(c.Context(), "AuthHandlers.GetUUID")
		defer span.End()

		userID, ok := c.Locals("userID").(int)
		if !ok {
			err := customErrors.ErrGetUserIDFromCtx
			h.log.Errorf("Error: %v\nauth.delivery.http.handlers.GetUUID() cannot get userID from fiber ctx", err)
			return err
		}

		result, err := h.authUC.CreateWSToken(ctx, userID, auth.BrokerSessionTypeID)
		if err != nil {
			h.log.Error(err)
			return err
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
