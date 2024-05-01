package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) OwnedClassroomTeamMemberMiddleware(c *fiber.Ctx) error {
	ctx := context.Get(c)

	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.TeamID == nil || param.MemberID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryUserClassrooms := query.UserClassrooms
	member, err := queryUserClassrooms.
		WithContext(c.Context()).
		Where(queryUserClassrooms.TeamID.Eq(param.TeamID)).
		Where(queryUserClassrooms.ClassroomID.Eq(param.ClassroomID)).
		Where(queryUserClassrooms.UserID.Eq(*param.MemberID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetOwnedClassroomTeamMember(member)
	ctx.SetGitlabUserID(member.UserID)

	return ctx.Next()
}
