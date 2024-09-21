package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
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
// @Id				JoinClassroomV2
// @Tags			classroom
// @Accept			json
// @Param			classroomId		path	string						true	"Classroom ID"	Format(uuid)
// @Param			invitation		body	api.joinClassroomRequest	true	"Invitation"
// @Param			X-Csrf-Token	header	string						true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v1/classroom/{classroomId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/join [post]
func (ctrl *DefaultController) JoinClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()

	var params Params
	if err = c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if params.ClassroomID == nil {
		return fiber.ErrBadRequest
	}

	var requestBody joinClassroomRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	queryClassroomInvitation := query.ClassroomInvitation
	invitation, err := queryClassroomInvitation.
		WithContext(c.Context()).
		Preload(queryClassroomInvitation.Classroom).
		Where(queryClassroomInvitation.ClassroomID.Eq(*params.ClassroomID)).
		Where(queryClassroomInvitation.ID.Eq(requestBody.InvitationID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch invitation.Status {
	case database.ClassroomInvitationRevoked:
		return fiber.NewError(fiber.StatusForbidden, "This invitation has been revoked.")
	case database.ClassroomInvitationPending:
		break
	default:
		return fiber.NewError(fiber.StatusForbidden, "This invitation has already been processed.")
	}

	if time.Now().After(invitation.ExpiryDate) {
		return fiber.NewError(fiber.StatusForbidden, "The link to this classroom expired. Please ask the owner for a new invitation link.")
	}

	userID := ctx.GetUserID()
	queryUser := query.User
	user, err := queryUser.WithContext(c.Context()).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryUserClassrooms := query.UserClassrooms
	_, err = queryUserClassrooms.WithContext(c.Context()).
		Where(queryUserClassrooms.UserID.Eq(userID)).
		Where(queryUserClassrooms.ClassroomID.Eq(invitation.ClassroomID)).
		First()
	if err == nil {
		if _, err := queryClassroomInvitation.WithContext(c.Context()).Delete(invitation); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fiber.NewError(fiber.StatusForbidden, "You are already a member of this classroom.")
	}

	if requestBody.Action == reject {
		invitation.Status = database.ClassroomInvitationRejected
		invitation.Email = user.GitlabEmail
		err = queryClassroomInvitation.WithContext(c.Context()).Save(invitation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(fiber.StatusAccepted)
	}

	// reauthenticate the repo with the group access token
	if err = repo.GroupAccessLogin(invitation.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = query.Q.Transaction(func(tx *query.Query) (err error) {
		member := &database.UserClassrooms{
			UserID:    userID,
			Classroom: invitation.Classroom,
			Role:      database.Student,
		}
		if err = tx.UserClassrooms.WithContext(c.Context()).Create(member); err != nil {
			return err
		}

		invitation.Status = database.ClassroomInvitationAccepted
		invitation.Email = user.GitlabEmail
		if err = tx.ClassroomInvitation.WithContext(c.Context()).Save(invitation); err != nil {
			return err
		}

		groupRole := gitlabModel.GuestPermissions
		if invitation.Classroom.StudentsViewAllProjects {
			groupRole = gitlabModel.ReporterPermissions
		}

		if err = repo.AddUserToGroup(invitation.Classroom.GroupID, userID, groupRole); err != nil {
			return err
		}
		defer func() {
			if recover() != nil || err != nil {
				repo.RemoveUserFromGroup(invitation.Classroom.GroupID, userID)
			}
		}()

		if invitation.Classroom.MaxTeamSize == 1 {
			var subgroup *gitlabModel.Group
			subgroup, err = repo.CreateSubGroup(
				user.Name,
				user.GitlabUsername,
				invitation.Classroom.GroupID,
				gitlabModel.Private,
				fmt.Sprintf("Team %s of classroom %s", user.Name, invitation.Classroom.Name),
			)
			if err != nil {
				return err
			}
			defer func() {
				if recover() != nil || err != nil {
					repo.DeleteGroup(subgroup.ID)
				}
			}()

			team := &database.Team{
				ClassroomID: invitation.Classroom.ID,
				Name:        user.Name,
				GroupID:     subgroup.ID,
				Member:      []*database.UserClassrooms{member},
			}
			if err = tx.Team.WithContext(c.Context()).Create(team); err != nil {
				return err
			}

			repo.ChangeGroupDescription(subgroup.ID, utils.CreateTeamGitlabDescription(&invitation.Classroom, team, ctrl.config.PublicURL))

			if err = repo.AddUserToGroup(subgroup.ID, userID, gitlabModel.ReporterPermissions); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/%s", invitation.ClassroomID.String()))
	return c.SendStatus(fiber.StatusAccepted)
}
