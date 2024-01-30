package api

import "github.com/gofiber/fiber/v2"

type Controller interface {
	CreateClassroom(c *fiber.Ctx) error
	CreateAssignment(c *fiber.Ctx) error
}
