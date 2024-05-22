package api

import (
	"github.com/gofiber/fiber/v2"
	apiV2 "gitlab.hs-flensburg.de/gitlab-classroom/controller/api_v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) CreatorMiddleware() fiber.Handler {
	var validateCreator apiV2.ValidateUserFunc = func(user database.UserClassrooms) bool {
		return user.UserID == user.Classroom.OwnerID
	}

	return ctrl.ValidateUserMiddleware(validateCreator)
}
