package api

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

func (r createAssignmentRequest) isValid() bool {
	return r.Name != "" && r.TemplateProjectId != 0
}

// @Summary		CreateAssignment
// @Description	CreateAssignment
// @Id				CreateAssignmentV2
// @Tags			assignment
// @Accept			json
// @Param			classroomId		path	string						true	"Classroom ID"	Format(uuid)
// @Param			assignmentInfo	body	api.createAssignmentRequest	true	"Assignment Info"
// @Param			X-Csrf-Token	header	string						true	"Csrf-Token"
// @Success		201
// @Header			201	{string}	Location	"/api/v2/classroom/{classroomId}/assignments/{assignmentId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments [post]
func (ctrl *DefaultController) CreateAssignment(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()
	classroom := ctx.GetUserClassroom()

	var requestBody createAssignmentRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	// Check if template repository exists
	if _, err = repo.GetProjectById(requestBody.TemplateProjectId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create assigment
	assignmentQuery := query.Assignment
	assignment := &database.Assignment{
		ClassroomID:       classroom.ClassroomID,
		TemplateProjectID: requestBody.TemplateProjectId,
		Name:              requestBody.Name,
		Description:       requestBody.Description,
		DueDate:           requestBody.DueDate,
	}

	// Persist assigment
	if err = assignmentQuery.WithContext(c.Context()).Create(assignment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s", classroom.ClassroomID.String(), assignment.ID.String()))
	return c.SendStatus(fiber.StatusCreated)
}
