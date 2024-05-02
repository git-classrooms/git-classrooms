package default_controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type createClassroomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateTeams *bool  `json:"createTeams"`
	MaxTeams    *int   `json:"maxTeams"`
	MaxTeamSize int    `json:"maxTeamSize"`
} //@Name CreateClassroomRequest

func (r createClassroomRequest) isValid() bool {
	return r.Name != "" &&
		r.Description != "" &&
		r.CreateTeams != nil &&
		r.MaxTeamSize > 0 &&
		r.MaxTeams != nil && *r.MaxTeams >= 0
}

// @Summary		Create a new classroom
// @Description	Create a new classroom
// @Id				CreateClassroom
// @Tags			classroom
// @Accept			json
// @Param			classroom		body	default_controller.createClassroomRequest	true	"Classroom Info"
// @Param			X-Csrf-Token	header	string										true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v1/classroom/owned/{classroomId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/classrooms/owned [post]
func (ctrl *DefaultController) CreateClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()

	userID := ctx.GetUserID()

	requestBody := &createClassroomRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	group, err := repo.CreateGroup(
		requestBody.Name,
		model.Private,
		requestBody.Description,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	expiresAt := time.Now().AddDate(0, 0, 364)

	accessToken, err := repo.CreateGroupAccessToken(group.ID, "Gitlab Classrooms", model.OwnerPermissions, expiresAt, "api")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	classroomQuery := query.Classroom
	classRoom := &database.Classroom{
		Name:               requestBody.Name,
		Description:        requestBody.Description,
		OwnerID:            userID,
		CreateTeams:        *requestBody.CreateTeams,
		MaxTeamSize:        requestBody.MaxTeamSize,
		MaxTeams:           *requestBody.MaxTeams,
		GroupID:            group.ID,
		GroupAccessTokenID: accessToken.ID,
		GroupAccessToken:   accessToken.Token,
	}

	err = classroomQuery.WithContext(c.Context()).Create(classRoom)
	if err != nil {
		newErr := repo.DeleteGroup(group.ID)
		if newErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, newErr.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/owned/%s", classRoom.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
