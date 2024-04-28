package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func ownedClassroomTeamQuery(c *fiber.Ctx, classroomId uuid.UUID) query.ITeamDo {
	queryTeam := query.Team
	return queryTeam.
		WithContext(c.Context()).
		Preload(queryTeam.Member).
		Preload(queryTeam.Member.User).
		Where(queryTeam.ClassroomID.Eq(classroomId))
}

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
		Preload(queryTeam.Member.User).
		Preload(queryTeam.AssignmentProjects).
		Where(queryTeam.ID.Eq(param.TeamID)).
		Where(queryTeam.ClassroomID.Eq(param.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetOwnedClassroomTeam(team)
	ctx.SetGitlabGroupID(team.GroupID)
	return ctx.Next()
}
