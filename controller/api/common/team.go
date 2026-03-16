package common

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/logging"
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
	log := logging.GetLogger(ctx)
	var err error
	tx := query.Q.Begin()
	defer tx.Rollback()

	log.Info("adding member to team")

	queryTeam := tx.Team
	newTeam, err := queryTeam.WithContext(ctx).Where(queryTeam.ID.Eq(teamID)).First()
	if err != nil {
		log.Error("error getting team", "teamID", teamID, "error", err)
		return fiber.NewError(fiber.StatusNotFound, err.Error()) // TODO: do not use fiber error
	}
	log = log.With("team", newTeam)

	if err = repo.AddUserToGroup(newTeam.GroupID, member.UserID, model.ReporterPermissions); err != nil {
		log.Error("error adding user to group", "groupID", newTeam.GroupID)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
	}
	defer func() {
		if err != nil {
			if err := repo.RemoveUserFromGroup(newTeam.GroupID, member.UserID); err != nil {
				log.Error("error removing user from group", "groupID", newTeam.GroupID, "memberID", member.UserID, "error", err)
			}
		}
	}()
	log.Info("added user to group", "groupID", newTeam.GroupID)

	member.TeamID = &newTeam.ID
	queryAssignmentProjects := tx.AssignmentProjects
	projects, err := queryAssignmentProjects.
		WithContext(ctx).
		Preload(queryAssignmentProjects.Assignment).
		Where(queryAssignmentProjects.TeamID.Eq(*member.TeamID)).
		Where(queryAssignmentProjects.ProjectStatus.Eq(string(database.Accepted))).
		Find()
	if err != nil {
		log.Error("error getting projects", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
	}

	projectCount := len(projects)
	if projectCount > 0 {
		log.Info(fmt.Sprintf("adding member to %d projects", projectCount))
	}

	for _, project := range projects {
		accessLevel := model.DeveloperPermissions
		if project.Assignment.Closed {
			accessLevel = model.ReporterPermissions
		}
		if err = repo.AddProjectMember(project.ProjectID, member.UserID, accessLevel); err != nil {
			log.Error("error adding member to project", "projectID", project.ProjectID, "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
		}
		defer func() {
			if err != nil {
				if err := repo.RemoveUserFromProject(project.ProjectID, member.UserID); err != nil {
					log.Error("error removing user from project", "projectID", project.ProjectID, "memberID", member.UserID, "error", err)
				}
			}
		}()
		log.Info("added member to project", "projectID", project.ProjectID)
	}

	queryUserClassrooms := tx.UserClassrooms
	if err = queryUserClassrooms.
		WithContext(ctx).
		Save(member); err != nil {
		log.Info("error saving member", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // TODO: do not use fiber error
	}

	err = tx.Commit()
	if err != nil {
		log.Info("error committing transaction", "error", err)
	}
	return err
}
