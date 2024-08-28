package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type createClassroomRequest struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	CreateTeams             *bool  `json:"createTeams"`
	MaxTeams                *int   `json:"maxTeams"`
	MaxTeamSize             int    `json:"maxTeamSize"`
	StudentsViewAllProjects *bool  `json:"studentsViewAllProjects"`
} //@Name CreateClassroomRequest

func (r createClassroomRequest) isValid() bool {
	return r.Name != "" &&
		r.Description != "" &&
		r.CreateTeams != nil &&
		r.MaxTeamSize > 0 &&
		r.MaxTeams != nil && *r.MaxTeams >= 0 &&
		r.StudentsViewAllProjects != nil
}

// @Summary		Create a new classroom
// @Description	Create a new classroom
// @Id				CreateClassroomV2
// @Tags			classroom
// @Accept			json
// @Param			classroom		body	api.createClassroomRequest	true	"Classroom Info"
// @Param			X-Csrf-Token	header	string						true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v2/classroom/{classroomId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms [post]
func (ctrl *DefaultController) CreateClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()

	userID := ctx.GetUserID()

	var requestBody createClassroomRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	queryUser := query.User
	user, err := queryUser.WithContext(c.Context()).Where(queryUser.ID.Eq(userID)).First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	group, err := repo.CreateGroup(
		requestBody.Name,
		model.Private,
		requestBody.Description,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		if recover() != nil || err != nil {
			if err := repo.DeleteGroup(group.ID); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	expiresAt := time.Now().AddDate(0, 0, 364)

	accessToken, err := repo.CreateGroupAccessToken(group.ID, "Gitlab Classrooms", model.OwnerPermissions, expiresAt, "api")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to delete the accessToken because it will be deleted when the group is deleted

	var classroom *database.Classroom
	err = query.Q.Transaction(func(tx *query.Query) error {
		classroomQuery := tx.Classroom
		classroom = &database.Classroom{
			Name:                    requestBody.Name,
			Description:             requestBody.Description,
			OwnerID:                 userID,
			CreateTeams:             *requestBody.CreateTeams,
			MaxTeamSize:             requestBody.MaxTeamSize,
			MaxTeams:                *requestBody.MaxTeams,
			GroupID:                 group.ID,
			GroupAccessTokenID:      accessToken.ID,
			GroupAccessToken:        accessToken.Token,
			StudentsViewAllProjects: *requestBody.StudentsViewAllProjects,
			Member:                  []*database.UserClassrooms{{UserID: userID, Role: database.Owner}},
		}

		if err = classroomQuery.WithContext(c.Context()).Create(classroom); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		invitation := &database.ClassroomInvitation{
			Status:      database.ClassroomInvitationAccepted,
			ClassroomID: classroom.ID,
			Email:       user.GitlabEmail,
			ExpiryDate:  time.Now().AddDate(0, 0, 14),
		}
		if err = tx.ClassroomInvitation.WithContext(c.Context()).Create(invitation); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v2/classrooms/%s", classroom.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
