package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) ArchivedMiddleware(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	switch c.Method() {
	case fiber.MethodPost:
		fallthrough
	case fiber.MethodPut:
		fallthrough
	case fiber.MethodPatch:
		fallthrough
	case fiber.MethodDelete:
		if classroom.Classroom.Archived {
			return fiber.NewError(fiber.StatusForbidden, "Classroom is archived")
		}
	default:
	}

	return c.Next()
}
