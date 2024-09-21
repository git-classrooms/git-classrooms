package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomProject
// @Description	GetClassroomProject
// @Id				GetClassroomProject
// @Tags			project
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			projectId	path		string	true	"Project ID"	Format(uuid)
// @Success		200			{object}	api.ProjectResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/projects/{projectId} [get]
func (ctrl *DefaultController) GetClassroomProject(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	project := ctx.GetAssignmentProject()

	response := &ProjectResponse{
		AssignmentProjects: project,
		WebURL:             fmt.Sprintf("/api/v2/classrooms/%s/projects/%s/gitlab", ctx.GetUserClassroom().ClassroomID, project.ID.String()),
		ReportWebURL:       fmt.Sprintf("/api/v2/classrooms/%s/projects/%s/report/gitlab", ctx.GetUserClassroom().ClassroomID, project.ID.String()),
	}

	return c.JSON(response)
}
