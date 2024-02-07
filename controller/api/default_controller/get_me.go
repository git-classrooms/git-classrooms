package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func (ctrl *DefaultController) GetMe(c *fiber.Ctx) error {
	userID, err := session.Get(c).GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	q := query.User
	user, err := q.WithContext(c.Context()).Where(q.ID.Eq(userID)).First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.JSON(user)
}
