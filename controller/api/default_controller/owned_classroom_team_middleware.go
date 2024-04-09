package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) OwnedClassroomTeamMiddleware(c *fiber.Ctx) error {
	ctx := context.Get(c)

	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.TeamID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryTeam := query.Team
	team, err := queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Classroom).
		Preload(queryTeam.Member).
		Where(queryTeam.ID.Eq(param.TeamID)).
		Where(queryTeam.ClassroomID.Eq(param.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetOwnedClassroomTeam(team)

	return ctx.Next()
}
