package api

import (
	"github.com/gofiber/fiber/v2"
	apiV2 "gitlab.hs-flensburg.de/gitlab-classroom/controller/api_v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) ValidateUserMiddleware(validateFunc apiV2.ValidateUserFunc) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ctx := context.Get(c)
		userClassroom := ctx.GetUserClassroom()

		if !validateFunc(*userClassroom) {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
