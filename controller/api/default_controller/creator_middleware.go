package api

import (
	"github.com/gofiber/fiber/v2"

	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) CreatorMiddleware() fiber.Handler {
	var validateCreator apiController.ValidateUserFunc = func(user database.UserClassrooms) bool {
		return user.UserID == user.Classroom.OwnerID
	}

	return ctrl.ValidateUserMiddleware(validateCreator)
}
