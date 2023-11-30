package handler

import (
	"backend/api/repository"
	"backend/model"

	"github.com/gofiber/fiber/v2"
)

type FiberHandler struct {
	repo repository.Repository
}

func NewFiberHandler(repository repository.Repository) *FiberHandler {
	return &FiberHandler{repository}
}

type ClassroomRequest struct {
	Name         string   `json:"name"`
	MemberEmails []string `json:"memberEmails"`
	Description  string   `json:"description"`
}

func (handler *FiberHandler) CreateClassroom(c *fiber.Ctx) error {
	var err error
	requestBody := new(ClassroomRequest)

	// TODO clarify how it should be known which user wants to create the classroom, so that this user could be user for the gitlab api
	// handler.repo.Login()

	err = c.BodyParser(requestBody)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	group, err := handler.repo.CreateGroup(
		requestBody.Name,
		model.Private,
		requestBody.Description,
		requestBody.MemberEmails,
	)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return err
	}

	for _, memberEmail := range requestBody.MemberEmails {
		err = handler.repo.CreateGroupInvite(group.ID, memberEmail)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			return err
		}
	}

	return nil
}
