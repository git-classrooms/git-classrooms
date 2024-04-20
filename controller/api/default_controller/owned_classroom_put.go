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
	oldclassroom := classroom

	repo := ctx.GetGitlabRepository()

	requestBody := &UpdateClassroomRequest{}
	err := c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "request requires name and description")
	}

	if requestBody.Name != oldclassroom.Name {
		group, err := repo.ChangeGroupName(classroom.GroupID, requestBody.Name)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		defer func() {
			if recover() != nil || err != nil {
				repo.ChangeGroupName(classroom.GroupID, oldclassroom.Name)
			}
		}()

		classroom.Name = group.Name
	}

	if requestBody.Description != oldclassroom.Description {
		group, err := repo.ChangeGroupDescription(classroom.GroupID, requestBody.Description)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		defer func() {
			if recover() != nil || err != nil {
				repo.ChangeGroupDescription(classroom.GroupID, oldclassroom.Description)
			}
		}()

		classroom.Description = group.Description
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		_, err := query.Classroom.WithContext(c.Context()).Updates(classroom)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return err
	})

	c.SendStatus(fiber.StatusAccepted)
	return c.JSON(classroom)
}
