package worker

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gorm.io/gen/field"
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
	assignmentDates := w.getAssignmentDates2Close(ctx)
	for _, assignmentDate := range assignmentDates {
		repo, err := GetWorkerRepo(w.gitlabConfig, assignmentDate.Assignment.Classroom.GroupAccessToken)
		if err != nil {
			log.Default().Printf("Error occurred while login into gitlab: %s", err.Error())
			continue
		}

		err = w.closeAssignment(ctx, assignmentDate, repo)
		if err != nil {
			log.Default().Printf("Error occurred while closing assignment date %s of %s: %s", assignmentDate.DueDate.String(), assignmentDate.Assignment.Name, err.Error())
			continue
		}
	}
}

// getAssignments2Close retrieves assignments that are due and not yet closed from the database.
func (w *DueAssignmentWork) getAssignmentDates2Close(ctx context.Context) []*database.AssignmentDate {
	assignmentDates, err := query.AssignmentDate.
		WithContext(ctx).
		Preload(query.AssignmentDate.AssignmentProjectGradingDate).
		Preload(query.AssignmentDate.AssignmentProjectGradingDate.AssignmentProject.Team).
		Preload(query.AssignmentDate.AssignmentProjectGradingDate.AssignmentProject.Team.Member).
		Preload(query.AssignmentDate.Assignment).
		Preload(field.NewRelation("Assignment.Classroom", "")).
		Preload(field.NewRelation("Assignment.AssignmentDate", "")).
		Where(query.AssignmentDate.DueDate.Lt(time.Now())).
		Where(query.AssignmentDate.Closed.Is(false)).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching assignments to close: %s", err.Error())
		return []*database.AssignmentDate{}
	}

	return assignmentDates
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
func (w *DueAssignmentWork) closeAssignment(ctx context.Context, assignmentDate *database.AssignmentDate, repo gitlab.Repository) (err error) {
	log.Printf("DueAssignmentWorker: Closing assignmentDate %s of %s", assignmentDate.DueDate.String(), assignmentDate.Assignment.Name)

	otherDates := assignmentDate.Assignment.AssignmentDates
	otherDatesSorted := slices.SortedFunc(slices.Values(otherDates), func(a, b *database.AssignmentDate) int {
		return a.DueDate.Compare(b.DueDate)
	})

	idx := slices.IndexFunc(otherDates, func(a *database.AssignmentDate) bool {
		return a.ID == assignmentDate.ID
	})

	if idx < 0 {
		panic("unreachable, assignmentDate is not in assignmentDates")
	}

	branchName := assignmentDate.DueDate.Format(time.DateOnly)
	previousBranchName := "init"
	if idx != 0 {
		previousBranchName = otherDatesSorted[idx-1].DueDate.Format(time.DateOnly)
	}

	for _, projectGradings := range assignmentDate.AssignmentProjectGradingDate {
		project := projectGradings.AssignmentProject

		branch, err := repo.CreateBranch(project.ProjectID, branchName, "main")
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				if otherErr := repo.DeleteBranch(project.ProjectID, branchName); otherErr != nil {
					fmt.Println(otherErr)
				}
			}
		}()

		if err = repo.ProtectBranch(project.ProjectID, branchName, model.MaintainerPermissions); err != nil {
			return err
		}
		defer func() {
			if err != nil {
				if otherErr := repo.UnprotectBranch(project.ProjectID, branchName); otherErr != nil {
					fmt.Println(otherErr)
				}
			}
		}()

		repo.CreateMergeRequest(project.ProjectID, branchName)
	}

	// TODO: create branch, protect it, create MR

	if idx == len(otherDatesSorted)-1 {
		// TODO: remove access
		panic("TODO: last assignmentDate")
	}

	// TODO: this needs to be reworked completely
	//       if it is not the last assignmentDate we need to create protected branches with MR to the one before
	//       when its the last, only then we will change the permissions of the user and create a last branch and MR

	assignmentDate.Closed = true
	_, err = query.Assignment.WithContext(ctx).Updates(assignmentDate)
	if err != nil {
		return err
	}

	log.Printf("DueAssignmentWorker: Assignment %s of %s has been closed", assignmentDate.DueDate, assignmentDate.Assignment.Name)
	return nil
}
