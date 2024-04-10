package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomTeamProjectResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

func (ctrl *DefaultController) GetOwnedClassroomTeamProjects(c *fiber.Ctx) error {
	ctx := context.Get(c)

	team := ctx.GetOwnedClassroomTeam()

	responses := make([]getOwnedClassroomTeamProjectResponse, len(team.AssignmentProjects))
	for i, project := range team.AssignmentProjects {
		responses[i] = getOwnedClassroomTeamProjectResponse{
			AssignmentProjects: *project,
			ProjectPath:        fmt.Sprintf("/api/v1/classrooms/owned/%s/assignments/%s/projects/%s/gitlab", team.ClassroomID.String(), project.AssignmentID.String(), project.ID.String()),
		}
	}
	return ctx.JSON(responses)
}
