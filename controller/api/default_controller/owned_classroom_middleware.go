package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func ownedClassroomQuery(userID int, c *fiber.Ctx) query.IClassroomDo {
	queryClassroom := query.Classroom
	return queryClassroom.
		WithContext(c.Context()).
		Preload(queryClassroom.Owner).
		Where(queryClassroom.OwnerID.Eq(userID))
}

func (ctrl *DefaultController) OwnedClassroomMiddleware(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()

	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.ClassroomID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryClassroom := query.Classroom
	classroom, err := ownedClassroomQuery(userID, c).
		Preload(queryClassroom.Member).
		Preload(queryClassroom.Invitations).
		Preload(queryClassroom.Teams).
		Where(queryClassroom.ID.Eq(param.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ctx.SetOwnedClassroom(classroom)
	return ctx.Next()
}
