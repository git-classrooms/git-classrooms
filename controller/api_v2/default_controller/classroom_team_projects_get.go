package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomTeamProjects(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	projects, err := teamProjectQuery(c, team.ID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(projects, func(project *database.AssignmentProjects) *ProjectResponse {
		return &ProjectResponse{
			AssignmentProjects: project,
			WebURL:             fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/projects/%s", classroom.ClassroomID.String(), team.ID.String(), project.ID.String()),
		}
	})

	return c.JSON(response)
}
