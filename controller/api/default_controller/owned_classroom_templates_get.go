package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetOwnedClassroomTemplates
// @Description	GetOwnedClassroomTemplates
// @Id				GetOwnedClassroomTemplates
// @Tags			classroom
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		default_controller.GetOwnedClassroomTemplates.templateResponse
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		403			{object}	httputil.HTTPError
// @Failure		404			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/templates [get]
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
