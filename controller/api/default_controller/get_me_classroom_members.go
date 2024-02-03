package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func (ctrl *DefaultController) GetMeClassroomMembers(c *fiber.Ctx) error {
	classroom := context.GetClassroom(c)

	queryUserClassrooms := query.UserClassrooms
	fetchedMembers, err := queryUserClassrooms.
		WithContext(c.Context()).
		Preload(queryUserClassrooms.User).
		Where(queryUserClassrooms.ClassroomID.Eq(classroom.ClassroomID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ids := make([]int, len(fetchedMembers))
	for i, member := range fetchedMembers {
		ids[i] = member.UserID
	}

	queryUser := query.User
	users, err := queryUser.WithContext(c.Context()).Where(queryUser.ID.In(ids...)).Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(users)
}
