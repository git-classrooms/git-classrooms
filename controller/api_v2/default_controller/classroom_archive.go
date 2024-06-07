package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		ArchiveClassroom
// @Description	ArchiveClassroom
// @Id				ArchiveClassroom
// @Tags			classroom
// @Produce		json
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
//
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
//
// @Success		202
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/archive [patch]
func (ctrl *DefaultController) ArchiveClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	userClassroom := ctx.GetUserClassroom()
	classroom := userClassroom.Classroom
	repo := ctx.GetGitlabRepository()

	log.Println("ArchiveClassroom")

	if classroom.Archived {
		return c.SendStatus(fiber.StatusNoContent)
	}
	classroom.Archived = true

	oldPermissions := make(map[int]model.AccessLevelValue)
	defer func() {
		if recover() != nil || err != nil {
			for userID, permission := range oldPermissions {
				repo.ChangeUserAccessLevelInGroup(classroom.GroupID, userID, permission)
			}
		}
	}()
	for _, member := range classroom.Member {
		if member.UserID == classroom.OwnerID {
			continue
		}

		permission, err := repo.GetAccessLevelOfUserInGroup(classroom.GroupID, member.UserID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err := repo.ChangeUserAccessLevelInGroup(classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		oldPermissions[member.UserID] = permission
	}

	if _, err := query.Classroom.WithContext(c.Context()).Updates(classroom); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
