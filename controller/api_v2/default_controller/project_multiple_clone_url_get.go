package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type MultipleProjectCloneUrlResponse struct {
	Projects []ProjectCloneUrlResponse `json:"projects"`
}

// @Summary		GetMultipleProjectCloneUrls
// @Description	GetMultipleProjectCloneUrls
// @Id				GetMultipleProjectCloneUrls
// @Tags			assignment
// @Produce		json
// @Param			classroomId			path		string	true	"Classroom ID"			Format(uuid)
// @Param			assignmentProjectId	path		string	true	"Assignment Project ID"	Format(uuid)
// @Success		200					{object}	ProjectCloneUrlResponse
// @Failure		500					{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentProjectId}/repos [get]
func (ctrl *DefaultController) GetMultipleProjectCloneUrls(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()
	repo := ctx.GetGitlabRepository()

	assignmentProjects, err := query.AssignmentProjects.WithContext(c.Context()).Where(query.AssignmentProjects.AssignmentID.Eq(assignment.ID)).Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response MultipleProjectCloneUrlResponse
	for _, assignment_project := range assignmentProjects {
		project, err := repo.GetProjectById(assignment_project.ProjectID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		response.Projects = append(response.Projects, ProjectCloneUrlResponse{
			ProjectId:     project.ID,
			SshUrlToRepo:  project.SSHURLToRepo,
			HttpUrlToRepo: project.HTTPURLToRepo,
		})
	}

	return c.JSON(response)
}
