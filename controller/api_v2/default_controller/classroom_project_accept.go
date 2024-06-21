package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Accept the assignment
// @Description	Accept the assignment and work on your repository
// @Id				AcceptAssignmentV2
// @Tags			project
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
// @Param			projectId		path	string	true	"Project ID"	Format(uuid)
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
// @Success		201
// @Success		202
// @Header			202	{string}	Location	"/api/v2/classroom/{classroomId}/projects/{projectId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/projects/{projectId}/accept [post]
func (ctrl *DefaultController) AcceptAssignment(c *fiber.Ctx) (err error) {
	ctx := fiberContext.Get(c)
	classroom := ctx.GetUserClassroom()
	userID := ctx.GetUserID()
	assignmentProject := ctx.GetAssignmentProject()

	if assignmentProject.ProjectStatus == database.Accepted {
		return c.SendStatus(fiber.StatusNoContent) // You or your teammate have already accepted the assignment
	}

	if assignmentProject.ProjectStatus == database.Creating {
		return fiber.NewError(fiber.StatusForbidden, "The project is still being created")
	}

	if assignmentProject.Assignment.DueDate != nil && assignmentProject.Assignment.DueDate.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "The assignment is already over")
	}

	repo := ctx.GetGitlabRepository()

	if err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Check if template repository still exists
	templateProject, err := repo.GetProjectById(assignmentProject.Assignment.TemplateProjectID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignmentProject.ProjectStatus = database.Creating

	queryAssignmentProjects := query.AssignmentProjects
	if err = queryAssignmentProjects.WithContext(c.Context()).Save(assignmentProject); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Make the actual project creation and assignment acceptance async
	go ctrl.acceptAssignment(repo, userID, classroom.Classroom.OwnerID, templateProject, assignmentProject)

	c.Set("Location", fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s", classroom.ClassroomID.String(), assignmentProject.AssignmentID.String()))
	return c.SendStatus(fiber.StatusAccepted)
}

const (
	mergeRequestDescription string = `
👋! GitLab Classroom created this merge request as a place for your teacher to leave feedback on your work. It will update automatically. **Don't close or merge this merge request**, unless you're instructed to do so by your teacher.
In this merge request, your teacher can leave comments and feedback on your code.
Click the **Changes** or **Commits** tab to see all of the changes pushed to ` + "`main`" + ` since the assignment started. Your teacher can see this too.

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

func (ctrl *DefaultController) acceptAssignment(repo gitlab.Repository, userID int, classroomOwnerID int, templateProject *gitlabModel.Project, assignmentProject *database.AssignmentProjects) {
	fmt.Println("Template defaultbranch ", templateProject.DefaultBranch)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	queryAssignmentProjects := query.AssignmentProjects

	var err error
	defer func() {
		if recover() != nil || err != nil {
			assignmentProject.ProjectStatus = database.Failed
			if err := queryAssignmentProjects.WithContext(ctx).
				Save(assignmentProject); err != nil {
				log.Println("Error while setting Project to Failed!", err)
			}
		}
	}()

	project, err := repo.ForkProject(assignmentProject.Assignment.TemplateProjectID, gitlabModel.Private, assignmentProject.Team.GroupID, assignmentProject.Assignment.Name, assignmentProject.Assignment.Description)
	if err != nil {
		log.Println("Error while forking the template Project", err)
		return
	}
	defer func() {
		if recover() != nil || err != nil {
			if err := repo.DeleteProject(project.ID); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	// wait till default branch of forked project is the same as the template project
	// this is necessary because the default branch is not immediately available after forking
	// TODO?: do not wait the whole 5 Minutes for this
	err = waitForDefaultBranch(ctx, repo, project.ID, templateProject.DefaultBranch)
	if err != nil {
		log.Println("Error while waiting for defaultBranch", err)
		return
	}

	memberIds := utils.Map(assignmentProject.Team.Member, func(member *database.UserClassrooms) int {
		return member.UserID
	})

	gitlabMember := utils.Map(memberIds, func(member int) gitlabModel.User {
		return gitlabModel.User{ID: member}
	})

	project, err = repo.AddProjectMembers(project.ID, gitlabMember)
	if err != nil {
		log.Println("Error while adding members to the project", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	_, err = repo.CreateBranch(project.ID, "feedback", project.DefaultBranch)
	if err != nil {
		log.Println("Error while creating feedback branch", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	queryUsers := query.User
	members, err := queryUsers.
		WithContext(ctx).
		Where(queryUsers.ID.In(memberIds...)).
		Find()
	if err != nil {
		log.Println("Error while fetching members", err)
		return
	}

	mentions := utils.Map(members, func(member *database.User) string {
		return fmt.Sprintf("/cc @%s", member.GitlabUsername)
	})
	description := fmt.Sprintf(mergeRequestDescription, strings.Join(mentions, "\n"))
	err = repo.CreateMergeRequest(project.ID, project.DefaultBranch, "feedback", "Feedback", description, userID, classroomOwnerID)
	if err != nil {
		log.Println("Error while creating merge request", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	// In a few cases the main branch isn't available directly after the creation, this would cause an error when setting up protection rules for it, there we wait for the default branch to exist
	// TODO?: do not wait the whole 5 Minutes for this
	err = waitForProtectedBranch(ctx, repo, project.ID, project.DefaultBranch)
	if err != nil {
		log.Println("Error while waiting for protected main branch", err)
		return
	}

	err = repo.UnprotectBranch(project.ID, project.DefaultBranch)
	if err != nil {
		log.Println("Error while unprotecting default branch", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	err = repo.ProtectBranch(project.ID, project.DefaultBranch, gitlabModel.DeveloperPermissions)
	if err != nil {
		log.Println("Error while protecting default branch", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	err = repo.ProtectBranch(project.ID, "feedback", gitlabModel.MaintainerPermissions)
	if err != nil {
		log.Println("Error while protecting feedback branch", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	assignmentProject.ProjectID = project.ID
	assignmentProject.ProjectStatus = database.Accepted

	if err = queryAssignmentProjects.WithContext(ctx).Save(assignmentProject); err != nil {
		log.Println("Error while setting Project to Accepted", err)
		return
	}
}

func waitForDefaultBranch(ctx context.Context, repo gitlab.Repository, projectID int, defaultBranch string) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout while waiting for default branch to be the same as the template project")
		case <-ticker.C:
			project, err := repo.GetProjectById(projectID)
			if err != nil {
				return err
			}
			if project.DefaultBranch == defaultBranch {
				return nil
			}
		}
	}
}

func waitForProtectedBranch(ctx context.Context, repo gitlab.Repository, projectID int, branch string) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout while waiting for protected branch to exist")
		case <-ticker.C:
			protectedBranchExists, err := repo.ProtectedBranchExists(projectID, branch)
			if err != nil {
				return err
			}
			if protectedBranchExists {
				return nil
			}
		}
	}
}
