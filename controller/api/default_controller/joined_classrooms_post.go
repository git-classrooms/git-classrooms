package default_controller

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type action string //@Name action

const (
	accept action = "accept"
	reject action = "reject"
)

type joinClassroomRequest struct {
	InvitationID uuid.UUID `json:"invitationId"`
	Action       action    `json:"action"`
} //@Name JoinClassroomRequest

func (r *joinClassroomRequest) isValid() bool {
	return r.InvitationID != uuid.Nil &&
		(r.Action == accept || r.Action == reject)
}

// @Summary		JoinClassroom
// @Description	JoinClassroom
// @Id				JoinClassroom
// @Tags			classroom
// @Accept			json
// @Param			invitation		body	default_controller.joinClassroomRequest	true	"Invitation"
// @Param			X-Csrf-Token	header	string									true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v1/classroom/joined/{classroomId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/classrooms/joined [post]
func (*DefaultController) JoinClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()

	requestBody := &joinClassroomRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	log.Println(requestBody)

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	queryClassroomInvitation := query.ClassroomInvitation

	invitation, err := queryClassroomInvitation.
		WithContext(c.Context()).
		Preload(queryClassroomInvitation.Classroom).
		Where(queryClassroomInvitation.ID.Eq(requestBody.InvitationID)).
		Where(queryClassroomInvitation.Status.Eq(uint8(database.ClassroomInvitationPending))).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if time.Now().After(invitation.ExpiryDate) {
		return fiber.NewError(fiber.StatusForbidden, "The link to this classroom expired. Please ask the owner for a new invitation link.")
	}

	currentUser, err := repo.GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// reauthenticate the repo with the group access token
	err = repo.GroupAccessLogin(invitation.Classroom.GroupAccessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if requestBody.Action == reject {
		log.Println("rejected invitation")
		invitation.Status = database.ClassroomInvitationRejected
		invitation.Email = currentUser.Email
		err = queryClassroomInvitation.WithContext(c.Context()).Save(invitation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(fiber.StatusAccepted)
	}

	log.Println("we should not be here")

	err = query.Q.Transaction(func(tx *query.Query) (err error) {
		member := &database.UserClassrooms{
			UserID:    currentUser.ID,
			Classroom: invitation.Classroom,
			Role:      database.Student,
		}
		err = tx.UserClassrooms.
			WithContext(c.Context()).
			Create(member)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		queryUserClassrooms := tx.UserClassrooms
		member, err = queryUserClassrooms.
			WithContext(c.Context()).
			Preload(queryUserClassrooms.User).
			Where(queryUserClassrooms.ClassroomID.Eq(member.ClassroomID)).
			Where(queryUserClassrooms.UserID.Eq(member.UserID)).
			First()

		invitation.Status = database.ClassroomInvitationAccepted
		invitation.Email = currentUser.Email
		err = tx.ClassroomInvitation.WithContext(c.Context()).Save(invitation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		groupRole := gitlabModel.GuestPermissions
		if invitation.Classroom.StudentsViewAllProjects {
			groupRole = gitlabModel.ReporterPermissions
		}

		err = repo.AddUserToGroup(invitation.Classroom.GroupID, currentUser.ID, groupRole)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				repo.RemoveUserFromGroup(invitation.Classroom.GroupID, currentUser.ID)
			}
		}()

		if invitation.Classroom.MaxTeamSize == 1 {
			var subgroup *gitlabModel.Group
			subgroup, err = repo.CreateSubGroup(
				member.User.Name,
				invitation.Classroom.GroupID,
				gitlabModel.Private,
				fmt.Sprintf("Team %s of classroom %s", member.User.Name, invitation.Classroom.Name),
			)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			defer func() {
				if recover() != nil || err != nil {
					repo.DeleteGroup(subgroup.ID)
				}
			}()

			team := &database.Team{
				ClassroomID: invitation.Classroom.ID,
				Name:        currentUser.Username,
				GroupID:     subgroup.ID,
				Member:      []*database.UserClassrooms{member},
			}
			err = tx.Team.WithContext(c.Context()).Create(team)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			err = repo.AddUserToGroup(subgroup.ID, currentUser.ID, gitlabModel.ReporterPermissions)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/joined/%s", invitation.ClassroomID.String()))
	return c.SendStatus(fiber.StatusAccepted)
}
