package default_controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type createAssignmentRequest struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	TemplateProjectId int        `json:"templateProjectId"`
	DueDate           *time.Time `json:"dueDate" validate:"optional"`
} //@Name CreateAssignmentRequest

func (r createAssignmentRequest) isValid() (bool, string) {
	if r.Name == "" || r.Description == "" || r.TemplateProjectId != 0 || r.DueDate == nil {
		return false, "Request can not be empty, requires name, description, dueDate and templateProjectId"
	}

	if r.DueDate.Before(time.Now()) {
		return false, "DueDate must be in the future"
	}

	return true, ""
}

// @Summary		CreateAssignment
// @Description	CreateAssignment
// @Id				CreateAssignment
// @Tags			assignment
// @Accept			json
// @Param			classroomId		path	string										true	"Classroom ID"	Format(uuid)
// @Param			assignmentInfo	body	default_controller.createAssignmentRequest	true	"Assignment Info"
// @Param			X-Csrf-Token	header	string										true	"Csrf-Token"
// @Success		201
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/classrooms/owned/{classroomId}/assignments [post]
func (ctrl *DefaultController) CreateAssignment(c *fiber.Ctx) error {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()
	classroom := ctx.GetOwnedClassroom()

	requestBody := &createAssignmentRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	isValid, reason := requestBody.isValid()
	if !isValid {
		return fiber.NewError(fiber.StatusBadRequest, reason)
	}

	// Check if template repository exists
	_, err = repo.GetProjectById(requestBody.TemplateProjectId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create assigment
	assignmentQuery := query.Assignment
	assignment := &database.Assignment{
		ClassroomID:       classroom.ID,
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

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
