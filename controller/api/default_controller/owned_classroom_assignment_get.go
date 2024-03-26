package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassroomAssignment(c *fiber.Ctx) error {
	ctx := context.Get(c)
	assignment := ctx.GetOwnedClassroomAssignment()
	return c.JSON(assignment)
}
