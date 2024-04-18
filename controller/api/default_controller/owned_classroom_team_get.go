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

// @Summary		Get current Team
// @Description	Get current Team
// @Id				GetOwnedClassroomTeam
// @Tags			team
// @Accept			json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{object}	default_controller.getOwnedClassroomTeamResponse
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		404			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/teams/{teamId} [get]
func (ctrl *DefaultController) GetOwnedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	team := ctx.GetOwnedClassroomTeam()

	response := &getOwnedClassroomTeamResponse{
		Team:      *team,
		GitlabUrl: fmt.Sprintf("/api/v1/classrooms/owned/%s/teams/%s/gitlab", team.ClassroomID.String(), team.ID.String()),
	}

	return c.JSON(response)
}
