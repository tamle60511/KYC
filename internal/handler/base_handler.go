package handler

import "github.com/gofiber/fiber/v3"

type BaseHandler interface {
	SetupRoutes(router fiber.Router, ms ...fiber.Handler)
}
