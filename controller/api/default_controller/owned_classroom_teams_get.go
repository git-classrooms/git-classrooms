package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassroomTeams(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	response := utils.Map(classroom.Teams, func(team *database.Team) *getOwnedClassroomTeamResponse {
		return &getOwnedClassroomTeamResponse{
			Team:      *team,
			GitlabUrl: fmt.Sprintf("/api/v1/classrooms/owned/%s/teams/%s/gitlab", classroom.ID.String(), team.ID.String()),
		}
	})

	return c.JSON(response)
}
