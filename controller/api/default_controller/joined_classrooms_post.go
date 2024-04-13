package default_controller

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type joinClassroomRequest struct {
	InvitationID uuid.UUID `json:"invitationId"`
}

func (*DefaultController) JoinClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()

	requestBody := &joinClassroomRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	log.Println(requestBody)

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
		return fiber.NewError(fiber.StatusForbidden, err.Error())
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

	err = query.Q.Transaction(func(tx *query.Query) (err error) {
		member := &database.UserClassrooms{
			UserID:    currentUser.ID,
			Classroom: invitation.Classroom,
			Role:      database.Student,
		}
		err = tx.UserClassrooms.
			WithContext(c.Context()).
			Preload(query.UserClassrooms.User).
			Create(member)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		invitation.Status = database.ClassroomInvitationAccepted
		invitation.Email = currentUser.Email
		err = tx.ClassroomInvitation.WithContext(c.Context()).Save(invitation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = repo.AddUserToGroup(invitation.Classroom.GroupID, currentUser.ID, gitlabModel.ReporterPermissions)
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
				model.Private,
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

			err = repo.AddUserToGroup(subgroup.ID, currentUser.ID, model.DeveloperPermissions)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
