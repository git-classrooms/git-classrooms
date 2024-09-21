package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTeamMembers
// @Description	GetClassroomTeamMembers
// @Id				GetClassroomTeamMembers
// @Tags			member
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{array}		api.UserClassroomResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/teams/{teamId}/members [get]
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
			UserClassrooms:   member,
			WebURL:           fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/members/%d/gitlab", classroom.ClassroomID.String(), team.ID.String(), member.UserID),
			AssignmentsCount: len(classroom.Classroom.Assignments),
		}
	})

	return c.JSON(response)
}
