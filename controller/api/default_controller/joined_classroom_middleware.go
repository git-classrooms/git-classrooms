package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/gen/field"
)

func joinedClassroomQuery(userID int, c *fiber.Ctx) query.IUserClassroomsDo {
	queryUserClassrooms := query.UserClassrooms
	return queryUserClassrooms.
		WithContext(c.Context()).
		Preload(queryUserClassrooms.Classroom).
		Preload(queryUserClassrooms.Team).
		Preload(field.NewRelation("Classroom.Owner", "")).
		Where(queryUserClassrooms.UserID.Eq(userID))
}

func (ctrl *DefaultController) JoinedClassroomMiddleware(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()

	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryUserClassrooms := query.UserClassrooms
	classroom, err := joinedClassroomQuery(userID, c).
		Where(queryUserClassrooms.ClassroomID.Eq(param.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetJoinedClassroom(classroom)
	ctx.SetGitlabGroupID(classroom.Classroom.GroupID)
	return ctx.Next()
}
