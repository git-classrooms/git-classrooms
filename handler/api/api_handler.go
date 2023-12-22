package apiHandler

import (
	"backend/api/repository"
	"backend/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type FiberApiHandler struct {
}

func NewFiberApiHandler() *FiberApiHandler {
	return &FiberApiHandler{}
}

type ClassroomRequest struct {
	Name         string   `json:"name"`
	MemberEmails []string `json:"memberEmails"`
	Description  string   `json:"description"`
}

func (handler *FiberApiHandler) CreateClassroom(c *fiber.Ctx) error {
	repo := handler.getRepo(c)

	var err error
	requestBody := new(ClassroomRequest)

	err = c.BodyParser(requestBody)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	group, err := repo.CreateGroup(
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
		err = repo.CreateGroupInvite(group.ID, memberEmail)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			return err
		}
	}

	c.Status(http.StatusCreated)
	return nil
}

func (handler *FiberApiHandler) CreateAssignment(c *fiber.Ctx) error {
	c.Status(http.StatusNotImplemented)
	return nil
}

func (handler *FiberApiHandler) getRepo(ctx *fiber.Ctx) repository.Repository {
	repo := ctx.Locals("gitlab-repo").(repository.Repository)
	return repo
}
