package auth

import "github.com/gofiber/fiber/v2"

type Controller interface {
	AuthMiddleware(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	Callback(c *fiber.Ctx) error
	GetAuth(c *fiber.Ctx) error
}
