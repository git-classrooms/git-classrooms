package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func (ctrl *DefaultController) RevokeClassroomInvitation(c *fiber.Ctx) (err error) {
	var params Params

	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.InvitationID == nil {
		return fiber.ErrBadRequest
	}

	queryInvitation := query.ClassroomInvitation
	invitation, err := queryInvitation.
		WithContext(c.Context()).
		Where(queryInvitation.ClassroomID.Eq(*params.ClassroomID)).
		Where(queryInvitation.ID.Eq(*params.InvitationID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch invitation.Status {
	case database.ClassroomInvitationRevoked:
		return c.SendStatus(fiber.StatusNoContent)
	case database.ClassroomInvitationAccepted:
		return fiber.NewError(fiber.StatusBadRequest, "Cannot revoke accepted invitation")
	case database.ClassroomInvitationRejected:
		fallthrough
	case database.ClassroomInvitationPending:
		invitation.Status = database.ClassroomInvitationRevoked
	}

	_, err = queryInvitation.WithContext(c.Context()).Updates(invitation)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
