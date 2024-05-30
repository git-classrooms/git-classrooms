package default_controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getJoinedClassroomTeamResponse struct {
	database.Team
	UserMember []*database.User `json:"members"`
	GitlabURL  string           `json:"gitlabUrl"`
} //@Name GetJoinedClassroomTeamResponse

// @Summary		Get current Team
// @Description	Get current Team
// @Id				GetJoinedClassroomTeam
// @Tags			team
// @Produces		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{object}	default_controller.getJoinedClassroomTeamResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/joined/{classroomId}/teams/{teamId} [get]
func (ctrl *DefaultController) GetJoinedClassroomTeam(c *fiber.Ctx) error {
	ctx := context.Get(c)
	team := ctx.GetJoinedClassroom().Team

	member := utils.Map(team.Member, func(u *database.UserClassrooms) *database.User {
		return &u.User
	})

	log.Println("GetJoinedClassroomTeam")

	response := &getJoinedClassroomTeamResponse{
		Team:       *team,
		UserMember: member,
		GitlabURL:  fmt.Sprintf("/api/v1/classrooms/joined/%s/teams/%s/gitlab", team.ClassroomID.String(), team.ID.String()),
	}

	return c.JSON(response)
}
