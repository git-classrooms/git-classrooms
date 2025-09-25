package common

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

func AddToTeam(
	ctx context.Context,
	repo gitlab.Repository,
	member *database.UserClassrooms,
	teamID uuid.UUID,
) error {
	var err error
	tx := query.Q.Begin()
	defer tx.Rollback()

	queryTeam := tx.Team
	newTeam, err := queryTeam.WithContext(ctx).Where(queryTeam.ID.Eq(teamID)).First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error()) // TODO: do not use fiber error
	}

	if err = repo.AddUserToGroup(newTeam.GroupID, member.UserID, model.ReporterPermissions); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
	}
	defer func() {
		if err != nil {
			if err := repo.RemoveUserFromGroup(newTeam.GroupID, member.UserID); err != nil {
				log.Println(err)
			}
		}
	}()

	member.TeamID = &newTeam.ID
	queryAssignmentProjects := tx.AssignmentProjects
	projects, err := queryAssignmentProjects.
		WithContext(ctx).
		Preload(queryAssignmentProjects.Assignment).
		Where(queryAssignmentProjects.TeamID.Eq(*member.TeamID)).
		Where(queryAssignmentProjects.ProjectStatus.Eq(string(database.Accepted))).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
	}

	for _, project := range projects {
		accessLevel := model.DeveloperPermissions
		if project.Assignment.Closed {
			accessLevel = model.ReporterPermissions
		}
		if err = repo.AddProjectMember(project.ProjectID, member.UserID, accessLevel); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
		}

		defer func() {
			if err != nil {
				if err := repo.RemoveUserFromProject(project.ProjectID, member.UserID); err != nil {
					log.Println(err)
				}
			}
		}()
	}

	queryUserClassrooms := tx.UserClassrooms
	if err = queryUserClassrooms.
		WithContext(ctx).
		Save(member); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
	}

	err = tx.Commit()
	return err
}
