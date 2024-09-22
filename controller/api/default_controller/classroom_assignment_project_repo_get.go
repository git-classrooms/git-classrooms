package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type ProjectCloneUrlResponse struct {
	ProjectId     uuid.UUID `json:"projectId"`
	SshUrlToRepo  string    `json:"sshUrlToRepo"`
	HttpUrlToRepo string    `json:"httpUrlToRepo"`
}

// @Summary		GetProjectCloneUrls
// @Description	GetProjectCloneUrls
// @Id				GetProjectCloneUrls
// @Tags			project
// @Produce		json
// @Param			classroomId			path		string	true	"Classroom ID"			Format(uuid)
// @Param			assignmentId		path		string	true	"Assignment ID"			Format(uuid)
// @Param			assignmentProjectId	path		string	true	"Assignment Project ID"	Format(uuid)
// @Success		200					{object}	ProjectCloneUrlResponse
// @Failure		500					{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentId}/projects/{assignmentProjectId}/repo [get]
func (ctrl *DefaultController) GetProjectCloneUrls(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	project := ctx.GetAssignmentProject()
	repo := ctx.GetGitlabRepository()

	if project.ProjectStatus != database.Accepted {
		return fiber.ErrNotFound
	}

	gitlabProject, err := repo.GetProjectById(project.ProjectID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := ProjectCloneUrlResponse{
		ProjectId:     project.ID,
		SshUrlToRepo:  gitlabProject.SSHURLToRepo,
		HttpUrlToRepo: gitlabProject.HTTPURLToRepo,
	}
	return c.JSON(response)
}
