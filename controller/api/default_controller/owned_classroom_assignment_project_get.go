package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomAssignmentProjectResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
} //@Name GetOwnedClassroomAssignmentProjectResponse

// @Summary		GetOwnedClassroomAssignmentProject
// @Description	GetOwnedClassroomAssignmentProject
// @Id				GetOwnedClassroomAssignmentProject
// @Tags			assignment
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Param			projectId		path		string	true	"Project ID"	Format(uuid)
// @Success		200				{object}	default_controller.getOwnedClassroomAssignmentProjectResponse
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/classrooms/owned/{classroomId}/assignments/{assignmentId}/projects/{projectId} [get]
func (ctrl *DefaultController) GetOwnedClassroomAssignmentProject(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	project := ctx.GetOwnedClassroomAssignmentProject()
	response := getOwnedClassroomAssignmentProjectResponse{
		AssignmentProjects: *project,
		ProjectPath:        fmt.Sprintf("/api/v1/classrooms/owned/%s/assignments/%s/projects/%s/gitlab", classroom.ID.String(), project.AssignmentID.String(), project.ID.String()),
	}
	return c.JSON(response)
}
