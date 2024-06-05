package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) ArchiveClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	repo := ctx.GetGitlabRepository()

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
