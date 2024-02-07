package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func (ctrl *DefaultController) GetMeClassroomInvitations(c *fiber.Ctx) error {
	classroom := context.Get(c).GetClassroom()
	if classroom.Role != database.Owner {
		return fiber.NewError(fiber.StatusForbidden, "only the owner can access the invitations")
	}

	queryClassroomInvitation := query.ClassroomInvitation
	invitations, err := queryClassroomInvitation.
		WithContext(c.Context()).
		Where(queryClassroomInvitation.ClassroomID.Eq(classroom.ClassroomID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(invitations)
}
