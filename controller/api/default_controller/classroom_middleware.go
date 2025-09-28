package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/gen/field"
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
	cleanLogger := ctx.GetLogger()
	log := ctx.GetLoggerForHandler("ClassroomMiddleware")
	userID := ctx.GetUserID()

	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil {
		return fiber.ErrBadRequest
	}

	log.Debug("retrieving userClassroom from db", "classroomID", *params.ClassroomID)

	classroom, err := userClassroomQuery(c, userID).
		Where(query.UserClassrooms.ClassroomID.Eq(params.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetLogger(cleanLogger.With(
		"classroom", classroom.Classroom,
		"role", classroom.Role.String(),
		"teamID", classroom.TeamID,
	))

	ctx.SetUserClassroom(classroom)
	ctx.SetGitlabGroupID(classroom.Classroom.GroupID)

	return c.Next()
}
