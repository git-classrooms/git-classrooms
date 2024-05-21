package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomTeamMember(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()
	member := ctx.GetClassroomMember()

	response := &UserClassroomResponse{
		UserClassrooms: member,
		WebURL:         fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/members/%d", classroom.ClassroomID.String(), team.ID.String(), member.UserID),
	}

	return c.JSON(response)
}
