package default_controller

import (
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type inviteToClassroomRequest struct {
	MemberEmails []string `json:"memberEmails"`
}

func (r inviteToClassroomRequest) isValid() bool {
	return len(r.MemberEmails) != 0
}

func (ctrl *DefaultController) InviteToClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	repo := ctx.GetGitlabRepository()
	if err := ctrl.RotateAccessToken(c, classroom); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	requestBody := &inviteToClassroomRequest{}
	if err := c.BodyParser(requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Validate email addresses
	validatedEmailAddresses := make([]*mail.Address, len(requestBody.MemberEmails))
	for i, email := range requestBody.MemberEmails {
		var err error
		validatedEmailAddresses[i], err = mail.ParseAddress(email)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	queryClassroomInvite := query.ClassroomInvitation
	invites, err := queryClassroomInvite.WithContext(c.Context()).Where(queryClassroomInvite.ClassroomID.Eq(classroom.ID)).Find()
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
				Where(tx.ClassroomInvitation.ClassroomID.Eq(classroom.ID)).
				Delete()
			if err != nil {
				return err
			}

			newInvitation := &database.ClassroomInvitation{
				Status:      database.ClassroomInvitationPending,
				ClassroomID: classroom.ID,
				Email:       email.Address,
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

		invitePath := fmt.Sprintf("/classrooms/joined/%s/invitations/%s", classroom.ID.String(), invitations[i].ID.String())
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
