package api

import "github.com/gofiber/fiber/v2"

type Controller interface {
	CreateClassroom(c *fiber.Ctx) error
	CreateAssignment(c *fiber.Ctx) error
	JoinClassroom(c *fiber.Ctx) error
	InviteToClassroom(c *fiber.Ctx) error
}
