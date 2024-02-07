package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
)

func (ctrl *DefaultController) GetMeClassroom(c *fiber.Ctx) error {
	classroom := context.Get(c).GetClassroom()

	return c.JSON(classroom)
}
