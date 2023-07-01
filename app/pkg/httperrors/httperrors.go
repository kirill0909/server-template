package httperrors

import (
	"auth-svc/config"
	"auth-svc/pkg/errors"
	"auth-svc/pkg/logger"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type responseMsg struct {
	Error string `json:"error"`
}

func Init(config *config.Config, logger logger.Logger) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		var (
			response   responseMsg
			statusCode int
		)

		if err, ok := err.(*errors.CustomError); ok {
			if err.IsNeedToBeShown() || config.Server.ShowUnknownErrorsInResponse {
				statusCode = err.GetStatusCode()
				response.Error = err.Error()
				return c.Status(statusCode).JSON(response)
			}
			statusCode = 500
			response.Error = "Some error"
			return c.Status(statusCode).JSON(response)
		}

		statusCode = 500
		if errStatusCode := c.Context().Response.StatusCode(); errStatusCode >= 400 {
			statusCode = errStatusCode
		}

		if config.Server.ShowUnknownErrorsInResponse {
			response.Error = fmt.Sprintf("[DEBUG MODE] %s", err.Error())
		} else if response.Error == "" {
			logger.Error(fmt.Errorf("%s %v", c.OriginalURL(), err))
			response.Error = "unknown error"
		}

		return c.Status(statusCode).JSON(response)
	}
}
