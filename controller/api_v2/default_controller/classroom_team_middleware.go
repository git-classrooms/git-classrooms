package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func classroomTeamQuery(c *fiber.Ctx, classroomID uuid.UUID) query.ITeamDo {
	queryTeam := query.Team
	return queryTeam.
		WithContext(c.Context()).
		Where(queryTeam.ClassroomID.Eq(classroomID))
}

func (ctrl *DefaultController) ClassroomTeamMiddleware(c *fiber.Ctx) (err error) {
	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.TeamID == nil {
		return fiber.ErrBadRequest
	}

	team, err := classroomTeamQuery(c, *params.ClassroomID).
		Where(query.Team.ID.Eq(*params.TeamID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetTeam(team)
	ctx.SetGitlabGroupID(team.GroupID)

	return c.Next()
}
