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
	return r.Name != "" && r.Description != ""
}

func (ctrl *DefaultController) PutOwnedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	oldName := classroom.Name
	oldDescription := classroom.Description

	repo := ctx.GetGitlabRepository()

	requestBody := &UpdateClassroomRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "request requires name and description")
	}

	if requestBody.Name != oldName {
		group, err := repo.ChangeGroupName(classroom.GroupID, requestBody.Name)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		classroom.Name = group.Name
	}

	if requestBody.Description != oldDescription {
		group, err := repo.ChangeGroupDescription(classroom.GroupID, requestBody.Description)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		classroom.Description = group.Description
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		defer func() {
			if recover() != nil || err != nil {
				repo.ChangeGroupName(classroom.GroupID, oldName)
				repo.ChangeGroupDescription(classroom.GroupID, oldDescription)
			}
		}()

		_, err := query.Classroom.WithContext(c.Context()).Updates(classroom)
		return err
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.SendStatus(fiber.StatusAccepted)
	return c.JSON(classroom)
}
