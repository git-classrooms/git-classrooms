package api

import (
	"github.com/gofiber/fiber/v2"
)

type getInfoGitlabResponse struct {
	GitlabUrl string `json:"gitlabUrl"`
} //@Name GetInfoGitlabResponse

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
		GitlabUrl: ctrl.config.GitLab.GetURL(),
	}

	return c.JSON(response)
}
