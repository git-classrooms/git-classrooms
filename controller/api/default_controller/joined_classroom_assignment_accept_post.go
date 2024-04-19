package default_controller

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) JoinAssignment(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()
	userID := ctx.GetUserID()
	team := classroom.Team
	assignmentProject := ctx.GetJoinedClassroomAssignment()

	if assignmentProject.AssignmentAccepted {
		return c.SendStatus(fiber.StatusNoContent) // You or your teammate have already accepted the assignment
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
	// We don't need to clean up this step because the branch will be deleted when the project is deleted

	_, err = repo.CreateBranch(project.ID, "feedback", "main")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the branch will be deleted when the project is deleted

	err = repo.UnprotectBranch(project.ID, "main")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the branch will be deleted when the project is deleted

	err = repo.ProtectBranch(project.ID, "main", gitlabModel.DeveloperPermissions)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the branch will be deleted when the project is deleted

	err = repo.ProtectBranch(project.ID, "feedback", gitlabModel.MaintainerPermissions)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// We don't need to clean up this step because the branch will be deleted when the project is deleted

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
	// We don't need to clean up this step because the branch will be deleted when the project is deleted

	assignmentProject.ProjectID = project.ID
	assignmentProject.AssignmentAccepted = true
	queryAssignmentProjects := query.AssignmentProjects
	err = queryAssignmentProjects.WithContext(c.Context()).Save(assignmentProject)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

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
