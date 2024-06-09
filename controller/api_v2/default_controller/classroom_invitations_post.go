package api

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
} //@Name InviteToClassroomRequest

func (r inviteToClassroomRequest) isValid() bool {
	return len(r.MemberEmails) != 0
}

// @Summary		InviteToClassroom
// @Description	InviteToClassroom
// @Id				InviteToClassroomV2
// @Tags			classroom
// @Accept			json
// @Param			classroomId		path	string							true	"Classroom ID"	Format(uuid)
// @Param			memberEmails	body	api.inviteToClassroomRequest	true	"Member Emails"
// @Param			X-Csrf-Token	header	string							true	"Csrf-Token"
// @Success		201
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/invitations [post]
func (ctrl *DefaultController) InviteToClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	if err := ctrl.RotateAccessToken(c, &classroom.Classroom); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var requestBody inviteToClassroomRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
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
	invites, err := queryClassroomInvite.WithContext(c.Context()).
		Where(queryClassroomInvite.ClassroomID.Eq(classroom.ClassroomID)).
		Find()
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
				Where(tx.ClassroomInvitation.ClassroomID.Eq(classroom.ClassroomID)).
				Delete()
			if err != nil {
				return err
			}

			newInvitation := &database.ClassroomInvitation{
				Status:      database.ClassroomInvitationPending,
				ClassroomID: classroom.ClassroomID,
				Email:       email.Address,
				ExpiryDate:  time.Now().AddDate(0, 0, 14),
			}
			if err = tx.ClassroomInvitation.WithContext(c.Context()).Create(newInvitation); err != nil {
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
	for i, email := range invitableEmails {
		log.Println("Sending invitation to", email.Address)

		invitePath := fmt.Sprintf("/classrooms/%s/invitations/%s", classroom.ClassroomID.String(), invitations[i].ID.String())
		err = ctrl.mailRepo.SendClassroomInvitation(email.Address,
			fmt.Sprintf(`New Invitation for Classroom "%s"`,
				classroom.Classroom.Name),
			mailRepo.ClassroomInvitationData{
				ClassroomName:      classroom.Classroom.Name,
				ClassroomOwnerName: classroom.Classroom.Owner.Name,
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
