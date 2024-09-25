package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) RedirectUserGitlab(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	userID := ctx.GetGitlabUserID()
	repo := ctx.GetGitlabRepository()

	user, err := repo.GetUserById(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Redirect(user.WebUrl)
}
