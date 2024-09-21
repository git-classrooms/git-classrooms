package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
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
// @Id				UpdateClassroomV2
// @Tags			classroom
// @Accept			json
// @Param			classroomId		path	string						true	"Classroom ID"	Format(uuid)
// @Param			classroom		body	api.updateClassroomRequest	true	"Classroom Update Info"
// @Param			X-Csrf-Token	header	string						true	"Csrf-Token"
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId} [put]
func (ctrl *DefaultController) UpdateClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom().Classroom
	oldclassroom := ctx.GetUserClassroom().Classroom

	repo := ctx.GetGitlabRepository()

	var requestBody updateClassroomRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	if requestBody.Name != oldclassroom.Name {
		classroom.Name = requestBody.Name
		if _, err := repo.ChangeGroupName(classroom.GroupID, requestBody.Name); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		defer func() {
			if recover() != nil || err != nil {
				repo.ChangeGroupName(classroom.GroupID, oldclassroom.Name)
			}
		}()
	}

	if requestBody.Description != oldclassroom.Description {
		classroom.Description = requestBody.Description

		if _, err := repo.ChangeGroupDescription(classroom.GroupID, utils.CreateClassroomGitlabDescription(&classroom, ctrl.config.PublicURL)); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		defer func() {
			if recover() != nil || err != nil {
				repo.ChangeGroupDescription(classroom.GroupID, oldclassroom.Description)
			}
		}()
	}

	if _, err = query.Classroom.WithContext(c.Context()).Updates(classroom); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.SendStatus(fiber.StatusAccepted)
	return c.JSON(classroom)
}
