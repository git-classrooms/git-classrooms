package api

import (
	"github.com/gofiber/fiber/v2"

	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) RedirectReportGitlab(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	projectID := ctx.GetGitlabProjectID()
	repo := ctx.GetGitlabRepository()

	pipeline, err := repo.GetProjectLatestPipeline(projectID, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Redirect(pipeline.WebURL + "/test_report")
}
