package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) UpdateClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	return c.JSON(classroom)
}
