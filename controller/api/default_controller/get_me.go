package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetMe(c *fiber.Ctx) error {

	// TODO: Add Avatar-URL and json tags to the User struct
	gitlabUser, err := context.Get(c).GetGitlabRepository().GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(gitlabUser)
}
