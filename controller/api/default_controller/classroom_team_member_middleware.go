package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func teamMemberQuery(c *fiber.Ctx, classroomID uuid.UUID, teamID uuid.UUID) query.IUserClassroomsDo {
	queryUserClassrooms := query.UserClassrooms

	return queryUserClassrooms.
		WithContext(c.Context()).
		Preload(queryUserClassrooms.User).
		Where(queryUserClassrooms.ClassroomID.Eq(classroomID)).
		Where(queryUserClassrooms.TeamID.Eq(teamID))
}

func (ctrl *DefaultController) ClassroomTeamMemberMiddleware(c *fiber.Ctx) (err error) {
	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.TeamID == nil || params.MemberID == nil {
		return fiber.ErrBadRequest
	}

	member, err := teamMemberQuery(c, *params.ClassroomID, *params.TeamID).
		Where(query.UserClassrooms.UserID.Eq(*params.MemberID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetClassroomMember(member)
	ctx.SetGitlabUserID(member.UserID)

	return c.Next()
}
