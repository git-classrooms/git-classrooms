package api

import (
	"github.com/gofiber/fiber/v2"
	apiV2 "gitlab.hs-flensburg.de/gitlab-classroom/controller/api_v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) ArchivedMiddleware() fiber.Handler {
	var validateArchived apiV2.ValidateUserFunc = func(classroom database.UserClassrooms) bool {
		return !classroom.Classroom.Archived
	}

	return ctrl.ValidateUserMiddleware(validateArchived)
}
