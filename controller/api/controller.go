package api

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetMe(*fiber.Ctx) error
	CreateClassroom(*fiber.Ctx) error
	CreateAssignment(*fiber.Ctx) error
	JoinClassroom(*fiber.Ctx) error
	JoinAssignment(*fiber.Ctx) error
	InviteToClassroom(*fiber.Ctx) error
	InviteToAssignment(*fiber.Ctx) error
}
