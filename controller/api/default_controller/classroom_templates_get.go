package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTemplates
// @Description	GetClassroomTemplates
// @Id				GetClassroomTemplates
// @Tags			classroom
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		TemplateResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/templateProjects [get]
func (ctrl *DefaultController) GetClassroomTemplates(c *fiber.Ctx) error {
	type templateResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} //@Name TemplateResponse

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
