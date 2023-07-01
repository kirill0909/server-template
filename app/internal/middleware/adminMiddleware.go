package middleware

import (
	"auth-svc/internal/auth"
	models "auth-svc/internal/models/auth"
	"auth-svc/pkg/utils"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
)

func (mw *MDWManager) AdminAuthedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := utils.StartFiberTrace(c, "MDWManager.AdminAuthedMiddleware")
		defer span.End()

		authHeaders := models.AuthHeaders{
			IP:            c.Get(mw.cfg.Server.IPHeader),
			UserAgent:     c.Get("User-Agent"),
			Fingerprint:   c.Get("Fingerprint"),
			Authorization: c.Cookies(AdminAccessTokenCookieName),
		}

		if err := authHeaders.Validate(); err != nil {
			return err
		}

		session, err := mw.authUC.GetSession(ctx, authHeaders.Authorization, auth.AdminSessionTypeID)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				mw.log.Info("MDWManager.AdminAuthedMiddleware. Unable to find session session")
				c.Status(fiber.StatusUnauthorized)
				err = nil
			} else {
				mw.log.Error(err)
				c.Status(fiber.StatusInternalServerError)
			}
			return err
		}
		// TODO: Убрать? Тк редис автоматом удалит запись при истечении TTL
		if session.ExpireAt < int(time.Now().Unix()) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userID, err := strconv.Atoi(session.ClientID)
		if err != nil {
			mw.log.Error(err)
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		mw.log.Infof("sess: %s", utils.StructToPrettyJsonString(session))
		mw.log.Infof("authHeaders: %s", utils.StructToPrettyJsonString(authHeaders))

		if !session.IsAuthDataValid(authHeaders) {
			c.Status(fiber.StatusUnauthorized)
			return err
		}

		authHeadersBytes, err := json.Marshal(authHeaders)
		if err != nil {
			mw.log.Error(err)
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		span.SetAttributes(attribute.Int("userID", userID))
		span.SetAttributes(attribute.String("authHeaders", string(authHeadersBytes)))
		c.Locals("userID", userID)
		c.Locals("authHeaders", authHeaders)
		c.Locals("traceCtx", ctx)
		return c.Next()
	}
}
