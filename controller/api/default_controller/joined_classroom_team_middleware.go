package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func joinedClassroomTeamQuery(c *fiber.Ctx, classroomId uuid.UUID) query.ITeamDo {
	queryTeam := query.Team
	return queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Preload(queryTeam.Member.User).
		Where(queryTeam.ClassroomID.Eq(classroomId))
}

func (ctrl *DefaultController) JoinedClassroomTeamMiddleware(c *fiber.Ctx) error {
	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.TeamID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	team, err := joinedClassroomTeamQuery(c, *param.ClassroomID).
		Where(query.Team.ID.Eq(*param.TeamID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetJoinedClassroomTeam(team)
	return ctx.Next()
}
