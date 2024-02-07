package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"time"
)

type CreateAssignmentRequest struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	TemplateProjectId int        `json:"templateProjectId"`
	DueDate           *time.Time `json:"dueDate"`
}

func (r CreateAssignmentRequest) isValid() bool {
	return r.Name != "" && r.TemplateProjectId != 0
}

func (ctrl *DefaultController) CreateAssignment(c *fiber.Ctx) error {
	repo := context.Get(c).GetGitlabRepository()

	// Parse parameters
	classroomId, err := uuid.Parse(c.Params("classroomId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Parse body
	requestBody := new(CreateAssignmentRequest)
	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Check if template repository exists
	_, err = repo.GetProjectById(requestBody.TemplateProjectId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create assigment
	assignmentQuery := query.Assignment
	assignment := &database.Assignment{
		ClassroomID:       classroomId,
		TemplateProjectID: requestBody.TemplateProjectId,
		Name:              requestBody.Name,
		Description:       requestBody.Description,
		DueDate:           requestBody.DueDate,
	}

	// Persist assigment
	err = assignmentQuery.WithContext(c.Context()).Create(assignment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s", classroomId.String(), assignment.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
