package api

import (
	"github.com/gofiber/fiber/v2"
)

type getInfoGitlabResponse struct {
	GitlabURL string `json:"gitlabUrl"`
} // @Name GetInfoGitlabResponse

// @Summary		GetGitlabInfo
// @Description	GetGitlabInfo
// @Id				GetGitlabInfo
// @Tags			info
// @Produce		json
// @Success		200	{object}	api.getInfoGitlabResponse
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/info/gitlab [get]
func (ctrl *DefaultController) GetGitlabInfo(c *fiber.Ctx) error {
	response := getInfoGitlabResponse{
		GitlabURL: ctrl.config.GitLab.GetURL(),
	}

	return c.JSON(response)
}
