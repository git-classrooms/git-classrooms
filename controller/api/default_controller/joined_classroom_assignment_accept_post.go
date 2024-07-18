package default_controller

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Accept the assignment
// @Description	Accept the assignment and work on your repository
// @Id				AcceptAssignment
// @Tags			assignment
// @Param			classroomId		path	string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string	true	"Assignment ID"	Format(uuid)
// @Param			X-Csrf-Token	header	string	true	"Csrf-Token"
// @Success		201
// @Success		202
// @Header			202	{string}	Location	"/api/v1/classroom/joined/{classroomId}/assignments/{assignmentId}"
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/joined/{classroomId}/assignments/{assignmentId}/accept [post]
func (ctrl *DefaultController) AcceptAssignment(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()
	userID := ctx.GetUserID()
	team := classroom.Team
	assignmentProject := ctx.GetJoinedClassroomAssignment()

	if assignmentProject.ProjectStatus == database.Accepted {
		return c.SendStatus(fiber.StatusNoContent) // You or your teammate have already accepted the assignment
	}

	if assignmentProject.Assignment.DueDate.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "Assignment is already overdue")
	}

	repo := context.Get(c).GetGitlabRepository()

	if err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Check if template repository still exists
	_, err = repo.GetProjectById(assignmentProject.Assignment.TemplateProjectID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	project, err := repo.ForkProject(assignmentProject.Assignment.TemplateProjectID, gitlabModel.Private, team.GroupID, assignmentProject.Assignment.Name, assignmentProject.Assignment.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		if recover() != nil || err != nil {
			if err := repo.DeleteProject(project.ID); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	memberIds := utils.Map(assignmentProject.Team.Member, func(member *database.UserClassrooms) int {
		return member.UserID
	})

	gitlabMember := utils.Map(memberIds, func(member int) gitlabModel.User {
		return gitlabModel.User{ID: member}
	})

	_, err = repo.AddProjectMembers(project.ID, gitlabMember)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the project will be deleted

	_, err = repo.CreateBranch(project.ID, "feedback", "main")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the project will be deleted

	queryUsers := query.User
	members, err := queryUsers.
		WithContext(c.Context()).
		Where(queryUsers.ID.In(memberIds...)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	mentions := utils.Map(members, func(member *database.User) string {
		return fmt.Sprintf("/cc @%s", member.GitlabUsername)
	})
	description := fmt.Sprintf(mergeRequestDescription, strings.Join(mentions, "\n"))
	err = repo.CreateMergeRequest(project.ID, "main", "feedback", "Feedback", description, userID, classroom.Classroom.OwnerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the project will be deleted

	// In a few cases the main branch isn't available directly after the creation, this would cause an error when setting up protection rules for it, there we wait for the main branch to exist
	timeout := time.After(1 * time.Second)
	for keep_on_waiting := true; keep_on_waiting; {
		select {
		case <-timeout:
			log.Default().Println("Timeout while waiting for protected main branch to exist")
			keep_on_waiting = false
		default:
			protectedMainExists, err := repo.ProtectedBranchExists(project.ID, "main")
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			if protectedMainExists {
				keep_on_waiting = false
			} else {
				log.Default().Println("Protected main branch does not exists directly")
			}
		}
	}

	err = repo.UnprotectBranch(project.ID, "main")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the project will be deleted

	err = repo.ProtectBranch(project.ID, "main", gitlabModel.DeveloperPermissions)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the project will be deleted

	err = repo.ProtectBranch(project.ID, "feedback", gitlabModel.MaintainerPermissions)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the project will be deleted

	assignmentProject.ProjectID = project.ID
	assignmentProject.ProjectStatus = database.Accepted
	queryAssignmentProjects := query.AssignmentProjects
	err = queryAssignmentProjects.WithContext(c.Context()).Save(assignmentProject)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Location", fmt.Sprintf("/api/v1/classrooms/joined/%s/assignments/%s", classroom.ClassroomID.String(), assignmentProject.AssignmentID.String()))
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
