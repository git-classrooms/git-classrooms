package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
)

func (ctrl *DefaultController) GetMe(c *fiber.Ctx) error {
	user, err := session.Get(c).GetUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(user)
}
