package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTeamMember
// @Description	GetClassroomTeamMember
// @Id				GetClassroomTeamMember
// @Tags			member
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Param			memberId		path		int	true	"Member ID"
// @Success		200			{object}	api.UserClassroomResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/teams/{teamId}/members/{memberId} [get]
func (ctrl *DefaultController) GetClassroomTeamMember(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()
	member := ctx.GetClassroomMember()

	response := &UserClassroomResponse{
		UserClassrooms: member,
		WebURL:         fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/members/%d/gitlab", classroom.ClassroomID.String(), team.ID.String(), member.UserID),
	}

	return c.JSON(response)
}
