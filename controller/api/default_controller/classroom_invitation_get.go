package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomInvitation
// @Description	GetClassroomInvitation
// @Id				GetClassroomInvitation
// @Tags			classroom
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			invitationId	path		string	true	"Invitation ID"	Format(uuid)
// @Param			groupLink		query		boolean	false	"is Group Link"
// @Success		200				{object}	ClassroomInvitation
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		403				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/invitations/{invitationId} [get]
func (ctrl *DefaultController) GetClassroomInvitation(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	var params Params

	userID := ctx.GetUserID()
	queryUser := query.User
	user, err := queryUser.WithContext(c.Context()).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil || params.InvitationID == nil {
		return fiber.ErrBadRequest
	}

	groupLink := c.QueryBool("groupLink", false)
	var invitation *database.ClassroomInvitation
	if groupLink {
		queryClassroom := query.Classroom
		classroom, err := queryClassroom.WithContext(c.Context()).
			Preload(queryClassroom.Owner).
			Where(queryClassroom.ID.Eq((*params.ClassroomID))).
			First()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		if classroom.InviteCode != *params.InvitationID {
			return fiber.NewError(fiber.StatusNotFound, "invitation code revoked or not found")
		}

		status := database.ClassroomInvitationPending
		queryUserClassrooms := query.UserClassrooms
		_, err = queryUserClassrooms.WithContext(c.Context()).
			Where(queryUserClassrooms.UserID.Eq(userID)).
			Where(queryUserClassrooms.ClassroomID.Eq(*params.ClassroomID)).
			First()
		if err == nil {
			status = database.ClassroomInvitationAccepted
		}

		invitation = &database.ClassroomInvitation{
			ID:          classroom.InviteCode,
			ClassroomID: *params.ClassroomID,
			Classroom:   *classroom,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Status:      status,
			Email:       user.GitlabEmail,
			ExpiryDate:  time.Now().AddDate(0, 0, 14),
			Role:        database.Student,
		}
	} else {
		queryInvitation := query.ClassroomInvitation
		invitation, err = queryInvitation.
			WithContext(c.Context()).
			Preload(queryInvitation.Classroom).
			Preload(queryInvitation.Classroom.Owner).
			Where(queryInvitation.ClassroomID.Eq(*params.ClassroomID)).
			Where(queryInvitation.ID.Eq(*params.InvitationID)).
			First()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}

	return c.JSON(invitation)
}
