package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomAssignmentProjectsResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

func (ctrl *DefaultController) GetOwnedClassroomAssignmentProjects(c *fiber.Ctx) error {
	ctx := context.Get(c)
	assignment := ctx.GetOwnedClassroomAssignment()
	repo := ctx.GetGitlabRepository()

	responses := make([]getOwnedClassroomAssignmentProjectsResponse, len(assignment.Projects))
	for i, project := range assignment.Projects {
		webURL := ""
		if project.AssignmentAccepted {
			projectFromGitLab, err := repo.GetProjectById(project.ProjectID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			webURL = projectFromGitLab.WebUrl
		}
		responses[i] = getOwnedClassroomAssignmentProjectsResponse{
			AssignmentProjects: *project,
			ProjectPath:        webURL,
		}
	}

	return c.JSON(responses)
}
