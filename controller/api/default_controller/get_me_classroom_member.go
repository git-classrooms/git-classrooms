package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"strconv"
)

func (ctrl *DefaultController) GetMeClassroomMember(c *fiber.Ctx) error {
	classroom := context.Get(c).GetClassroom()

	memberId, err := strconv.ParseInt(c.Params("memberId"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no member specified")
	}

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
	user, err := queryUser.WithContext(c.Context()).
		Where(queryUser.ID.In(ids...)).
		Where(queryUser.ID.Eq(int(memberId))).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(user)
}
