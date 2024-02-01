package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"net/http"
	"time"
)

func (handler *DefaultController) JoinClassroom(c *fiber.Ctx) error {
	invitationIdParameter := c.Params("invitationId")

	invitationId, err := uuid.Parse(invitationIdParameter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	queryClassroomInvitation := query.ClassroomInvitation
	invitation, err := queryClassroomInvitation.
		WithContext(c.Context()).
		Where(queryClassroomInvitation.ClassroomID.Eq(invitationId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if time.Now().After(invitation.ExpiryDate) {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	currentUser, err := context.GetGitlabRepository(c).GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomIdParameter := c.Params("classroomId")
	classroomId, err := uuid.Parse(classroomIdParameter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if classroomId != invitation.ClassroomID {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if invitation.Email == currentUser.Email {
		repo := context.GetGitlabRepository(c)
		err := repo.AddUserToGroup(invitation.Classroom.GroupID, currentUser.ID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		queryClassroom := query.UserClassrooms
		err = queryClassroom.Create(&database.UserClassrooms{
			UserID:    currentUser.ID,
			Classroom: invitation.Classroom,
			Role:      database.Student,
		})
		if err != nil {
			err := repo.RemoveUserFromGroup(invitation.Classroom.GroupID, currentUser.ID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		_, err = queryClassroomInvitation.WithContext(c.Context()).
			Where(queryClassroomInvitation.ID.Eq(invitation.ID)).
			Update(queryClassroomInvitation.Status, database.InvitationAccepted)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(http.StatusAccepted)
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
}
