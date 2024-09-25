package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func classroomMemberQuery(c *fiber.Ctx, classroomID uuid.UUID) query.IUserClassroomsDo {
	queryUserClassrooms := query.UserClassrooms
	return queryUserClassrooms.
		WithContext(c.Context()).
		Preload(queryUserClassrooms.User).
		Preload(queryUserClassrooms.Team).
		Where(queryUserClassrooms.ClassroomID.Eq(classroomID))
}

func (ctrl *DefaultController) ClassroomMemberMiddleware(c *fiber.Ctx) (err error) {
	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.MemberID == nil {
		return fiber.ErrBadRequest
	}

	member, err := classroomMemberQuery(c, *params.ClassroomID).
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
