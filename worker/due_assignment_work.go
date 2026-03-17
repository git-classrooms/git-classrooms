package worker

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
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

		err = w.closeAssignmentDate(ctx, assignmentDate, repo)
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
		Preload(field.NewRelation("AssignmentProjectGradingDate.AssignmentProject.Team.Member.User", "")).
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

// closeAssignmentDate marks the assignment as closed and performs necessary updates in the repository.
func (w *DueAssignmentWork) closeAssignmentDate(ctx context.Context, assignmentDate *database.AssignmentDate, repo gitlab.Repository) (err error) {
	log.Printf("DueAssignmentWorker: Closing assignmentDate %s of %s", assignmentDate.DueDate.String(), assignmentDate.Assignment.Name)

	tx := query.Q.Begin()
	defer tx.Rollback()

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

		_, err := repo.CreateBranch(project.ProjectID, branchName, "main")
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

		mentions := utils.Map(project.Team.Member, func(member *database.UserClassrooms) string {
			return fmt.Sprintf("/cc @%s", member.User.GitlabUsername)
		})
		description := fmt.Sprintf(mergeRequestDescription, branchName, strings.Join(mentions, "\n"))

		var userID *int
		if len(project.Team.Member) > 0 {
			userID = &project.Team.Member[0].UserID
		}

		if err = repo.CreateMergeRequest(project.ProjectID, branchName, previousBranchName, fmt.Sprintf("Feedback (%s)", branchName), description, userID, assignmentDate.Assignment.Classroom.OwnerID); err != nil {
			return err
		}
		defer func() {
			if err != nil {
				if otherErr := repo.DeleteMergeRequest(project.ProjectID, branchName, previousBranchName); otherErr != nil {
					fmt.Println(otherErr)
				}
			}
		}()

		projectGradings.Branchname = branchName
		if err = tx.AssignmentProjectGradingDate.WithContext(ctx).Save(projectGradings); err != nil {
			return err
		}
	}

	if idx == len(otherDatesSorted)-1 {

		caches := []utils.ProjectAccessLevelCache{}
		defer func() {
			if err != nil {
				log.Default().Printf("DueAssignmentWorker: Error occurred while closing assignment %s of %s: %s", assignmentDate.DueDate.String(), assignmentDate.Assignment.Name, err.Error())
				for _, cache := range caches {
					err := repo.ChangeUserAccessLevelInProject(cache.ProjectID, cache.UserID, cache.AccessLevel)
					if err != nil {
						log.Default().Printf("DueAssignmentWorker: Error occurred while changing access level for %d in assignment %s of %s: %s", cache.UserID, assignmentDate.DueDate, assignmentDate.Assignment.Name, err.Error())
					}
					// TODO: when this fails, we lose the sync between our database and the gitlab. We should handle this in the future
				}
			}
		}()

		for _, projectGrading := range assignmentDate.AssignmentProjectGradingDate {
			project := projectGrading.AssignmentProject
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

	}

	assignmentDate.Closed = true
	_, err = tx.Assignment.WithContext(ctx).Updates(assignmentDate)
	if err != nil {
		return err
	}

	log.Printf("DueAssignmentWorker: Assignment %s of %s has been closed", assignmentDate.DueDate, assignmentDate.Assignment.Name)
	err = tx.Commit()
	return err
}

const (
	mergeRequestDescription string = `
👋! GitLab Classroom created this merge request as a place for your teacher to leave feedback on your work. It will update automatically. **Don't close or merge this merge request**, unless you're instructed to do so by your teacher.
In this merge request, your teacher can leave comments and feedback on your code.
Click the **Changes** or **Commits** tab to see all of the changes pushed until the duedate ` + "`%s`" + ` for this assignmentDate was reached. Your teacher can see this too.

<details>
<summary>
<strong>Notes for teachers</strong>
</summary>

Use this MR to leave feedback. Here are some tips:
  - Click the **Changes** tab to see all of the changes pushed to ` + "`main`" + `since the assignment started. To leave comments on specific lines of code, put your cursor over a line of code and click the blue **comment sign**. To learn more about comments, read "[Add a comment to a merge request diff](https://docs.gitlab.com/ee/user/discussions/#add-a-comment-to-a-merge-request-diff)".
  - Click the **Commits** tab to see the commits pushed to ` + "`main`" + `. Click a commit to see specific changes.
  - ?? If you turned on autograding, then click the **Checks** tab to see the results. ??
  - This page is an overview. It shows commits, line comments, and general comments. You can leave a general comment below.

</details>

%s
`
)
