package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/gen/field"
)

// @Summary		ArchiveClassroom
// @Description	ArchiveClassroom
// @Id				ArchiveClassroom
// @Tags			classroom
// @Produce		json
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
// @Success		202
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/archive [patch]
func (ctrl *DefaultController) ArchiveClassroom(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	userClassroom := ctx.GetUserClassroom()
	classroom := userClassroom.Classroom
	repo := ctx.GetGitlabRepository()
	log := ctx.GetLoggerForHandler("ArchiveClassroom")

	if classroom.Archived {
		log.Info("classroom is already archived")
		return c.SendStatus(fiber.StatusNoContent)
	}
	classroom.Archived = true

	teams, err := query.Team.
		WithContext(c.Context()).
		Preload(query.Team.Member).
		Preload(query.Team.AssignmentProjects).
		Preload(field.NewRelation("AssignmentProjects.Assignment", "")).
		Where(query.Team.ClassroomID.Eq(classroom.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, team := range teams {
		for _, project := range team.AssignmentProjects {
			if project.ProjectStatus == database.Creating {
				return fiber.NewError(fiber.StatusConflict, fmt.Sprintf("project %q in team %q is still being created, please try again in a few minutes", project.Assignment.Name, team.Name))
			}
		}
	}

	caches := []utils.ProjectAccessLevelCache{}
	defer func() {
		if recover() != nil || err != nil {
			for _, cache := range caches {
				repo.ChangeUserAccessLevelInProject(cache.ProjectID, cache.UserID, cache.AccessLevel)
			}
		}
	}()
	for _, team := range teams {
		for _, project := range team.AssignmentProjects {
			if project.ProjectStatus != database.Accepted {
				continue
			}
			for _, member := range team.Member {
				if member.Role != database.Student {
					continue
				}

				permission, err := repo.GetAccessLevelOfUserInProject(project.ProjectID, member.UserID)
				if err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}

				if err := repo.ChangeUserAccessLevelInProject(project.ProjectID, member.UserID, model.ReporterPermissions); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}

				caches = append(caches, utils.ProjectAccessLevelCache{UserID: member.UserID, ProjectID: project.ProjectID, AccessLevel: permission})
			}
		}
	}

	if _, err := query.Classroom.WithContext(c.Context()).Updates(classroom); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
