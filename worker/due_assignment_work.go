package worker

import (
	"context"
	"log"
	"time"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type DueAssignmentWork struct {
	gitlabConfig gitlabConfig.Config
}

func NewDueAssignmentWork(config gitlabConfig.Config) *DueAssignmentWork {
	return &DueAssignmentWork{gitlabConfig: config}
}

func (w *DueAssignmentWork) Do(ctx context.Context) {
	assignments := w.getAssignments2Close(ctx)
	for _, assignment := range assignments {
		repo, err := GetWorkerRepo(w.gitlabConfig, assignment.Classroom.GroupAccessToken)
		if err != nil {
			log.Default().Printf("Error occurred while login into gitlab: %s", err.Error())
			continue
		}

		w.closeAssignment(ctx, assignment, repo)
	}
}

func (w *DueAssignmentWork) getAssignments2Close(ctx context.Context) []*database.Assignment {
	assignments, err := query.Assignment.
		WithContext(ctx).
		Preload(query.Assignment.Projects).
		Preload(query.Assignment.Projects.Team).
		Preload(query.Assignment.Projects.Team.Member).
		Preload(query.Assignment.Classroom).
		Where(query.Assignment.DueDate.Lt(time.Now())).
		Where(query.Assignment.Closed.Is(false)).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching assignments to close: %s", err.Error())
		return []*database.Assignment{}
	}

	return assignments
}

type RestoreCache struct {
	userID        int
	projectID     int
	oldPermission model.AccessLevelValue
}

func (w *DueAssignmentWork) closeAssignment(ctx context.Context, assignment *database.Assignment, repo gitlab.Repository) (err error) {
	log.Printf("DueAssignmentWorker: Closing assignment %s", assignment.Name)

	caches := []RestoreCache{}
	defer func() {
		if recover() != nil || err != nil {
			log.Default().Printf("DueAssignmentWorker: Error occurred while closing assignment %s: %s", assignment.Name, err.Error())
			for _, cache := range caches {
				repo.ChangeUserAccessLevelInProject(cache.projectID, cache.userID, cache.oldPermission)
				// TODO: when this fails, we lose the sync between our database and the gitlab. We should handle this in the future
			}
		}
	}()

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

			caches = append(caches, RestoreCache{userID: member.UserID, projectID: project.ProjectID, oldPermission: oldPermission})
		}
	}

	assignment.Closed = true
	_, err = query.Assignment.WithContext(ctx).Updates(assignment)
	if err != nil {
		return err
	}

	log.Printf("DueAssignmentWorker: Assignment %s has been closed", assignment.Name)
	return nil
}
