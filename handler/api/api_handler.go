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

type CreateClassroomRequest struct {
	Name         string   `json:"name"`
	MemberEmails []string `json:"memberEmails"`
	Description  string   `json:"description"`
}

type CreateAssignmentRequest struct {
	AssigneeUserIds   []int `json:"assigneeUserIds"`
	TemplateProjectId int   `json:"templateProjectId"`
}

func (handler *FiberApiHandler) CreateClassroom(c *fiber.Ctx) error {
	repo := handler.getRepo(c)

	var err error
	requestBody := new(CreateClassroomRequest)

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
	repo := handler.getRepo(c)

	var err error
	requestBody := new(CreateAssignmentRequest)

	err = c.BodyParser(requestBody)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	templateProject, err := repo.GetProjectById(requestBody.TemplateProjectId)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	name := templateProject.Name

	assignees := make([]model.User, len(requestBody.AssigneeUserIds))
	for i, id := range requestBody.AssigneeUserIds {
		user, err := repo.GetUserById(id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}
		assignees[i] = *user
		name += "_" + user.Username
	}

	project := &model.Project{}
	project, err = repo.ForkProject(requestBody.TemplateProjectId, name)
	if err != err {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	project, err = repo.AddProjectMembers(project.ID, assignees)
	if err != err {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	c.Status(http.StatusCreated)
	return nil
}

func (handler *FiberApiHandler) getRepo(ctx *fiber.Ctx) repository.Repository {
	repo := ctx.Locals("gitlab-repo").(repository.Repository)
	return repo
}
