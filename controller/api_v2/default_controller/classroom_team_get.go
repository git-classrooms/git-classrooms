package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
// @Router			/api/v2/classrooms/{classroomId}/teams/{teamId} [get]
func (ctrl *DefaultController) GetClassroomTeam(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	response := &TeamResponse{
		Team:   team,
		WebURL: fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/gitlab", classroom.ClassroomID.String(), team.ID.String()),
	}

	return c.JSON(response)
}
