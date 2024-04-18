package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getMeResponse struct {
	*database.User
	GitlabWebURL string `json:"gitlabWeburl"`
} //@Name GetMe

// @Summary		Show your user account
// @Description	Get your user account
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	default_controller.getMeResponse
// @Failure		401	{object}	httputil.HTTPError
// @Failure		500	{object}	httputil.HTTPError
// @Router			/me [get]
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
		User:         user,
		GitlabWebURL: "/api/v1/me/gitlab",
	}
	return c.JSON(response)
}
