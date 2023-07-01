package auth

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	CreateAdminWsToken() fiber.Handler
	CreateBrokerWsToken() fiber.Handler
}
