package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassroomAssignmentProjectGitlab(c *fiber.Ctx) error {
	ctx := context.Get(c)
	project := ctx.GetOwnedClassroomAssignmentProject()
	repo := ctx.GetGitlabRepository()

	projectFromGitLab, err := repo.GetProjectById(project.ProjectID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect(projectFromGitLab.WebUrl)
}
