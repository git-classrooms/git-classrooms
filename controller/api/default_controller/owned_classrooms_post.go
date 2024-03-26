package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"time"
)

type CreateClassroomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r CreateClassroomRequest) isValid() bool {
	return r.Name != "" && r.Description != ""
}

func (ctrl *DefaultController) CreateClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()

	userID := ctx.GetUserID()

	requestBody := &CreateClassroomRequest{}
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
		OwnerID:            userID,
		Description:        requestBody.Description,
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
