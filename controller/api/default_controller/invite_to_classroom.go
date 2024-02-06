package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"log"
	"net/mail"
	"time"
)

type InviteToClassroomRequest struct {
	MemberEmails []string `json:"memberEmails"`
}

func (r InviteToClassroomRequest) isValid() bool {
	return len(r.MemberEmails) != 0
}

func (ctrl *DefaultController) InviteToClassroom(c *fiber.Ctx) error {
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

	repo := context.GetGitlabRepository(c)
	if err := ctrl.RotateAccessToken(c, classroom); err != nil {
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
	invites, err := queryClassroomInvite.WithContext(c.Context()).Where(queryClassroomInvite.ClassroomID.Eq(classroomId)).Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	invitableEmails := filterInvitableEmails(validatedEmailAddresses, invites)

	invitations := make([]*database.ClassroomInvitation, len(invitableEmails))

	// Create invitations
	err = query.Q.Transaction(func(tx *query.Query) error {
		for i, email := range invitableEmails {
			_, err := tx.ClassroomInvitation.
				WithContext(c.Context()).
				Where(tx.ClassroomInvitation.Email.Eq(email.Address)).
				Where(tx.ClassroomInvitation.ClassroomID.Eq(classroomId)).
				Delete()
			if err != nil {
				return err
			}

			newInvitation := &database.ClassroomInvitation{
				Status:      database.ClassroomInvitationPending,
				ClassroomID: classroomId,
				Email:       email.Address,
				Enabled:     true,
				ExpiryDate:  time.Now().AddDate(0, 0, 14),
			}
			if err := tx.ClassroomInvitation.WithContext(c.Context()).Create(newInvitation); err != nil {
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
		log.Println("Sending invitation to", email.Address)

		invitePath := fmt.Sprintf("/classrooms/%s/invitations/%s", classroom.ID.String(), invitations[i].ID.String())
		err = ctrl.mailRepo.SendClassroomInvitation(email.Address,
			fmt.Sprintf(`New Invitation for Classroom "%s"`,
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

func filterInvitableEmails(emails []*mail.Address, invitations []*database.ClassroomInvitation) []*mail.Address {
	invitableEmails := make([]*mail.Address, 0)

	for _, email := range emails {
		found := false
		for _, invite := range invitations {
			if invite.Email == email.Address {
				if invite.Status != database.ClassroomInvitationAccepted {
					invitableEmails = append(invitableEmails, email)
				}
				found = true
				break
			}
		}
		if !found {
			invitableEmails = append(invitableEmails, email)
		}
	}

	return invitableEmails
}
