package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getMeResponse struct {
	*database.User
	GitlabURL string `json:"gitlabUrl"`
} //@Name GetMeResponse

// @Summary		Show your user account
// @Description	Get your user account
// @Id				GetMeV2
// @Tags			auth
// @Produce		json
// @Success		200	{object}	api.getMeResponse
// @Failure		401	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/me [get]
func (ctrl *DefaultController) GetMe(c *fiber.Ctx) error {
	queryUser := query.User
	user, err := queryUser.WithContext(c.Context()).
		Preload(queryUser.GitLabAvatar).
		Where(queryUser.ID.Eq(context.Get(c).GetUserID())).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := getMeResponse{
		User:      user,
		GitlabURL: "/api/v1/me/gitlab",
	}
	return c.JSON(response)
}
