package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) RedirectGroupGitlab(c *fiber.Ctx) error {
	ctx := context.Get(c)
	groupID := ctx.GetGitlabGroupID()
	repo := ctx.GetGitlabRepository()

	group, err := repo.GetGroupById(groupID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Redirect(group.WebUrl)
}
