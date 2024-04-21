package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"strconv"
)

func (ctrl *DefaultController) GetUserGitlab(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))
	gitlabUser, err := context.Get(c).GetGitlabRepository().GetUserById(userId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect(gitlabUser.WebUrl)
}
