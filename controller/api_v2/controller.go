package apiV2

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

type Controller interface {
	RoleMiddleware(validRoles ...database.Role) fiber.Handler

	GetClassrooms(c *fiber.Ctx) (err error)
	ClassroomMiddleware(c *fiber.Ctx) (err error)
	GetClassroom(c *fiber.Ctx) (err error)
}
