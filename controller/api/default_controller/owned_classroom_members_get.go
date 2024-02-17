package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassroomMembers(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()

	ids := make([]int, len(classroom.Member))
	for i, member := range classroom.Member {
		ids[i] = member.UserID
	}

	queryUser := query.User
	users, err := queryUser.
		WithContext(c.Context()).
		Where(queryUser.ID.In(ids...)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(users)
}
