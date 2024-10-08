package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomAssignmentProject
// @Description	GetClassroomAssignmentProject
// @Id				GetClassroomAssignmentProject
// @Tags			project
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Param			projectId		path		string	true	"Project ID"	Format(uuid)
// @Success		200				{object}	api.ProjectResponse
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		403				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentId}/projects/{projectId} [get]
func (ctrl *DefaultController) GetClassroomAssignmentProject(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	assignment := ctx.GetAssignment()
	project := ctx.GetAssignmentProject()

	response := &ProjectResponse{
		AssignmentProjects: project,
		WebURL:             fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s/projects/%s/gitlab", classroom.ClassroomID.String(), assignment.ID.String(), project.ID.String()),
		ReportWebURL:       fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s/projects/%s/report/gitlab", classroom.ClassroomID.String(), assignment.ID.String(), project.ID.String()),
	}

	return c.JSON(response)
}
