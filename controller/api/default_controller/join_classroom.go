package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"time"
)

func (ctrl *DefaultController) JoinClassroom(c *fiber.Ctx) error {
	invitationIdParameter := c.Params("invitationId")
	classroomIdParameter := c.Params("classroomId")

	invitationId, err := uuid.Parse(invitationIdParameter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	classroomId, err := uuid.Parse(classroomIdParameter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryClassroomInvitation := query.ClassroomInvitation

	invitation, err := queryClassroomInvitation.
		WithContext(c.Context()).
		Preload(queryClassroomInvitation.Classroom).
		Where(queryClassroomInvitation.ID.Eq(invitationId)).
		Where(queryClassroomInvitation.ClassroomID.Eq(classroomId)).
		Where(queryClassroomInvitation.Status.Eq(uint8(database.InvitationPending))).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if time.Now().After(invitation.ExpiryDate) {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	repo := context.GetGitlabRepository(c)
	currentUser, err := repo.GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if invitation.Email != currentUser.Email {
		return fiber.NewError(fiber.StatusForbidden, "You are not the chosen one")
	}

	// reauthenticate the repo with the group access token
	err = repo.GroupAccessLogin(invitation.Classroom.GroupAccessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = repo.AddUserToGroup(invitation.Classroom.GroupID, currentUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		err = tx.UserClassrooms.Create(&database.UserClassrooms{
			UserID:    currentUser.ID,
			Classroom: invitation.Classroom,
			Role:      database.Student,
		})
		if err != nil {
			newErr := repo.RemoveUserFromGroup(invitation.Classroom.GroupID, currentUser.ID)
			if newErr != nil {
				return fiber.NewError(fiber.StatusInternalServerError, newErr.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		invitation.Status = database.InvitationAccepted
		err = tx.ClassroomInvitation.Save(invitation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
