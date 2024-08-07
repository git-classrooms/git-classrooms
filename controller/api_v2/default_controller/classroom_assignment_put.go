package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateAssignmentRequest struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
} //@Name UpdateAssignmentRequest

func (r updateAssignmentRequest) isValid() (bool, string) {
	if r.Name == "" && r.Description == "" && r.DueDate == nil {
		return false, "Request can not be empty, requires name, description or dueDate"
	}

	if r.DueDate.Before(time.Now()) {
		return false, "DueDate must be in the future"
	}

	return true, ""
}

// @Summary		UpdateAssignment
// @Description	UpdateAssignment
// @Id				UpdateAssignment
// @Tags			assignment
// @Accept			json
// @Param			classroomId		path		string						true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string						true	"Assignment ID"	Format(uuid)
// @Param			assignmentInfo	body		api.updateAssignmentRequest	true	"Assignment Update Info"
// @Param			X-Csrf-Token	header		string						true	"Csrf-Token"
// @Success		202				{object}	Assignment
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		403				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId} [put]
func (ctrl *DefaultController) UpdateAssignment(c *fiber.Ctx) error {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()
	var err error

	requestBody := &updateAssignmentRequest{}
	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	isValid, reason := requestBody.isValid()
	if !isValid {
		return fiber.NewError(fiber.StatusBadRequest, reason)
	}

	projectLinks, err := query.AssignmentProjects.WithContext(c.Context()).Where(query.AssignmentProjects.AssignmentID.Eq(assignment.ID)).Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	assignmentAcceptedByStudents := false
	for _, projectLink := range projectLinks {
		if projectLink.ProjectStatus == database.Accepted || projectLink.ProjectStatus == database.Creating {
			assignmentAcceptedByStudents = true
			break
		}
	}

	if assignmentAcceptedByStudents {
		if requestBody.Name != "" || requestBody.Description != "" {
			return fiber.NewError(fiber.StatusBadRequest, "Assignment name and description can not be changed after it has been accepted by students")
		}
	} else {
		if requestBody.Name != "" {
			assignment.Name = requestBody.Name
		}

		if requestBody.Description != "" {
			assignment.Description = requestBody.Description
		}
	}

	if requestBody.DueDate != nil {
		assignment.DueDate = requestBody.DueDate
	}

	if _, err = query.Assignment.WithContext(c.Context()).Updates(assignment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.SendStatus(fiber.StatusAccepted)
	return c.JSON(assignment)
}
