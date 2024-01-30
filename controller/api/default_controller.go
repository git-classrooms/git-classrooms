package api

import (
	"de.hs-flensburg.gitlab/gitlab-classroom/context"
	"de.hs-flensburg.gitlab/gitlab-classroom/model"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ApiController struct {
}

func NewApiController() *ApiController {
	return &ApiController{}
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

func (handler *ApiController) CreateClassroom(c *fiber.Ctx) error {
	repo := context.GetGitlabRepository(c)

	var err error
	requestBody := new(CreateClassroomRequest)

	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	group, err := repo.CreateGroup(
		requestBody.Name,
		model.Private,
		requestBody.Description,
		requestBody.MemberEmails,
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

func (handler *ApiController) CreateAssignment(c *fiber.Ctx) error {
	repo := context.GetGitlabRepository(c)

	var err error
	requestBody := new(CreateAssignmentRequest)

	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	templateProject, err := repo.GetProjectById(requestBody.TemplateProjectId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	name := templateProject.Name

	assignees := make([]model.User, len(requestBody.AssigneeUserIds))
	for i, id := range requestBody.AssigneeUserIds {
		user, err := repo.GetUserById(id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		assignees[i] = *user
		name += "_" + user.Username
	}

	project := &model.Project{}
	project, err = repo.ForkProject(requestBody.TemplateProjectId, name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	project, err = repo.AddProjectMembers(project.ID, assignees)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(http.StatusCreated)
}
