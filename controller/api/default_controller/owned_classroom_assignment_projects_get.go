package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetOwnedClassroomAssignmentProjects
// @Description	GetOwnedClassroomAssignmentProjects
// @Id				GetOwnedClassroomAssignmentProjects
// @Tags			assignment
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{array}		default_controller.getOwnedClassroomAssignmentProjectResponse
// @Failure		400				{object}	httputil.HTTPError
// @Failure		401				{object}	httputil.HTTPError
// @Failure		404				{object}	httputil.HTTPError
// @Failure		500				{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/assignments/{assignmentId}/projects [get]
func (ctrl *DefaultController) GetOwnedClassroomAssignmentProjects(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	assignment := ctx.GetOwnedClassroomAssignment()

	responses := make([]getOwnedClassroomAssignmentProjectResponse, len(assignment.Projects))
	for i, project := range assignment.Projects {
		responses[i] = getOwnedClassroomAssignmentProjectResponse{
			AssignmentProjects: *project,
			ProjectPath:        fmt.Sprintf("/api/v1/classrooms/owned/%s/assignments/%s/projects/%s/gitlab", classroom.ID.String(), project.AssignmentID.String(), project.ID.String()),
		}
	}

	return c.JSON(responses)
}
