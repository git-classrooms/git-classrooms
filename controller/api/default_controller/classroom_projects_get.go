package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomProjects
// @Description	GetClassroomProjects
// @Id				GetClassroomProjects
// @Tags			project
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{array}		api.ProjectResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/projects [get]
func (ctrl *DefaultController) GetClassroomProjects(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	if classroom.TeamID == nil {
		return c.JSON([]*ProjectResponse{})
	}

	projects, err := classroomProjectQuery(c, classroom.ClassroomID, *classroom.TeamID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(projects, func(project *database.AssignmentProjects) *ProjectResponse {
		return &ProjectResponse{
			AssignmentProjects: project,
			WebURL:             fmt.Sprintf("/api/v1/classrooms/%s/projects/%s/gitlab", classroom.ClassroomID, project.ID.String()),
			ReportWebURL:       fmt.Sprintf("/api/v1/classrooms/%s/projects/%s/report/gitlab", ctx.GetUserClassroom().ClassroomID, project.ID.String()),
		}
	})

	return c.JSON(response)
}
