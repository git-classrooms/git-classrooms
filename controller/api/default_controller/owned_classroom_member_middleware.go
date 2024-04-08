package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) OwnedClassroomMemberMiddleware(c *fiber.Ctx) error {
	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.MemberID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryUserClassroom := query.UserClassrooms
	user, err := queryUserClassroom.
		WithContext(c.Context()).
		Preload(queryUserClassroom.Team).
		Where(queryUserClassroom.UserID.Eq(*param.MemberID)).
		Where(queryUserClassroom.ClassroomID.Eq(*param.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx := context.Get(c)
	ctx.SetOwnedClassroomMember(user)

	return ctx.Next()
}
