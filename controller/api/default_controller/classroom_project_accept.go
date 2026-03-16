package api

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/logging"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Accept the assignment
// @Description	Accept the assignment and work on your repository
// @Id				AcceptAssignment
// @Tags			project
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
// @Param			projectId		path	string	true	"Project ID"	Format(uuid)
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
// @Success		201
// @Success		202
// @Header			202	{string}	Location	"/api/v1/classroom/{classroomId}/projects/{projectId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/projects/{projectId}/accept [post]
func (ctrl *DefaultController) AcceptAssignment(c *fiber.Ctx) (err error) {
	ctx := fiberContext.Get(c)
	log := ctx.GetLoggerForHandler("AcceptAssignment")
	classroom := ctx.GetUserClassroom()
	userID := ctx.GetUserID()
	assignmentProject := ctx.GetAssignmentProject()

	if assignmentProject.ProjectStatus == database.Accepted {
		log.Info("assignment is already accepted")
		return c.SendStatus(fiber.StatusNoContent) // You or your teammate have already accepted the assignment
	}

	if assignmentProject.ProjectStatus == database.Creating {
		return fiber.NewError(fiber.StatusForbidden, "The project is still being created")
	}

	if assignmentProject.Assignment.DueDate != nil && assignmentProject.Assignment.DueDate.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "The assignment is already over")
	}

	repo := ctx.GetGitlabRepository()

	log.Debug("authenticate the repo with group access token")

	if err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken, log); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Check if template repository still exists
	templateProject, err := repo.GetProjectById(assignmentProject.Assignment.TemplateProjectID)
	if err != nil {
		log.Error("error getting template project", "templateProjectID", assignmentProject.Assignment.TemplateProjectID, "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	log.Info("fetched templateProject from gitlab", "name", templateProject.Name, "templateProjectID", assignmentProject.Assignment.TemplateProjectID)

	assignmentProject.ProjectStatus = database.Creating

	log.Debug("change project status to creating")

	queryAssignmentProjects := query.AssignmentProjects
	if err = queryAssignmentProjects.WithContext(c.Context()).Save(assignmentProject); err != nil {
		log.Error("error saving assignmentProject from db", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Make the actual project creation and assignment acceptance async
	go ctrl.acceptAssignment(c.Context(), repo, userID, classroom.Classroom.OwnerID, templateProject, assignmentProject)

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s", classroom.ClassroomID.String(), assignmentProject.AssignmentID.String()))
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

func (ctrl *DefaultController) acceptAssignment(ctx context.Context, repo gitlab.Repository, userID int, classroomOwnerID int, templateProject *gitlabModel.Project, assignmentProject *database.AssignmentProjects) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Minute)
	defer cancel()

	log := logging.GetLogger(ctx).With(
		"templateProjectID", assignmentProject.Assignment.TemplateProjectID,
		"groupID", assignmentProject.Team.GroupID,
	)

	queryAssignmentProjects := query.AssignmentProjects

	var err error
	defer func() {
		if recover() != nil || err != nil {
			assignmentProject.ProjectStatus = database.Failed
			if err := queryAssignmentProjects.WithContext(ctx).
				Save(assignmentProject); err != nil {
				log.Error("error while setting project to failed!", "error", err)
			}
		}
	}()

	project, err := repo.ForkProjectWithOnlyDefaultBranch(assignmentProject.Assignment.TemplateProjectID, gitlabModel.Private, assignmentProject.Team.GroupID, assignmentProject.Assignment.Name, assignmentProject.Assignment.Description)
	if err != nil {
		log.Error("error while forking the template Project", "error", err)
		return
	}
	log = log.With("forkedProjectID", project.ID)
	defer func() {
		if recover() != nil || err != nil {
			if err := repo.DeleteProject(project.ID); err != nil {
				log.Error("error while deleting project", "error", err)
			}
		}
	}()

	log.Info("forked project for team")

	// wait till default branch of forked project is the same as the template project
	// this is necessary because the default branch is not immediately available after forking
	// TODO?: do not wait the whole 5 Minutes for this
	err = waitForDefaultBranch(ctx, repo, project.ID, templateProject.DefaultBranch)
	if err != nil {
		log.Error("error while waiting for defaultBranch", "defaultBranch", templateProject.DefaultBranch, "error", err)
		return
	}

	memberIds := utils.Map(assignmentProject.Team.Member, func(member *database.UserClassrooms) int {
		return member.UserID
	})

	gitlabMember := utils.Map(memberIds, func(member int) gitlabModel.User {
		return gitlabModel.User{ID: member}
	})

	log.Info(fmt.Sprintf("adding %d members to project", len(gitlabMember)))

	project, err = repo.AddProjectMembers(project.ID, gitlabMember)
	if err != nil {
		log.Error("error while adding members to the project", "error", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	log.Debug("create branch feedback")

	_, err = repo.CreateBranch(project.ID, "feedback", project.DefaultBranch)
	if err != nil {
		log.Error("error while creating feedback branch", "error", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	queryUsers := query.User
	members, err := queryUsers.
		WithContext(ctx).
		Where(queryUsers.ID.In(memberIds...)).
		Find()
	if err != nil {
		log.Error("error while fetching members", "error", err)
		return
	}

	log.Info("create merge request")

	mentions := utils.Map(members, func(member *database.User) string {
		return fmt.Sprintf("/cc @%s", member.GitlabUsername)
	})
	description := fmt.Sprintf(mergeRequestDescription, strings.Join(mentions, "\n"))
	err = repo.CreateMergeRequest(project.ID, project.DefaultBranch, "feedback", "Feedback", description, userID, classroomOwnerID)
	if err != nil {
		log.Error("error while creating merge request", "error", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	// In a few cases the main branch isn't available directly after the creation, this would cause an error when setting up protection rules for it, there we wait for the default branch to exist
	// TODO?: do not wait the whole 5 Minutes for this
	err = waitForProtectedBranch(ctx, repo, project.ID, project.DefaultBranch)
	if err != nil {
		log.Error("error while waiting for protected main branch", "error", err)
		return
	}

	log.Debug("unprotect branch on forked project", "branch", project.DefaultBranch)

	err = repo.UnprotectBranch(project.ID, project.DefaultBranch)
	if err != nil {
		log.Error("Error while unprotecting default branch", "error", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	log.Debug("protect branch on forked project to developers", "branch", project.DefaultBranch)

	err = repo.ProtectBranch(project.ID, project.DefaultBranch, gitlabModel.DeveloperPermissions)
	if err != nil {
		log.Error("error while protecting default branch", "error", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	log.Debug("protect feedback branch on forked project to maintainers")

	err = repo.ProtectBranch(project.ID, "feedback", gitlabModel.MaintainerPermissions)
	if err != nil {
		log.Error("error while protecting feedback branch", "error", err)
		return
	}
	// We don't need to clean up this step because the project will be deleted

	assignmentProject.ProjectID = project.ID
	assignmentProject.ProjectStatus = database.Accepted
	assignmentProject.SSHURLToRepo = project.SSHURLToRepo
	assignmentProject.HTTPURLToRepo = project.HTTPURLToRepo

	log.Info("setting project to accepted")

	if err = queryAssignmentProjects.WithContext(ctx).Save(assignmentProject); err != nil {
		log.Error("error while setting Project to Accepted", "error", err)
		return
	}
}

func waitForDefaultBranch(ctx context.Context, repo gitlab.Repository, projectID int, defaultBranch string) error {
	log := logging.GetLogger(ctx)

	start := time.Now()
	log.Debug("wait until default branch is available on forked project")

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			err := errors.New("timeout while waiting for default branch to be the same as the template project")
			log.Error("timeout waiting for default branch", "duration", time.Now().Sub(start), "error", err)
			return err
		case <-ticker.C:
			project, err := repo.GetProjectById(projectID)
			if err != nil {
				log.Error("error getting project from gitlab", "duration", time.Now().Sub(start), "error", err)
				return err
			}
			if project.DefaultBranch == defaultBranch {
				log.Debug("project default branch is available", "duration", time.Now().Sub(start))
				return nil
			}
		}
	}
}

func waitForProtectedBranch(ctx context.Context, repo gitlab.Repository, projectID int, branch string) error {
	log := logging.GetLogger(ctx)

	start := time.Now()
	log.Debug("wait until branch is protected on forked project")

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			err := errors.New("timeout while waiting for protected branch to exist")
			log.Error("timeout waiting for protected branch", "duration", time.Now().Sub(start), "error", err)
			return err
		case <-ticker.C:
			protectedBranchExists, err := repo.ProtectedBranchExists(projectID, branch)
			if err != nil {
				log.Error("error getting protected branch from gitlab", "duration", time.Now().Sub(start), "error", err)
				return err
			}
			if protectedBranchExists {
				log.Debug("protected branch is available", "duration", time.Now().Sub(start))
				return nil
			}
		}
	}
}
