package server_template

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	Ping() fiber.Handler
}
