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
		Preload(queryUserClassrooms.Classroom).
		Where(queryUserClassrooms.ClassroomID.Eq(param.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	queryTeam := query.Team
	team, err := queryTeam.
		WithContext(c.Context()).
		FindByUserIDAndClassroomID(userID, classroom.ClassroomID)
	if err != nil {
		team = nil
		// The user is not a member of a team
		// log.Println(err)
		// return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetJoinedTeam(team)
	ctx.SetJoinedClassroom(classroom)
	return ctx.Next()
}
