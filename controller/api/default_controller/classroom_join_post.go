package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/gorm"
)

type action string //@Name action

const (
	accept action = "accept"
	reject action = "reject"
)

type joinClassroomRequest struct {
	InvitationID  uuid.UUID `json:"invitationId"`
	Action        action    `json:"action"`
	ClassroomCode bool      `json:"classroomCode"`
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
	log := ctx.GetLoggerForHandler("JoinClassroom")
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
		return fiber.NewError(fiber.StatusBadRequest, "Body is not valid")
	}

	log.Debug("retrieving user from db")

	userID := ctx.GetUserID()
	queryUser := query.User
	user, err := queryUser.WithContext(c.Context()).
		Where(queryUser.ID.Eq(userID)).
		First()
	if err != nil {
		log.Error("error getting user from db", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryClassroomInvitation := query.ClassroomInvitation
	var invitation *database.ClassroomInvitation
	if requestBody.ClassroomCode {
		log.Info("received an classroom code invitation request")

		if requestBody.Action == reject {
			log.Info("join classroom was rejected")
			return c.SendStatus(fiber.StatusAccepted)
		}

		queryClassroom := query.Classroom
		classroom, err := queryClassroom.WithContext(c.Context()).
			Where(queryClassroom.ID.Eq(*params.ClassroomID)).
			First()
		if err != nil {
			log.Error("error getting classroom from db", "error", err)
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		log = log.With("classroom", classroom)

		if classroom.InviteCode != requestBody.InvitationID {
			log.Error("classroom code is not valid", "invitationCode", requestBody.InvitationID)
			return fiber.NewError(fiber.StatusForbidden, "This invitation has been revoked or does not exist.")
		}

		queryUserClassrooms := query.UserClassrooms
		_, err = queryUserClassrooms.WithContext(c.Context()).
			Where(queryUserClassrooms.UserID.Eq(userID)).
			Where(queryUserClassrooms.ClassroomID.Eq(*params.ClassroomID)).
			First()
		if err == nil {
			return fiber.NewError(fiber.StatusForbidden, "You are already a member of this classroom.")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("error getting userclassroom from db", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		invitation = &database.ClassroomInvitation{
			Status:      database.ClassroomInvitationAccepted,
			ClassroomID: *params.ClassroomID,
			Classroom:   *classroom,
			Email:       user.GitlabEmail,
			ExpiryDate:  time.Now().AddDate(0, 0, 14),
		}

		log.Debug("saving invitation in db", "invitation", invitation)

		if err := queryClassroomInvitation.
			WithContext(c.Context()).
			Create(invitation); err != nil {
			log.Error("error creating invitation in db", "error", err)
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				log.Warn("error occured, cleaning up: invitation")
				if queryClassroomInvitation.
					WithContext(c.Context()).
					Delete(invitation); err != nil {
					log.Error("error while deleting invitation", "error", err)
				}
			}
		}()
	} else {
		log.Info("received a normal invitation request")

		invitation, err = queryClassroomInvitation.
			WithContext(c.Context()).
			Preload(queryClassroomInvitation.Classroom).
			Where(queryClassroomInvitation.ClassroomID.Eq(*params.ClassroomID)).
			Where(queryClassroomInvitation.ID.Eq(requestBody.InvitationID)).
			First()
		if err != nil {
			log.Error("error getting invitation in db", "error", err)
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		log = log.With("classroom", invitation.Classroom)

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

		queryUserClassrooms := query.UserClassrooms
		_, err = queryUserClassrooms.WithContext(c.Context()).
			Where(queryUserClassrooms.UserID.Eq(userID)).
			Where(queryUserClassrooms.ClassroomID.Eq(invitation.ClassroomID)).
			First()
		if err == nil {
			log.Info("user is already member of this classroom")
			if _, err := queryClassroomInvitation.WithContext(c.Context()).Delete(invitation); err != nil {
				log.Error("error removing invitation", "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			return fiber.NewError(fiber.StatusForbidden, "You are already a member of this classroom.")
		}

		if requestBody.Action == reject {
			log.Info("invitation was rejected")
			invitation.Status = database.ClassroomInvitationRejected
			invitation.Email = user.GitlabEmail
			err = queryClassroomInvitation.WithContext(c.Context()).Save(invitation)
			if err != nil {
				log.Error("error setting invitation to rejected", "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			return c.SendStatus(fiber.StatusAccepted)
		}
	}

	log.Debug("authenticate the repo with group access token")

	// reauthenticate the repo with the group access token
	if err = repo.GroupAccessLogin(invitation.Classroom.GroupAccessToken, log); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = query.Q.Transaction(func(tx *query.Query) (err error) {
		member := &database.UserClassrooms{
			UserID:    userID,
			Classroom: invitation.Classroom,
			Role:      database.Student,
		}
		if err = tx.UserClassrooms.WithContext(c.Context()).Create(member); err != nil {
			log.Error("error creating classroom member", "error", err)
			return err
		}

		log.Info("member was added to the classroom")

		invitation.Status = database.ClassroomInvitationAccepted
		invitation.Email = user.GitlabEmail
		if err = tx.ClassroomInvitation.WithContext(c.Context()).Save(invitation); err != nil {
			log.Error("error setting invitation to accepted", "error", err)
			return err
		}

		log.Debug("invitation was set to accepted")

		groupRole := gitlabModel.GuestPermissions
		if invitation.Classroom.StudentsViewAllProjects {
			groupRole = gitlabModel.ReporterPermissions
		}

		if err = repo.AddUserToGroup(invitation.Classroom.GroupID, userID, groupRole); err != nil {
			log.Error("error adding user to group", "groupID", invitation.Classroom.GroupID, "error", err)
			return err
		}
		defer func() {
			if recover() != nil || err != nil {
				log.Warn("error occured, cleaning up: user added to group", "groupID", invitation.Classroom.GroupID)
				if err := repo.RemoveUserFromGroup(invitation.Classroom.GroupID, userID); err != nil {
					log.Error("error while removing user from group", "groupID", invitation.Classroom.GroupID, "error", err)
				}
			}
		}()

		log.Info("user added to gitlab group", "groupID", invitation.Classroom.GroupID)

		if invitation.Classroom.MaxTeamSize == 1 {
			log.Info("teams are disabled creating team for user")

			var subgroup *gitlabModel.Group
			subgroup, err = repo.CreateSubGroup(
				user.Name,
				user.GitlabUsername,
				invitation.Classroom.GroupID,
				gitlabModel.Private,
				fmt.Sprintf("Team %s of classroom %s", user.Name, invitation.Classroom.Name),
			)
			if err != nil {
				log.Error("error creating subgroup for team", "groupID", invitation.Classroom.GroupID, "error", err)
				return err
			}
			defer func() {
				if recover() != nil || err != nil {
					log.Warn("error occurred, cleaning up: subgroup")
					if err := repo.DeleteGroup(subgroup.ID); err != nil {
						log.Error("error deleting group", "error", err)
					}
				}
			}()

			log.Info("subgroup created for new team")

			team := &database.Team{
				ClassroomID: invitation.Classroom.ID,
				Name:        user.Name,
				GroupID:     subgroup.ID,
				Member:      []*database.UserClassrooms{member},
			}
			if err = tx.Team.WithContext(c.Context()).Create(team); err != nil {
				log.Error("error while creating team", "error", err)
				return err
			}
			ctrl.createAssignmentProjectsIfNeeded(c.Context(), query.Q, &invitation.Classroom, team.ID)

			log.Debug("changing group description")
			if _, err := repo.ChangeGroupDescription(subgroup.ID, utils.CreateTeamGitlabDescription(&invitation.Classroom, team, ctrl.config.PublicURL)); err != nil {
				log.Error("error changing subgroup description", "subGroupID", subgroup.ID, "error", err)
			}

			log.Debug("adding user to team subgroup", "subGroupID", subgroup.ID)

			if err = repo.AddUserToGroup(subgroup.ID, userID, gitlabModel.ReporterPermissions); err != nil {
				log.Error("error adding user to group", "subGroupID", subgroup.ID, "error", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Debug("error in transaction", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/%s", invitation.ClassroomID.String()))
	return c.SendStatus(fiber.StatusAccepted)
}
