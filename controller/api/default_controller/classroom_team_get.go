package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTeam
// @Description	GetClassroomTeam
// @Id				GetClassroomTeam
// @Tags			team
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{object}	api.TeamResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/teams/{teamId} [get]
func (ctrl *DefaultController) GetClassroomTeam(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	members := utils.Map(team.Member, func(member *database.UserClassrooms) *UserClassroomResponse {
		return &UserClassroomResponse{
			UserClassrooms:   member,
			WebURL:           fmt.Sprintf("/api/v1/classrooms/%s/users/%d", classroom.ClassroomID.String(), member.UserID),
			AssignmentsCount: len(classroom.Classroom.Assignments),
		}
	})

	response := &TeamResponse{
		Team:    team,
		Members: members,
		WebURL:  fmt.Sprintf("/api/v1/classrooms/%s/teams/%s/gitlab", classroom.ClassroomID.String(), team.ID.String()),
	}

	return c.JSON(response)
}
