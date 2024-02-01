package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"net/http"
)

type CreateClassroomRequest struct {
	Name         string   `json:"name"`
	MemberEmails []string `json:"memberEmails"`
	Description  string   `json:"description"`
}

func (handler *DefaultController) CreateClassroom(c *fiber.Ctx) error { //TODO: Rework, group-access-token and persisting data is needed
	repo := context.GetGitlabRepository(c)

	var err error
	requestBody := &CreateClassroomRequest{}

	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	group, err := repo.CreateGroup(
		requestBody.Name,
		model.Private,
		requestBody.Description,
		requestBody.MemberEmails, //TODO: User shouldn't be added in advance. They should be invited and have the chance not to accept
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, memberEmail := range requestBody.MemberEmails {
		err = repo.CreateGroupInvite(group.ID, memberEmail)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.SendStatus(http.StatusCreated)
}
