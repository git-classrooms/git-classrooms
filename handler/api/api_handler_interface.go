package apiHandler

import "github.com/gofiber/fiber/v2"

type ApiHandler interface {
	CreateClassroom(c *fiber.Ctx) error
	CreateAssignment(c *fiber.Ctx) error
}
