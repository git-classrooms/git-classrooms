package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getJoinedClassroomAssignmentResponse struct {
	*database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

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
