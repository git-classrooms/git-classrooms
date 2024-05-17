package default_controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateAssignmentRequest struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

func (r updateAssignmentRequest) isValid() bool {
	return r.Name != "" || r.Description != "" || r.DueDate != nil
}

func (ctrl *DefaultController) PutOwnedAssignments(c *fiber.Ctx) error {
	ctx := context.Get(c)
	assignment := ctx.GetOwnedClassroomAssignment()
	oldName := assignment.Name
	oldDescription := assignment.Description

	repo := ctx.GetGitlabRepository()
	var err error

	requestBody := &updateAssignmentRequest{}
	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "request requires name, description or dueDate")
	}

	if requestBody.Name != "" {
		for _, project := range assignment.Projects {
			_, err = repo.ChangeProjectName(project.ProjectID, requestBody.Name)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			defer func(projectId int) {
				if recover() != nil || err != nil {
					repo.ChangeProjectName(projectId, oldName)
				}
			}(project.ProjectID)
		}

		assignment.Name = requestBody.Name
	}

	if requestBody.Description != "" {
		for _, project := range assignment.Projects {
			_, err = repo.ChangeProjectDescription(project.ProjectID, requestBody.Description)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			defer func(projectId int) {
				if recover() != nil || err != nil {
					repo.ChangeProjectDescription(projectId, oldDescription)
				}
			}(project.ProjectID)
		}

		assignment.Description = requestBody.Description
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
