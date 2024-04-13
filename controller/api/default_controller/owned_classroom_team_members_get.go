package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Get all Members of the current Team
// @Description	Get all Members of the current Team
// @Tags			team, member
// @Accept			json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{array}		database.UserClassrooms
// @Failure		400			{object}	httputil.HTTPError
// @Failure		401			{object}	httputil.HTTPError
// @Failure		404			{object}	httputil.HTTPError
// @Failure		500			{object}	httputil.HTTPError
// @Router			/classrooms/owned/{classroomId}/teams/{teamId}/members [get]
func (ctrl *DefaultController) GetOwnedClassroomTeamMembers(c *fiber.Ctx) error {
	ctx := context.Get(c)
	team := ctx.GetOwnedClassroomTeam()

	members := team.Member

	return ctx.JSON(members)
}
