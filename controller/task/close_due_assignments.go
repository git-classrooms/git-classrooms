package task

import (
	"context"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type OldPermission struct {
	UserID     int
	ProjectID  int
	Permission model.AccessLevelValue
}

func CloseDueAssignments(repo gitlab.Repository, ctx context.Context) (err error) {
	assignments, err := query.Assignment.
		WithContext(ctx).
		Preload(query.Assignment.Projects).
		Preload(query.Assignment.Projects.Team).
		Preload(query.Assignment.Projects.Team.Member).
		Where(query.Assignment.DueDate.Lt(time.Now())).
		Where(query.Assignment.Closed.Is(false)).
		Find()
	if err != nil {
		return err
	}

	oldPermissions := []OldPermission{}
	defer func() {
		if recover() != nil || err != nil {
			for _, oldPermission := range oldPermissions {
				repo.ChangeUserAccessLevelInProject(oldPermission.ProjectID, oldPermission.UserID, oldPermission.Permission)
			}
		}
	}()

	for _, assignment := range assignments {
		for _, project := range assignment.Projects {
			if project.ProjectStatus != database.Accepted {
				continue
			}

			for _, member := range project.Team.Member {
				oldPermission, err := repo.GetAccessLevelOfUserInProject(project.ProjectID, member.UserID)
				if err != nil {
					return err
				}

				if err := repo.ChangeUserAccessLevelInProject(project.ProjectID, member.UserID, model.ReporterPermissions); err != nil {
					return err
				}

				oldPermissions = append(oldPermissions, OldPermission{UserID: member.UserID, ProjectID: project.ProjectID, Permission: oldPermission})
			}
		}

		assignment.Closed = true
		_, err := query.Assignment.WithContext(ctx).Updates(assignment)
		if err != nil {
			return err
		}
		oldPermissions = []OldPermission{}
	}

	return nil
}
