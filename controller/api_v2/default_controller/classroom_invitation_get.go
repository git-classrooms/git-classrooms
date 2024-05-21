package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func (ctrl *DefaultController) GetClassroomInvitation(c *fiber.Ctx) (err error) {
	var params Params

	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.InvitationID == nil {
		return fiber.ErrBadRequest
	}

	queryInvitation := query.ClassroomInvitation
	invitations, err := queryInvitation.
		WithContext(c.Context()).
		Where(queryInvitation.ClassroomID.Eq(*params.ClassroomID)).
		Where(queryInvitation.ID.Eq(*params.InvitationID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(invitations)
}
