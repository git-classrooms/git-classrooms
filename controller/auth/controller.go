package auth

import "github.com/gofiber/fiber/v2"

type Controller interface {
	AuthMiddleware(c *fiber.Ctx) error
	Auth(c *fiber.Ctx) error
	Callback(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}
