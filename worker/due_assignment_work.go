package worker

import (
	"context"
	"fmt"
	"time"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/logging"
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
	log := logging.GetLogger(ctx)

	log.Info("starting work for closed assignments")

	assignments := w.getAssignments2Close(ctx)
	assignmentCount := len(assignments)
	if assignmentCount > 0 {
		log.Info(fmt.Sprintf("%d assignments to close", assignmentCount))
	}

	for _, assignment := range assignments {
		log = log.With("assignment", assignment)
		ctx = logging.SetLogger(ctx, log)

		repo, err := GetWorkerRepo(ctx, w.gitlabConfig, assignment.Classroom.GroupAccessToken)
		if err != nil {
			log.Error("error occurred while creating repo", "error", err)
			continue
		}

		err = w.closeAssignment(ctx, assignment, repo)
		if err != nil {
			log.Error("error occurred while closing assignment", "error", err)
			continue
		}
	}
}

// getAssignments2Close retrieves assignments that are due and not yet closed from the database.
func (w *DueAssignmentWork) getAssignments2Close(ctx context.Context) []*database.Assignment {
	log := logging.GetLogger(ctx)
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
		log.Error("error occurred while fetching assignments to close", "error", err)
		return []*database.Assignment{}
	}

	return assignments
}

// getLoggedInRepo logs into the GitLab repository associated with the assignment and returns the repository object.
func (w *DueAssignmentWork) getLoggedInRepo(ctx context.Context, assignment *database.Assignment) (gitlab.Repository, error) {
	log := logging.GetLogger(ctx)
	repo := gitlab.NewGitlabRepo(w.gitlabConfig)
	err := repo.GroupAccessLogin(assignment.Classroom.GroupAccessToken, log)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// closeAssignment marks the assignment as closed and performs necessary updates in the repository.
func (w *DueAssignmentWork) closeAssignment(ctx context.Context, assignment *database.Assignment, repo gitlab.Repository) (err error) {
	log := logging.GetLogger(ctx)
	log.Info("closing assignment")

	caches := []utils.ProjectAccessLevelCache{}
	defer func() {
		if recover() != nil || err != nil {
			log.Error("error occurred while closing assignment", "error", err)
			for _, cache := range caches {
				err := repo.ChangeUserAccessLevelInProject(cache.ProjectID, cache.UserID, cache.AccessLevel)
				if err != nil {
					log.Error("error occurred while changing access level for member", "projectID", cache.ProjectID, "memberID", cache.UserID, "error", err)
				}
				// TODO: when this fails, we lose the sync between our database and the gitlab. We should handle this in the future
			}
		}
	}()

	for _, project := range assignment.Projects {
		if project.ProjectStatus != database.Accepted {
			continue
		}
		log := log.With("project", project, "team", project.Team)

		log.Info("closing assignmentProject")

		for _, member := range project.Team.Member {
			log = log.With("projectID", project.ProjectID, "memberID", member.UserID)
			log.Debug("getting access-level of member")
			oldAccessLevel, err := repo.GetAccessLevelOfUserInProject(project.ProjectID, member.UserID)
			if err != nil {
				return err
			}
			if oldAccessLevel == model.OwnerPermissions {
				continue
			}
			log.Debug("changing access-level of member to reporter")

			if err := repo.ChangeUserAccessLevelInProject(project.ProjectID, member.UserID, model.ReporterPermissions); err != nil {
				return err
			}

			caches = append(caches, utils.ProjectAccessLevelCache{UserID: member.UserID, ProjectID: project.ProjectID, AccessLevel: oldAccessLevel})
		}
	}

	log.Debug("setting assignment to closed")

	assignment.Closed = true
	_, err = query.Assignment.WithContext(ctx).Updates(assignment)
	if err != nil {
		return err
	}

	log.Info("assignment has been closed")
	return nil
}
