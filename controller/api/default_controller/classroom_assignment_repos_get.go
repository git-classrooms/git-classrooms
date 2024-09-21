package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetMultipleProjectCloneUrls
// @Description	GetMultipleProjectCloneUrls
// @Id				GetMultipleProjectCloneUrls
// @Tags			assignment
// @Produce		json
// @Param			classroomId			path		string	true	"Classroom ID"			Format(uuid)
// @Param			assignmentProjectId	path		string	true	"Assignment Project ID"	Format(uuid)
// @Success		200					{array}		ProjectCloneUrlResponse
// @Failure		500					{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentProjectId}/repos [get]
func (ctrl *DefaultController) GetMultipleProjectCloneUrls(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()
	repo := ctx.GetGitlabRepository()

	assignmentProjects, err := query.AssignmentProjects.
		WithContext(c.Context()).
		Where(query.AssignmentProjects.AssignmentID.Eq(assignment.ID)).
		Where(query.AssignmentProjects.ProjectStatus.Eq(string(database.Accepted))).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := make([]*ProjectCloneUrlResponse, len(assignmentProjects))
	for i, project := range assignmentProjects {
		gitlabProject, err := repo.GetProjectById(project.ProjectID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		response[i] = &ProjectCloneUrlResponse{
			ProjectId:     project.ID,
			SshUrlToRepo:  gitlabProject.SSHURLToRepo,
			HttpUrlToRepo: gitlabProject.HTTPURLToRepo,
		}
	}

	return c.JSON(response)
}
