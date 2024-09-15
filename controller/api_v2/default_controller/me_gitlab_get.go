package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetMeGitlab(c *fiber.Ctx) (err error) {
	gitlabUser, err := context.Get(c).GetGitlabRepository().GetCurrentUser(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Redirect(gitlabUser.WebUrl)
}
