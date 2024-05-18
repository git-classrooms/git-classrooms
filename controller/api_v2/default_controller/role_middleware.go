package api

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) RoleMiddleware(validRoles ...database.Role) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ctx := context.Get(c)
		classroom := ctx.GetUserClassroom()

		if !slices.Contains(validRoles, classroom.Role) {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
