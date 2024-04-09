package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomTeamResponse struct {
	database.Team
	GitlabUrl string `json:"gitlabUrl"`
}

func (ctrl *DefaultController) GetOwnedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	team := ctx.GetOwnedClassroomTeam()

	response := &getOwnedClassroomTeamResponse{
		Team:      *team,
		GitlabUrl: fmt.Sprintf("/api/v1/classrooms/owned/%s/teams/%s/gitlab", team.ClassroomID.String(), team.ID.String()),
	}

	return c.JSON(response)
}
