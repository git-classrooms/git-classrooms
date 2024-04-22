package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassroomTemplates(c *fiber.Ctx) error {
	type templateResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	ctx := context.Get(c)
	search := c.Query("search")

	repo := ctx.GetGitlabRepository()
	projects, err := repo.GetAllProjects(search)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	response := utils.Map(projects, func(project *model.Project) *templateResponse {
		return &templateResponse{
			ID:   project.ID,
			Name: project.Name,
		}
	})

	return c.JSON(response)
}
