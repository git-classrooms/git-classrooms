package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Get all Members of the current Team
// @Description	Get all Members of the current Team
// @Id				GetOwnedClassroomTeamMembers
// @Tags			member
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{array}		UserClassrooms
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/owned/{classroomId}/teams/{teamId}/members [get]
func (ctrl *DefaultController) GetOwnedClassroomTeamMembers(c *fiber.Ctx) error {
	ctx := context.Get(c)
	team := ctx.GetOwnedClassroomTeam()

	members := team.Member
	response := utils.Map(members, func(member *database.UserClassrooms) *getOwnedClassroomMemberResponse {
		return &getOwnedClassroomMemberResponse{
			UserClassrooms: member,
			GitlabURL:      fmt.Sprintf("/api/v1/classrooms/owned/%s/teams/%s/members/%d/gitlab", team.ClassroomID.String(), team.ID.String(), member.UserID),
		}
	})
	return c.JSON(response)
}
