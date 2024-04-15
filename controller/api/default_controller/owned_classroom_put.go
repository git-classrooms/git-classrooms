package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type UpdateClassroomRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r UpdateClassroomRequest) isValid() bool {
	return r.Name != "" || r.Description != ""
}

func (ctrl *DefaultController) UpdateClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	oldClassroom := *classroom
	classroomQuery := query.Classroom

	repo := ctx.GetGitlabRepository()

	requestBody := &UpdateClassroomRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if requestBody.Name != "" {
		group, err := repo.ChangeGroupName(classroom.GroupID, requestBody.Name)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		classroom.Name = group.Name

	}

	if requestBody.Description != "" {
		group, err := repo.ChangeGroupDescription(classroom.GroupID, requestBody.Description)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		classroom.Description = group.Description
	}

	if requestBody.Name != "" || requestBody.Description != "" {
		info, err := classroomQuery.WithContext(c.Context()).Updates(classroom)
		if err != nil || info.Error != nil {
			if requestBody.Name != "" {
				_, newErr := repo.ChangeGroupName(oldClassroom.GroupID, oldClassroom.Name)
				if newErr != nil {
					return fiber.NewError(fiber.StatusInternalServerError, newErr.Error())
				}
			}
			if requestBody.Description != "" {
				_, newErr := repo.ChangeGroupDescription(oldClassroom.GroupID, oldClassroom.Description)
				if newErr != nil {
					return fiber.NewError(fiber.StatusInternalServerError, newErr.Error())
				}
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(classroom)
}
