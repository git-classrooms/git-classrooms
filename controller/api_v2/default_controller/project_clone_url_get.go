package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type ProjectCloneUrlResponse struct {
	ProjectId     int    `json:"projectId"`
	SshUrlToRepo  string `json:"sshUrlToRepo"`
	HttpUrlToRepo string `json:"httpUrlToRepo"`
}

//	@Summary		GetProjectCloneUrls
//	@Description	GetProjectCloneUrls
//	@Id				GetProjectCloneUrls
//	@Tags			default
//	@Produce		json
//	@Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
//	@Param			projectId	path		string	true	"Project ID"	Format(uuid)
//	@Success		200			{object}	ProjectCloneUrlResponse
//	@Failure		500			{object}	HTTPError
//	@Router			/classrooms/:classroomId/projects/:projectId/repo [get]
func (ctrl *DefaultController) GetProjectCloneUrls(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	projectId := ctx.GetGitlabProjectID()
	repo := ctx.GetGitlabRepository()

	project, err := repo.GetProjectById(projectId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := ProjectCloneUrlResponse{
		ProjectId:     project.ID,
		SshUrlToRepo:  project.SSHURLToRepo,
		HttpUrlToRepo: project.HTTPURLToRepo,
	}
	return c.JSON(response)
}
