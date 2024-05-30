package default_controller

import (
	"fmt"

	"gitlab.hs-flensburg.de/gitlab-classroom/utils"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomTeamResponse struct {
	database.Team
	UserMember []*database.User `json:"members"`
	GitlabURL  string           `json:"gitlabUrl"`
} //@Name GetOwnedClassroomTeamResponse

// @Summary		Get current Team
// @Description	Get current Team
// @Id				GetOwnedClassroomTeam
// @Tags			team
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{object}	default_controller.getOwnedClassroomTeamResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/owned/{classroomId}/teams/{teamId} [get]
func (ctrl *DefaultController) GetOwnedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	team := ctx.GetOwnedClassroomTeam()
	member := utils.Map(team.Member, func(u *database.UserClassrooms) *database.User {
		return &u.User
	})

	response := &getOwnedClassroomTeamResponse{
		Team:       *team,
		UserMember: member,
		GitlabURL:  fmt.Sprintf("/api/v1/classrooms/owned/%s/teams/%s/gitlab", team.ClassroomID.String(), team.ID.String()),
	}

	return c.JSON(response)
}
