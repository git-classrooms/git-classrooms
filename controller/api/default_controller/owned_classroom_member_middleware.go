package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func ownedClassroomMemberQuery(classroomID uuid.UUID, c *fiber.Ctx) query.IUserClassroomsDo {
	queryUserClassroom := query.UserClassrooms
	return queryUserClassroom.
		WithContext(c.Context()).
		Preload(queryUserClassroom.User).
		Preload(queryUserClassroom.User.GitLabAvatar).
		Preload(queryUserClassroom.Team).
		Where(queryUserClassroom.ClassroomID.Eq(classroomID))
}

func (ctrl *DefaultController) OwnedClassroomMemberMiddleware(c *fiber.Ctx) error {
	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil || param.MemberID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	queryUserClassroom := query.UserClassrooms
	user, err := ownedClassroomMemberQuery(classroom.ID, c).
		Where(queryUserClassroom.UserID.Eq(*param.MemberID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetOwnedClassroomMember(user)
	ctx.SetGitlabUserID(user.UserID)
	return ctx.Next()
}
