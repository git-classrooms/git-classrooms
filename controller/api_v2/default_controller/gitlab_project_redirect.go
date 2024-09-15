package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) RedirectProjectGitlab(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	projectId := ctx.GetGitlabProjectID()
	repo := ctx.GetGitlabRepository()

	project, err := repo.GetProjectById(c.Context(), projectId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Redirect(project.WebUrl)
}
