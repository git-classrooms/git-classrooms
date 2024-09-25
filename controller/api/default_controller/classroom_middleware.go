package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen/field"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func userClassroomQuery(ctx *fiber.Ctx, userID int) query.IUserClassroomsDo {
	queryUserClassroom := query.UserClassrooms
	return queryUserClassroom.
		WithContext(ctx.Context()).
		Preload(queryUserClassroom.Classroom).
		Preload(queryUserClassroom.User).
		Preload(queryUserClassroom.Team).
		Preload(field.NewRelation("Classroom.Owner", "")).
		Preload(field.NewRelation("Classroom.Assignments", "")).
		Preload(field.NewRelation("Classroom.ManualGradingRubrics", "")).
		Where(queryUserClassroom.UserID.Eq(userID))
}

func (ctrl *DefaultController) ClassroomMiddleware(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	userID := ctx.GetUserID()

	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil {
		return fiber.ErrBadRequest
	}

	classroom, err := userClassroomQuery(c, userID).
		Where(query.UserClassrooms.ClassroomID.Eq(params.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetUserClassroom(classroom)
	ctx.SetGitlabGroupID(classroom.Classroom.GroupID)

	return c.Next()
}
