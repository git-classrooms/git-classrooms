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
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
)

// DueAssignmentWork handles the processing of assignments that are due.
// It uses the GitLab API to close assignments after their due date has passed.
type DueAssignmentWork struct {
	gitlabConfig gitlabConfig.Config
}

// NewDueAssignmentWork creates a new instance of DueAssignmentWork with the given GitLab configuration.
func NewDueAssignmentWork(config gitlabConfig.Config) *DueAssignmentWork {
	return &DueAssignmentWork{gitlabConfig: config}
}

// Do processes and closes assignments that are due.
// It fetches assignments, logs into the corresponding GitLab repository, and closes each assignment.
func (w *DueAssignmentWork) Do(ctx context.Context) {
	assignments := w.getAssignments2Close(ctx)
	for _, assignment := range assignments {
		repo, err := GetWorkerRepo(w.gitlabConfig, assignment.Classroom.GroupAccessToken)
		if err != nil {
			log.Default().Printf("Error occurred while login into gitlab: %s", err.Error())
			continue
		}

		err = w.closeAssignment(ctx, assignment, repo)
		if err != nil {
			log.Default().Printf("Error occurred while closing assignment %s: %s", assignment.Name, err.Error())
			continue
		}
	}
}

// getAssignments2Close retrieves assignments that are due and not yet closed from the database.
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

// getLoggedInRepo logs into the GitLab repository associated with the assignment and returns the repository object.
func (w *DueAssignmentWork) getLoggedInRepo(assignment *database.Assignment) (gitlab.Repository, error) {
	repo := gitlab.NewGitlabRepo(w.gitlabConfig)
	err := repo.GroupAccessLogin(assignment.Classroom.GroupAccessToken)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// closeAssignment marks the assignment as closed and performs necessary updates in the repository.
func (w *DueAssignmentWork) closeAssignment(ctx context.Context, assignment *database.Assignment, repo gitlab.Repository) (err error) {
	log.Printf("DueAssignmentWorker: Closing assignment %s", assignment.Name)

	caches := []utils.ProjectAccessLevelCache{}
	defer func() {
		if recover() != nil || err != nil {
			log.Default().Printf("DueAssignmentWorker: Error occurred while closing assignment %s: %s", assignment.Name, err.Error())
			for _, cache := range caches {
				err := repo.ChangeUserAccessLevelInProject(cache.ProjectID, cache.UserID, cache.AccessLevel)
				if err != nil {
					log.Default().Printf("DueAssignmentWorker: Error occurred while changing access level for %d in assignment %s: %s", cache.UserID, assignment.Name, err.Error())
				}
				// TODO: when this fails, we lose the sync between our database and the gitlab. We should handle this in the future
			}
		}
	}()

	for _, project := range assignment.Projects {
		if project.ProjectStatus != database.Accepted {
			continue
		}

		for _, member := range project.Team.Member {
			oldAccessLevel, err := repo.GetAccessLevelOfUserInProject(project.ProjectID, member.UserID)
			if err != nil {
				return err
			}
			if oldAccessLevel == model.OwnerPermissions {
				continue
			}

			if err := repo.ChangeUserAccessLevelInProject(project.ProjectID, member.UserID, model.ReporterPermissions); err != nil {
				return err
			}

			caches = append(caches, utils.ProjectAccessLevelCache{UserID: member.UserID, ProjectID: project.ProjectID, AccessLevel: oldAccessLevel})
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
