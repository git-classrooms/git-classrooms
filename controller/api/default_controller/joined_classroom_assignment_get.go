package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getJoinedClassroomAssignmentResponse struct {
	*database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
} //@Name GetJoinedClassroomAssignmentResponse

// @Summary		GetJoinedClassroomAssignment
// @Description	GetJoinedClassroomAssignment
// @Id				GetJoinedClassroomAssignment
// @Tags			assignment
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{object}	default_controller.getJoinedClassroomAssignmentResponse
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v1/classrooms/joined/{classroomId}/assignments/{assignmentId} [get]
func (ctrl *DefaultController) GetJoinedClassroomAssignment(c *fiber.Ctx) error {
	ctx := context.Get(c)
	assignment := ctx.GetJoinedClassroomAssignment()

	repo := ctx.GetGitlabRepository()
	webURL := ""
	if assignment.AssignmentAccepted {
		projectFromGitLab, err := repo.GetProjectById(assignment.ProjectID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		webURL = projectFromGitLab.WebUrl
	}
	response := &getJoinedClassroomAssignmentResponse{
		AssignmentProjects: assignment,
		ProjectPath:        webURL,
	}

	return c.JSON(response)
}
