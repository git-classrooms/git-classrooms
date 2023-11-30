package handler

import "github.com/gofiber/fiber/v2"

type Handler interface {
	CreateClassroom(c *fiber.Ctx) error
}
