package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

// @Summary		GetClassroomInvitation
// @Description	GetClassroomInvitation
// @Id				GetClassroomInvitation
// @Tags			classroom
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			invitationId	path		string	true	"Invitation ID"	Format(uuid)
// @Success		200				{object}	ClassroomInvitation
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		403				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/invitations/{invitationId} [get]
func (ctrl *DefaultController) GetClassroomInvitation(c *fiber.Ctx) (err error) {
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
		Preload(queryInvitation.Classroom).
		Preload(queryInvitation.Classroom.Owner).
		Where(queryInvitation.ClassroomID.Eq(*params.ClassroomID)).
		Where(queryInvitation.ID.Eq(*params.InvitationID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(invitation)
}
