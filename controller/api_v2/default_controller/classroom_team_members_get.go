package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomTeamMembers(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	members, err := teamMemberQuery(c, classroom.ClassroomID, team.ID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(members, func(member *database.UserClassrooms) *UserClassroomResponse {
		return &UserClassroomResponse{
			UserClassrooms: member,
			WebURL:         fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/members/%d", classroom.ClassroomID.String(), team.ID.String(), member.UserID),
		}
	})

	return c.JSON(response)
}
