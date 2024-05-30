package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

// @Summary		GetInvitationInfo
// @Description	GetInvitationInfo
// @Id				GetInvitationInfo
// @Tags			classroom
// @Produce		json
// @Param			invitationId	path		string	true	"Invitation ID"	Format(uuid)
// @Success		200				{object}	database.ClassroomInvitation
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v1/invitations/{invitationId} [get]
func (ctrl *DefaultController) GetInvitationInfo(c *fiber.Ctx) error {
	param := &Params{}
	err := c.ParamsParser(param)
	if err != nil || param.InvitationID == nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryClassroomInvitation := query.ClassroomInvitation
	invitation, err := queryClassroomInvitation.
		WithContext(c.Context()).
		Preload(queryClassroomInvitation.Classroom).
		Preload(queryClassroomInvitation.Classroom.Owner).
		Where(queryClassroomInvitation.ID.Eq(*param.InvitationID)).
		Where(queryClassroomInvitation.Status.Eq(uint8(database.ClassroomInvitationPending))).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(invitation)
}
