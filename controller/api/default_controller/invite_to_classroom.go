package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"net/http"
	"net/mail"
	"time"
)

type InviteToClassroomRequest struct {
	MemberEmails []string `json:"memberEmails"`
}

func (handler *DefaultController) InviteToClassroom(c *fiber.Ctx) error {
	user, err := session.Get(c).GetUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryClassroom := query.Classroom
	classroom, err := queryClassroom.
		WithContext(c.Context()).
		Where(queryClassroom.ID.Eq(classroomId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Check if owner or moderator of specific classroom
	if classroom.OwnerID != user.ID {
		queryUserClassroom := query.UserClassrooms
		_, err := queryUserClassroom.
			WithContext(c.Context()).
			Where(queryUserClassroom.ClassroomID.Eq(classroomId)).
			Where(queryUserClassroom.UserID.Eq(user.ID)).
			Where(queryUserClassroom.Role.Eq(uint8(database.Moderator))).
			First()
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
	}

	requestBody := new(InviteToClassroomRequest)
	if err = c.BodyParser(requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate email addresses
	validatedEmailAddresses := make([]*mail.Address, len(requestBody.MemberEmails))
	for i, email := range requestBody.MemberEmails {
		validatedEmailAddresses[i], err = mail.ParseAddress(email)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	// Create invitations
	err = query.Q.Transaction(func(tx *query.Query) error {
		for _, email := range validatedEmailAddresses {
			newInvitation := database.ClassroomInvitation{
				Status:      database.InvitationPending,
				ClassroomID: classroomId,
				Email:       email.Address,
				Enabled:     true,
				ExpiryDate:  time.Now().Add(time.Hour * 24 * 14), // Two Weeks, TODO: Add to configuration
			}
			err := tx.ClassroomInvitation.Create(&newInvitation)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send invitations
	for _, email := range validatedEmailAddresses {
		err = handler.mailRepo.Send(email.Address, fmt.Sprintf("Test: New Invitation from Classroom %s", classroom.Name), mailRepo.MailData{}) // TODO: Add meaningful content
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(http.StatusAccepted)
}
