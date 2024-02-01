package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"net/mail"
	"time"
)

type InviteToClassroomRequest struct {
	MemberEmails []string `json:"memberEmails"`
}

func (r InviteToClassroomRequest) isValid() bool {
	return len(r.MemberEmails) != 0
}

func (handler *DefaultController) InviteToClassroom(c *fiber.Ctx) error {
	user, err := session.Get(c).GetUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomId, err := uuid.Parse(c.Params("classroomId"))
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

	expiresAt := time.Now().AddDate(0, 0, 14) // Two Weeks, TODO: Add to configuration

	repo := context.GetGitlabRepository(c)
	// TODO: check if an accessToken already exist and update the expiration date of that one instead of creating a new one
	accessToken, err := repo.CreateGroupAccessToken(classroom.GroupID, "Gitlab-Classroom-Access-Token", model.OwnerPermissions, expiresAt, "api")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroom.GroupAccessToken = accessToken.Token
	err = queryClassroom.WithContext(c.Context()).Save(classroom)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	requestBody := &InviteToClassroomRequest{}
	if err = c.BodyParser(requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Validate email addresses
	validatedEmailAddresses := make([]*mail.Address, len(requestBody.MemberEmails))
	for i, email := range requestBody.MemberEmails {
		validatedEmailAddresses[i], err = mail.ParseAddress(email)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	queryClassroomInvite := query.ClassroomInvitation
	invites, err := queryClassroomInvite.Where(queryClassroomInvite.ClassroomID.Eq(classroomId)).Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	invitableEmails := make([]*mail.Address, 0)

	for _, email := range validatedEmailAddresses {
		for _, invite := range invites {
			if invite.Email == email.Address {
				if invite.Status != database.InvitationAccepted {
					invitableEmails = append(invitableEmails, email)
					break
				}
			} else {
				invitableEmails = append(invitableEmails, email)
				break
			}
		}
	}

	invitations := make([]*database.ClassroomInvitation, len(invitableEmails))

	// Create invitations
	err = query.Q.Transaction(func(tx *query.Query) error {
		for i, email := range invitableEmails {
			_, err := tx.ClassroomInvitation.
				Where(tx.ClassroomInvitation.Email.Eq(email.Address)).
				Where(tx.ClassroomInvitation.ClassroomID.Eq(classroomId)).
				Delete()
			if err != nil {
				return err
			}

			newInvitation := &database.ClassroomInvitation{
				Status:      database.InvitationPending,
				ClassroomID: classroomId,
				Email:       email.Address,
				Enabled:     true,
				ExpiryDate:  expiresAt,
			}
			if err := tx.ClassroomInvitation.Create(newInvitation); err != nil {
				return err
			}
			invitations[i] = newInvitation
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send invitations
	owner, err := repo.GetUserById(classroom.OwnerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for i, email := range invitableEmails {
		invitePath := fmt.Sprintf("/api/classrooms/%s/invitations/%s", classroom.ID.String(), invitations[i].ID.String())
		err = handler.mailRepo.SendClassroomInvitation(email.Address,
			fmt.Sprintf(`Test: New Invitation for Classroom "%s"`,
				classroom.Name),
			mailRepo.ClassroomInvitationData{
				ClassroomName:      classroom.Name,
				ClassroomOwnerName: owner.Name,
				RecipientEmail:     invitations[i].Email,
				InvitationPath:     invitePath,
				ExpireDate:         invitations[i].ExpiryDate,
			})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}
