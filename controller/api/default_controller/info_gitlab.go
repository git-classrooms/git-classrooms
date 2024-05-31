package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"log"
)

type getInfoGitlabResponse struct {
	GitlabUrl string `json:"gitlabUrl"`
} //@Name getInfoGitlabResponse

// @Summary		getInfoGitlabResponse
// @Description	getInfoGitlabResponse
// @Id				getInfoGitlabResponse
// @Tags			info
// @Produce		json
// @Success		200	{object}	default_controller.getInfoGitlabResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/info/gitlab [get]
func (ctrl *DefaultController) GetGitlabInfo(c *fiber.Ctx) error {
	appConfig, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("failed to get application configuration", err)
	}

	response := getInfoGitlabResponse{
		GitlabUrl: appConfig.GitLab.GetURL(),
	}

	return c.JSON(response)
}
