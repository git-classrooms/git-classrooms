package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassroomTemplates(c *fiber.Ctx) error {
	ctx := context.Get(c)
	search := c.Query("search")

	repo := ctx.GetGitlabRepository()
	projects, err := repo.GetAllProjects(search)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(projects)
}
