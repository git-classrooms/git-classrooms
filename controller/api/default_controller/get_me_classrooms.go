package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
)

func (ctrl *DefaultController) GetMeClassrooms(c *fiber.Ctx) error {
	repo := context.GetGitlabRepository(c)
	user, err := repo.GetCurrentUser()
	if err != nil {
		return err
	}

	return c.JSON(user)
}
