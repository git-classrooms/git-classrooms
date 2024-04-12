package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		Show your user account
// @Description	Get your user account
// @Tags			user
// @Accept			json
// @Produce		json
// @Success		200	{object}	database.User
// @Failure		401	{object}	httputil.HTTPError
// @Failure		500	{object}	httputil.HTTPError
// @Router			/me [get]
func (ctrl *DefaultController) GetMe(c *fiber.Ctx) error {

	// TODO: Add Avatar-URL and json tags to the User struct
	gitlabUser, err := context.Get(c).GetGitlabRepository().GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(gitlabUser)
}
