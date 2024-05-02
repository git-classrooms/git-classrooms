package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateClassroomRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
} //@Name UpdateClassroomRequest

func (r updateClassroomRequest) isValid() bool {
	return r.Name != "" && r.Description != ""
}

// @Summary		UpdateClassroom
// @Description	UpdateClassroom
// @Id				UpdateClassroom
// @Tags			classroom
// @Accept			json
// @Param			classroomId		path	string										true	"Classroom ID"	Format(uuid)
// @Param 			classroom body default_controller.updateClassroomRequest true "Classroom Update Info"
// @Param			X-Csrf-Token	header	string										true	"Csrf-Token"
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/classrooms/owned/{classroomId} [put]
func (ctrl *DefaultController) PutOwnedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	oldclassroom := classroom

	repo := ctx.GetGitlabRepository()

	requestBody := &updateClassroomRequest{}
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

	_, err = query.Classroom.WithContext(c.Context()).Updates(classroom)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
